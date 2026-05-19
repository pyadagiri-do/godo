# Kubernetes Create Request Examples

Examples for testing the HA field behavior in `KubernetesClusterCreateRequest`.

## Option 1: Print JSON (no API call, no token needed)

Shows the request body that would be sent for different HA configurations:

```bash
go run ./examples/kubernetes_create_request/main.go
```

Output shows:
- **HA nil**: `"ha"` is omitted from JSON (API applies version-based default)
- **HA true**: `"ha": true` in JSON
- **HA false**: `"ha": false` in JSON

## Option 2: Call real API and capture request/response (6 scenarios)

Runs all combinations of version (`1.35`, `1.36`) and HA (`unset`, `true`, `false`), writes one JSON file with every request/response:

```bash
export DIGITALOCEAN_ACCESS_TOKEN=your_token
export OUTPUT_FILE=./captures.json   # optional, default: captures.json
go run ./examples/kubernetes_create_request/with_api/
```

Output: `captures.json` with a `runs` array — each entry has scenario metadata, HTTP request body, response status/body, and cluster `ha` when created.

**Warning**: This creates up to **6 real clusters** (costs money). To see the request without creating:
- Use an invalid token to get an auth error
- Use invalid params (e.g. non-existent region) to get a validation error before provisioning
