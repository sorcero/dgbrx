package dgraph

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func doSlashRequest(dgraph *DGraph, query string) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/slash", dgraph.Endpoint)
	logger.Debugf("Preparing request for backup of Dgraph Database: %s with query %s", endpoint, query)

	marshal, err := json.Marshal(GenericQuery{Query: query})
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

func DropAll(dgraph *DGraph) error {
	d, err := json.Marshal(map[string]bool{"drop_all": true})
	endpoint := fmt.Sprintf("%s/alter", dgraph.Endpoint)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"POST",
		endpoint,
		bytes.NewBuffer(d),
	)
	req.Header.Set("Dg-Auth", dgraph.ApiToken)
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return err
	}
	logger.Debugf("Authenticating")

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	logger.Debugf("Drop All request succeeded")
	return nil
}
