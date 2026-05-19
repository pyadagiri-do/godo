# Kubernetes Create Request Examples

Examples for testing the `ha` field on `KubernetesClusterCreateRequest` (versions `1.35` / `1.36`, HA unset / `true` / `false`).

## Option 1: Print JSON (no API call, no token)

Shows the JSON bodies godo would send for each HA setting. Uses your **local godo checkout** when run from the repo root:

```bash
# from godo repo root
go run ./examples/kubernetes_create_request/main.go
```

Output:
- **HA unset** (`nil`): `"ha"` omitted — API applies a version-based default
- **HA true**: `"ha": true`
- **HA false**: `"ha": false`

## Option 2: Call the real API (6 scenarios)

`with_api/` runs all six combinations, calls the DigitalOcean API, and writes request/response captures to a single file.

| Version | HA values tested      |
|---------|------------------------|
| `1.35`  | unset, `true`, `false` |
| `1.36`  | unset, `true`, `false` |

### Published godo (default)

`with_api/` is a small standalone module that depends on a released `github.com/digitalocean/godo` version (see `with_api/go.mod`).

```bash
cd examples/kubernetes_create_request/with_api

export DIGITALOCEAN_ACCESS_TOKEN=your_token
export OUTPUT_FILE=./captures.json   # optional, default: captures.json

go run .
```

To pin or upgrade the library version:

```bash
go get github.com/digitalocean/godo@latest   # or @v1.192.0
go run .
```

### Local godo checkout

To exercise **unreleased changes** in this repo, point the example module at the parent checkout:

```bash
cd examples/kubernetes_create_request/with_api

go mod edit -replace=github.com/digitalocean/godo=../../../
go mod tidy

export DIGITALOCEAN_ACCESS_TOKEN=your_token
go run .
```

Remove the replace when you want published godo again:

```bash
go mod edit -dropreplace=github.com/digitalocean/godo
go mod tidy
```

### Output

One file (default `captures.json`) with a `runs` array. Each entry includes:

- `scenario`, `version`, `ha`
- `request` — HTTP body sent to `POST /v2/kubernetes/clusters`
- `response` — `status`, `status_code`, `body`
- `cluster_id`, `cluster_ha` on success

### Warnings

- Creates up to **6 real clusters** (billable). Delete them when finished.
- Cluster names look like `godo-ha-v1-35-ha-unset-<user>` (lowercase, hyphens only).

To inspect traffic without provisioning, use an invalid token or invalid params (e.g. bad region) so the API returns an error before create completes — captures are still written for failed runs.
