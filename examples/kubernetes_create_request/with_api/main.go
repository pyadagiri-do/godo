// This example calls the real DO API and logs the request body via a custom RoundTripper.
// Requires DIGITALOCEAN_ACCESS_TOKEN env var. Use with caution - it will attempt to create a cluster.
//
// Run: go run ./examples/kubernetes_create_request/with_api/
//
// To test without creating: use invalid params (e.g. invalid region) to trigger validation error
// before the cluster is provisioned, or use a token with no balance to get a different error.
package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

// loggingRoundTripper logs the request body before delegating to the real transport
type loggingRoundTripper struct {
	transport http.RoundTripper
}

func (l *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body != nil && req.URL.Path == "/v2/kubernetes/clusters" && req.Method == http.MethodPost {
		body, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(body))
		println("\n--- HTTP Request ---")
		println("URL:", req.Method, req.URL.String())
		println("Body:")
		println(string(body))
		println("--- End Request ---\n")
	}
	return l.transport.RoundTrip(req)
}

func main() {
	token := os.Getenv("DIGITALOCEAN_ACCESS_TOKEN")
	if token == "" {
		println("Set DIGITALOCEAN_ACCESS_TOKEN to run this example")
		os.Exit(1)
	}

	oauthClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
	}))

	// Wrap with our logging transport
	oauthClient.Transport = &loggingRoundTripper{transport: oauthClient.Transport}

	client := godo.NewClient(oauthClient)

	// Create request with HA nil (omitted)
	req := &godo.KubernetesClusterCreateRequest{
		Name:        "godo-ha-test-3" + os.Getenv("USER"),
		RegionSlug:  "nyc1",
		VersionSlug: "1.34",
		NodePools: []*godo.KubernetesNodePoolCreateRequest{
			{Name: "pool-1", Size: "s-1vcpu-2gb", Count: 1},
		},
		// HA not set - will be omitted from request, API applies version-based default
		HA: godo.PtrTo(false),
	}

	ctx := context.Background()
	cluster, resp, err := client.Kubernetes.Create(ctx, req)
	if err != nil {
		println("Error (expected if invalid token/params):", err.Error())
		if resp != nil && resp.Body != nil {
			body, _ := io.ReadAll(resp.Body)
			println("Response body:", string(body))
		}
		os.Exit(1)
	}

	println("Created cluster:", cluster.ID, cluster.Name, "HA:", cluster.HA)
}
