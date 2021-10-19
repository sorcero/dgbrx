package dgraph

import "time"

// GenericError is a
type GenericError struct {
	Message string `json:"message"`
}

// DGraph instance definition for providing
// the endpoint and the Api Token
type DGraph struct {
	Endpoint string
	ApiToken string
}

// ResponseDataPayload is generic struct received
// from dgraph
type ResponseDataPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ExportDataPayload returns the export Id and task Id
// which can then be used to poll for jobs
type ExportDataPayload struct {
	ExportId string `json:"exportId"`
	TaskId   string `json:"taskId"`
}

// AdminExportDataPayload is a wrapper for ExportDataPayload
type AdminExportDataPayload struct {
	Export ExportDataPayload `json:"export"`
}

// AdminExportPromise is a wrapper for AdminExportDataPayload along with
// details for any errors that might have happened
type AdminExportPromise struct {
	Data  AdminExportDataPayload `json:"data"`
	Error []GenericError         `json:"errors,omitempty"`
}

// AdminExportStatus provides the response from dgraph
// which indicates the status of the backup,
// and URLs over which the backup can be downloaded
type AdminExportStatus struct {
	Kind        string    `json:"kind"`
	LastUpdated time.Time `json:"lastUpdated"`
	SignedUrls  []string  `json:"signedUrls"`
	Status      string    `json:"status"`
}

// AdminExportPayload is a wrapper for AdminExportStatus
type AdminExportPayload struct {
	ExportStatus AdminExportStatus `json:"exportStatus"`
}

// AdminExport is a wrapper AdminExportPayload along with
// details for any errors that might have happened
type AdminExport struct {
	Data  AdminExportPayload `json:"data"`
	Error []GenericError     `json:"errors,omitempty"`
}

// GenericQuery helps to send a generic query to dgraph
type GenericQuery struct {
	Query string `json:"query"`
}
