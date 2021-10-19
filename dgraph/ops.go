package dgraph

import "github.com/urfave/cli/v2"

// New creates a dgraph cloud instance from provided
// endpoint and api key
func New(endpoint string, apiKey string) *DGraph {
	return &DGraph{
		Endpoint: endpoint,
		ApiToken: apiKey,
	}
}

// FromCliContext creates a dgraph cloud instance from
// cli.Context
func FromCliContext(cli *cli.Context) *DGraph {
	return &DGraph{
		Endpoint: cli.String("url"),
		ApiToken: cli.String("api-key"),
	}
}
