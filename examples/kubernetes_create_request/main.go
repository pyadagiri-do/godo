// This example demonstrates the JSON request body for Kubernetes cluster creation
// with different HA configurations. Run with: go run main.go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/digitalocean/godo"
)

func main() {
	// Case 1: HA nil (omitted) - API will apply version-based default
	req1 := &godo.KubernetesClusterCreateRequest{
		Name:        "test-cluster",
		RegionSlug:  "nyc1",
		VersionSlug: "1.36",
		NodePools: []*godo.KubernetesNodePoolCreateRequest{
			{Name: "pool-1", Size: "s-1vcpu-2gb", Count: 1},
		},
	}
	body1, _ := json.MarshalIndent(req1, "", "  ")
	fmt.Println("=== HA nil (omitted) - request body ===")
	fmt.Println(string(body1))
	fmt.Println()

	// Case 2: HA explicitly true
	req2 := &godo.KubernetesClusterCreateRequest{
		Name:        "test-cluster",
		RegionSlug:  "nyc1",
		VersionSlug: "1.36",
		HA:         godo.PtrTo(true),
		NodePools: []*godo.KubernetesNodePoolCreateRequest{
			{Name: "pool-1", Size: "s-1vcpu-2gb", Count: 1},
		},
	}
	body2, _ := json.MarshalIndent(req2, "", "  ")
	fmt.Println("=== HA explicitly true - request body ===")
	fmt.Println(string(body2))
	fmt.Println()

	// Case 3: HA explicitly false
	req3 := &godo.KubernetesClusterCreateRequest{
		Name:        "test-cluster",
		RegionSlug:  "nyc1",
		VersionSlug: "1.36",
		HA:         godo.PtrTo(false),
		NodePools: []*godo.KubernetesNodePoolCreateRequest{
			{Name: "pool-1", Size: "s-1vcpu-2gb", Count: 1},
		},
	}
	body3, _ := json.MarshalIndent(req3, "", "  ")
	fmt.Println("=== HA explicitly false - request body ===")
	fmt.Println(string(body3))
}
