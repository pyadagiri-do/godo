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

## Option 2: Call real API and log request

Uses a custom `http.RoundTripper` to log the actual HTTP request body before sending:

```bash
export DIGITALOCEAN_ACCESS_TOKEN=your_token
go run ./examples/kubernetes_create_request/with_api/
```

**Warning**: This will attempt to create a real cluster (costs money). To see the request without creating:
- Use an invalid token to get an auth error
- Use invalid params (e.g. non-existent region) to get a validation error before provisioning
