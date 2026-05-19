# Kubernetes HA create example

Calls the real DigitalOcean API for six `KubernetesClusterCreateRequest` scenarios and saves every HTTP request/response to one JSON file. Used to verify `ha` behavior across versions `1.35` / `1.36` with HA unset, `true`, or `false`.

| Version | HA values tested      |
|---------|------------------------|
| `1.35`  | unset, `true`, `false` |
| `1.36`  | unset, `true`, `false` |

## Run (published godo)

```bash
cd examples/kubernetes_create_request

export DIGITALOCEAN_ACCESS_TOKEN=your_token
export OUTPUT_FILE=./captures.json   # optional, default: captures.json

go run .
```

Pin or upgrade godo:

```bash
go get github.com/digitalocean/godo@latest
go run .
```

## Run (local godo checkout)

```bash
cd examples/kubernetes_create_request

go mod edit -replace=github.com/digitalocean/godo=../..
go mod tidy

export DIGITALOCEAN_ACCESS_TOKEN=your_token
go run .
```

Remove the replace when finished:

```bash
go mod edit -dropreplace=github.com/digitalocean/godo
go mod tidy
```

## Output

`captures.json` contains a `runs` array. Each entry has `scenario`, `version`, `ha`, `request`, `response` (`status`, `status_code`, `body`), and `cluster_id` / `cluster_ha` on success.

## Warnings

- Creates up to **6 billable clusters**. Delete them when done.
- Names look like `godo-ha-v1-35-ha-unset-<user>` (lowercase, hyphens only).

To capture traffic without provisioning, use an invalid token or bad params — failed runs are still recorded.
