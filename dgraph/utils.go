package dgraph

import (
	"net/url"
	"strings"
)

const slashPort = "443"

// GetGrpcEndpoint converts a graphql or admin endpoint to a GRPC
// endpoint, as is commonly specified by the Dgraph Cloud Service
func GetGrpcEndpoint(dg *DGraph) (string, error) {
	u, err := url.Parse(dg.Endpoint)
	if err != nil {
		return "", err
	}

	urlParts := strings.SplitN(u.Host, ".", 2)

	return urlParts[0] + ".grpc." + urlParts[1] + ":" + slashPort, nil
}
