package ops

import "gitlab.com/sorcero/community/dgbrx/dgraph"

// Clean removes the schema and all the content from
// the Dgraph Database
func Clean(dg *dgraph.DGraph) error {
	return dgraph.DropAll(dg)
}
