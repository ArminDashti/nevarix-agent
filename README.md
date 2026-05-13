# nevarix-agent

HTTP API for agent management (`POST /api/v1/`, health, stats). OpenAPI: `api/openapi.yaml`.

## Docker Hub CI/CD

Pushes to `main` and version tags trigger [Docker Build and Push](.github/workflows/docker.yml). The workflow uses Docker Buildx, pushes to Docker Hub, and uses GitHub Actions cache for image layers.

### Required GitHub secrets

Configure these in the repository **Settings → Secrets and variables → Actions**:

| Secret | Purpose |
|--------|---------|
| `DOCKER_USERNAME` | Docker Hub username or organization name used in the image reference. |
| `DOCKER_PASSWORD` | **Docker Hub access token** (recommended), not your account password. Create under [Docker Hub → Account Settings → Security](https://hub.docker.com/settings/security). |

Images are pushed as:

- `<DOCKER_USERNAME>/<repository-name>:latest` — only for pushes to the default branch (`main`).
- `<DOCKER_USERNAME>/<repository-name>:sha-<short>` — every qualifying push (short Git commit SHA).
- `<DOCKER_USERNAME>/<repository-name>:<tag>` — when the push is a Git tag (for example `v1.0.0`).

The workflow never embeds credentials in the image or logs; login uses the secrets above only.

## Build and run locally (Docker)

Build:

```bash
docker buildx build --load -t nevarix-agent:local --build-arg VERSION="$(git rev-parse --short HEAD)" .
```

Run the API server (requires a non-empty API token):

```bash
docker run --rm \
  -p 8080:8080 \
  -e NEVARIX_AGENT_API_TOKEN="your-secure-token" \
  -e AGENT_HTTP_ADDR=":8080" \
  nevarix-agent:local
```

Health check (replace the token):

```bash
curl -sS -H "Authorization: Bearer your-secure-token" http://127.0.0.1:8080/api/v1/health
```

Optional environment variables:

- `NEVARIX_AGENT_API_TOKEN` or `AGENT_API_TOKEN` — bearer token for `/api/v1/*` (required for `agent` mode).
- `AGENT_HTTP_ADDR` — listen address (default `:8080`).
- `NEVARIX_HUB_BASE_URL` — hub URL for the `monitor` subcommand when not using `runtime.json`.

The container runs as non-root user `appuser` (UID `65532`) and uses a writable `/home/.nevarix-server` for runtime state.

## Build and run locally (Go)

```bash
go build -o nevarix-agent ./cmd/agent-server
NEVARIX_AGENT_API_TOKEN="your-secure-token" ./nevarix-agent agent
```

## Pull from Docker Hub

After CI has published an image:

```bash
docker pull <DOCKER_USERNAME>/<repository-name>:latest
# or a specific version:
docker pull <DOCKER_USERNAME>/<repository-name>:v1.0.0
docker pull <DOCKER_USERNAME>/<repository-name>:sha-abc1234
```

Replace `<DOCKER_USERNAME>` and `<repository-name>` with your Docker Hub namespace and this repository name as shown on Docker Hub.

## Security notes

- Prefer Docker Hub **access tokens** with minimal scope, rotated periodically.
- The runtime image is `alpine:3.19` with a dedicated non-root user and a static Go binary (`CGO_ENABLED=0`) to keep the attack surface small.
- Do not commit `.env` files or tokens; use secrets in GitHub Actions and environment variables at deploy time.
