package ops

import (
	"fmt"
	"gitlab.com/sorcero/community/dgbrx/dgraph"
	"io"
	"net/http"
	"os"
)

// Backup creates a backup request to Dgraph Instance, and polls for task
// completion provided the export request id and waits until the entire backup is complete
// once the backup is complete, it is written to local disk
func Backup(dg *dgraph.DGraph) error {
	promise, err := dgraph.ExportBackupQueueJob(dg)
	if err != nil {
		return err
	}
	data, err := dgraph.AwaitExportBackup(dg, promise)
	if err != nil {
		return err
	}
	urls := data.Data.ExportStatus.SignedUrls
	outputNames := []string{"g01.gql_schema.gz", "g01.json.gz", "g01.schema.gz"}
	for i := range urls {
		fmt.Println(urls[i])
		out, err := os.Create(outputNames[i])
		if err != nil {
			return err
		}
		defer out.Close()

		get, err := http.Get(urls[i])
		if err != nil {
			return err
		}
		defer get.Body.Close()
		n, err := io.Copy(out, get.Body)
		if err != nil {
			return err
		}
		fmt.Println(n, "bytes written for", outputNames[i])

	}
	return nil
}
