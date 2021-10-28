package dgraph

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func doGraphRequest(dgraph *DGraph, query interface{}, suffix string) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/%s", dgraph.Endpoint, suffix)
	logger.Debugf("Preparing request for backup of Dgraph Database: %s with query %s", endpoint, query)

	marshal, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		endpoint,
		bytes.NewBuffer(marshal),
	)
	req.Header.Set("Dg-Auth", dgraph.ApiToken)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	logger.Debugf("Authenticating")

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	logger.Debugf("Backup request succeeded, parsing response")
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	logger.Debugf("Backup request read successfully")
	return d, nil
}

func doSlashRequest(dgraph *DGraph, query string) ([]byte, error) {
	return doGraphRequest(dgraph, GenericQuery{Query: query}, "admin/slash")
}

func doQueryRequest(dgraph *DGraph, query string) ([]byte, error) {
	return doGraphRequest(dgraph, map[string]string{"query": query}, "query")
}

func ExportBackupQueueJob(dgraph *DGraph) (*AdminExportPromise, error) {
	logger.Debugf("Preparing backup queue")
	d, err := doSlashRequest(dgraph, graphExportMutationQuery)
	if err != nil {
		return nil, err
	}
	logger.Debugf("Parsing export queue response")
	admin := &AdminExportPromise{}
	fmt.Println(string(d))
	err = json.Unmarshal(d, admin)
	if err != nil {
		return nil, err
	}
	logger.Debugf("Parse export queue success")
	if len(admin.Error) != 0 {
		for i := range admin.Error {
			logger.Warn(admin.Error[i].Message)
		}
		return nil, errors.New(admin.Error[0].Message)
	}
	return admin, nil

}

func AwaitExportBackup(dgraph *DGraph, e *AdminExportPromise) (*AdminExport, error) {
	wait := 1
	for {
		logger.Debugf("Awaiting export backup")
		d, err := doSlashRequest(
			dgraph,
			fmt.Sprintf(graphExportStatusCheckMutationQuery, e.Data.Export.ExportId, e.Data.Export.TaskId))
		fmt.Println(string(d))
		if err != nil {
			return nil, err
		}
		logger.Debugf("Parsing received export backup data")
		admin := &AdminExport{}
		err = json.Unmarshal(d, admin)
		if err != nil {
			return nil, err
		}

		if len(admin.Error) != 0 {
			for i := range admin.Error {
				logger.Warn(admin.Error[i].Message)
			}
			return nil, errors.New(admin.Error[0].Message)
		}
		if admin.Data.ExportStatus.Status == "Running" {
			// the backup is still in progress
			// exponentially increase wait time
			if wait > 30 {
				// we waited for 30 retries
				return nil, errors.New(fmt.Sprintf("Waiting for backup timed out after %d retries", wait))
			}
			time.Sleep(time.Duration(wait*2) * time.Second)
			wait += 1
		} else {
			return admin, nil
		}

	}

}

/*func DropAll(dgraph *DGraph) error {
	logger.Debugf("Preparing drop all")


	ctx := context.Background()
	conn, err := dgo.DialCloud(dgraph.Endpoint, dgraph.ApiToken)
	if err != nil { return err }
	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))
	err = dg.Login(ctx, "groot", "password")
	if err != nil { return err }
	err = dg.Alter(ctx, &api.Operation{DropAll: true})
	if err != nil {
		return err
	}
	logger.Infof("Drop databases completely.")
	return nil
}*/

func DropPredicate(ctx context.Context, dg *dgo.Dgraph, predicate string) error {
	if strings.HasPrefix(predicate, "dgraph") {
		// some internal dgraph schema which we would like to omit
		logger.Infof("Skipping predicate '%s'.", predicate)
		return nil
	}
	logger.Debugf("Removing predicate '%s'", predicate)

	err := dg.Alter(ctx, &api.Operation{DropAttr: predicate})
	if err != nil {
		return err
	}
	logger.Infof("Removed predicate '%s'.", predicate)
	return nil
}

func DropAll(dgraph *DGraph) error {
	logger.Debugf("Preparing to dropAll data")
	d, err := doSlashRequest(dgraph, dropAllMutationQuery)
	if err != nil {
		return err
	}
	logger.Debugf("Fetching current schema")
	d, err = doQueryRequest(dgraph, "schema {}")
	if err != nil {
		return err
	}

	admin := &PredicateQuery{}
	logger.Debug(string(d))
	err = json.Unmarshal(d, admin)
	if err != nil {
		return err
	}
	if len(admin.Error) != 0 {
		for i := range admin.Error {
			logger.Warn(admin.Error[i].Message)
		}
		return errors.New(admin.Error[0].Message)
	}
	logger.Debugf("Preparing to drop schema predicates")
	ctx := context.Background()
	conn, err := dgo.DialCloud(dgraph.Endpoint, dgraph.ApiToken)
	if err != nil {
		return err
	}
	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	for i := range admin.Data.Schema {
		if err := DropPredicate(ctx, dg, admin.Data.Schema[i].Predicate); err != nil {
			return err
		}
	}
	logger.Infof("Drop all data completed successfully")
	return nil
}
