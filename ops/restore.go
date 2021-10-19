package ops

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/codeclysm/extract/v3"
	"gitlab.com/sorcero/community/dgbrx/dgraph"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Entry *map[string]string

func extractArchive(kind string, filepath string) (string, error) {
	jsonDir, err := os.MkdirTemp(".dgbrx_temp", fmt.Sprintf("dgbrx_%s", kind))
	if err != nil {
		return "", err
	}

	r, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	targetFile := path.Join(jsonDir, fmt.Sprintf("g01.%s", kind))

	err = extract.Gz(context.Background(), r, targetFile, nil)
	if err != nil {
		return "", err
	}
	logger.Debugf("Writing to %s file %s", kind, targetFile)
	return targetFile, nil
}

// Restore will restore the data from local path to the provided dgraph instance
func Restore(dg *dgraph.DGraph, i InputOptions, o OutputOptions) error {
	logger.Infof("Creating temporary directories")
	err := os.MkdirAll(".dgbrx_temp", 0o755)
	if err != nil {
		return err
	}

	jsonFile, err := extractArchive("json", i.JsonPath)
	if err != nil {
		return err
	}

	schemaFile, err := extractArchive("schema", i.SchemaPath)
	if err != nil {
		return err
	}

	j, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return err
	}
	var entry []map[string]interface{}

	err = json.Unmarshal(j, &entry)
	if err != nil {
		return err
	}

	var cleanedEntry []map[string]interface{}
	for i := range entry {
		e := entry[i]
		skip := false
		for k, _ := range e {
			if strings.HasPrefix(k, "dgraph") {
				skip = true
				break
			}
		}
		if skip {
			continue
		} else {
			cleanedEntry = append(cleanedEntry, e)
		}
	}

	b, err := json.Marshal(cleanedEntry)
	if err != nil {
		return err
	}
	cleanedJsonFile := fmt.Sprintf("%s-cleaned.json", jsonFile)
	err = ioutil.WriteFile(cleanedJsonFile, b, 0o644)
	if err != nil {
		return err
	}

	grpcEndpoint, err := dgraph.GetGrpcEndpoint(dg)
	if err != nil {
		return err
	}
	cmd := exec.Command("dgraph", "live", fmt.Sprintf("--slash_grpc_endpoint=%s", grpcEndpoint), "-f", cleanedJsonFile, "-t", dg.ApiToken, "-s", schemaFile)
	fmt.Println("COMMAND", cmd.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()

	// return nil
}
