package dgraph

const graphExportMutationQuery = `mutation {
  export {
	response { code message }
	exportId
	taskId
  }
}`

const graphExportStatusCheckMutationQuery = `query {
  exportStatus (
    exportId:"%s"
    taskId: "%s"
  ){
    kind
    lastUpdated
    signedUrls
    status
  }
}`

const dropAllMutationQuery = `mutation {
  dropData(allData: true) {
	response { code message }
  }
}`
