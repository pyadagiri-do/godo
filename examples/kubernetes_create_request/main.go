// Runs 6 Kubernetes create scenarios (versions 1.35/1.36 × HA unset/true/false),
// calls the real DO API, and writes all HTTP request/response captures to one JSON file.
//
// Requires DIGITALOCEAN_ACCESS_TOKEN. Warning: creates up to 6 real clusters.
//
// Run:
//
//	cd examples/kubernetes_create_request
//	export DIGITALOCEAN_ACCESS_TOKEN=your_token
//	export OUTPUT_FILE=./captures.json   # optional
//	go run .
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

type scenario struct {
	label   string
	version string
	ha      *bool // nil = omitted from JSON
}

var scenarios = []scenario{
	{label: "v1-35-ha-unset", version: "1.35", ha: nil},
	{label: "v1-35-ha-true", version: "1.35", ha: godo.PtrTo(true)},
	{label: "v1-35-ha-false", version: "1.35", ha: godo.PtrTo(false)},
	{label: "v1-36-ha-unset", version: "1.36", ha: nil},
	{label: "v1-36-ha-true", version: "1.36", ha: godo.PtrTo(true)},
	{label: "v1-36-ha-false", version: "1.36", ha: godo.PtrTo(false)},
}

// captureRoundTripper buffers the latest create request/response so godo can still decode the body.
type captureRoundTripper struct {
	transport      http.RoundTripper
	lastReq        []byte
	lastResp       []byte
	lastStatus     string
	lastStatusCode int
}

func (c *captureRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	logCreate := req.Method == http.MethodPost && req.URL.Path == "/v2/kubernetes/clusters"

	var reqBody []byte
	if logCreate && req.Body != nil {
		reqBody, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(reqBody))
		c.lastReq = reqBody
	}

	resp, err := c.transport.RoundTrip(req)
	if err != nil || !logCreate || resp == nil || resp.Body == nil {
		return resp, err
	}

	respBody, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(respBody))

	c.lastResp = respBody
	c.lastStatus = resp.Status
	c.lastStatusCode = resp.StatusCode

	return resp, err
}

func haLabel(ha *bool) string {
	if ha == nil {
		return "unset"
	}
	if *ha {
		return "true"
	}
	return "false"
}

type scenarioRun struct {
	Scenario  string          `json:"scenario"`
	Version   string          `json:"version"`
	HA        string          `json:"ha"`
	ClusterID string          `json:"cluster_id,omitempty"`
	ClusterHA bool            `json:"cluster_ha,omitempty"`
	Error     string          `json:"error,omitempty"`
	Request   json.RawMessage `json:"request"`
	Response  struct {
		Status     string          `json:"status"`
		StatusCode int             `json:"status_code"`
		Body       json.RawMessage `json:"body"`
	} `json:"response"`
}

func main() {
	token := os.Getenv("DIGITALOCEAN_ACCESS_TOKEN")
	if token == "" {
		println("Set DIGITALOCEAN_ACCESS_TOKEN to run this example")
		os.Exit(1)
	}

	outputFile := os.Getenv("OUTPUT_FILE")
	if outputFile == "" {
		outputFile = "captures.json"
	}
	if err := os.MkdirAll(filepath.Dir(outputFile), 0o755); err != nil && filepath.Dir(outputFile) != "." {
		fmt.Println("mkdir:", err)
		os.Exit(1)
	}

	oauthClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
	}))

	rt := &captureRoundTripper{transport: oauthClient.Transport}
	oauthClient.Transport = rt

	client := godo.NewClient(oauthClient)
	ctx := context.Background()
	user := os.Getenv("USER")
	if user == "" {
		user = "local"
	}

	runs := make([]scenarioRun, 0, len(scenarios))

	for i, sc := range scenarios {
		createReq := &godo.KubernetesClusterCreateRequest{
			Name:        fmt.Sprintf("godo-ha-%s-%s", sc.label, user),
			RegionSlug:  "nyc1",
			VersionSlug: sc.version,
			NodePools: []*godo.KubernetesNodePoolCreateRequest{
				{Name: "pool-1", Size: "s-1vcpu-2gb", Count: 1},
			},
			HA: sc.ha,
		}

		fmt.Printf("\n[%d/6] %s (version=%s ha=%s)\n", i+1, sc.label, sc.version, haLabel(sc.ha))

		cluster, _, err := client.Kubernetes.Create(ctx, createReq)

		run := scenarioRun{
			Scenario: sc.label,
			Version:  sc.version,
			HA:       haLabel(sc.ha),
			Request:  json.RawMessage(rt.lastReq),
		}
		run.Response.Status = rt.lastStatus
		run.Response.StatusCode = rt.lastStatusCode
		run.Response.Body = json.RawMessage(rt.lastResp)

		if err != nil {
			run.Error = err.Error()
			fmt.Println("  error:", err)
		} else {
			run.ClusterID = cluster.ID
			run.ClusterHA = cluster.HA
			fmt.Printf("  created: id=%s ha=%v\n", cluster.ID, cluster.HA)
		}
		runs = append(runs, run)
	}

	out, err := json.MarshalIndent(map[string]interface{}{"runs": runs}, "", "  ")
	if err != nil {
		fmt.Println("marshal:", err)
		os.Exit(1)
	}
	if err := os.WriteFile(outputFile, out, 0o644); err != nil {
		fmt.Println("write:", err)
		os.Exit(1)
	}

	fmt.Printf("\nDone. All captures written to %s\n", outputFile)

	var failed int
	for _, r := range runs {
		if r.Error != "" {
			failed++
		}
	}
	if failed > 0 {
		fmt.Printf("  %d/%d scenarios returned an error\n", failed, len(scenarios))
		os.Exit(1)
	}
}
