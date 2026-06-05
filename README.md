# SatuSehat Intero

SatuSehat Intero is a hospital operations dashboard and Go API wrapper for integrating local hospital workflows with Indonesia's SatuSehat services. The project combines a backend service for FHIR-oriented endpoints with a Next.js dashboard for operational teams.

## Capabilities

- Monitor local patient, practitioner, location, and encounter data.
- Search and create supported SatuSehat resources through backend endpoints.
- Validate SatuSehat OAuth token access from the dashboard.
- Run raw backend endpoint checks from the integration tools screen.
- Store local operational data in SQLite during development.

## Tech Stack

| Layer | Technology |
| --- | --- |
| Backend | Go, net/http, SQLite |
| Frontend | Next.js 15, React 19, Ant Design 6 |
| Package manager | pnpm 10 |
| API docs | Swagger UI |

## Project Structure

```text
.
├── apps/web                 # Next.js hospital dashboard
├── cmd/myapp                # Go backend entrypoint
├── internal/api             # HTTP handlers
├── internal/config          # Environment configuration
├── internal/database        # SQLite setup
├── internal/satusehat       # SatuSehat client and FHIR models
├── docs                     # Generated Swagger package
├── http                     # Request samples
└── app.db                   # Local SQLite database
```

## Prerequisites

- Go 1.24 or newer
- Node.js compatible with Next.js 15
- pnpm 10.11.0 or newer
- SatuSehat sandbox or production credentials

## Environment

Create or update the root `.env` file for backend credentials:

```env
AUTH_URL=https://api-satusehat-stg.dto.kemkes.go.id/oauth2/v1
BASE_URL=https://api-satusehat-stg.dto.kemkes.go.id/fhir-r4/v1
CONSENT_URL=https://api-satusehat-stg.dto.kemkes.go.id/consent/v1
CLIENT_ID=your-client-id
CLIENT_SECRET=your-client-secret
ORG_ID=your-organization-id
```

Create the frontend environment file when the API base URL differs from the default:

```powershell
Copy-Item apps\web\.env.example apps\web\.env.local
```

Default frontend API URL:

```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8083/api
```

## Local Development

Install frontend dependencies:

```powershell
pnpm install
```

Run the backend API:

```powershell
go run ./cmd/myapp
```

Backend runs at:

- API: [http://localhost:8083](http://localhost:8083)
- Swagger: [http://localhost:8083/swagger/](http://localhost:8083/swagger/)

Run the dashboard in another terminal:

```powershell
pnpm web:dev
```

Dashboard runs at:

- Web: [http://localhost:3000](http://localhost:3000)

## Build

Build the backend binary for Windows:

```powershell
$env:GOOS="windows"
$env:GOARCH="amd64"
go build -o myapp.exe ./cmd/myapp
```

Build the backend binary for Linux:

```powershell
go build -o myapp-linux ./cmd/myapp
```

Build the frontend:

```powershell
pnpm web:build
```

Run frontend quality checks:

```powershell
pnpm web:typecheck
pnpm web:lint
```

## API Surface

Main backend endpoints:

| Method | Endpoint | Purpose |
| --- | --- | --- |
| `POST` | `/api/token` | Request SatuSehat access token |
| `GET` | `/api/patients` | Search SatuSehat patient data |
| `POST` | `/api/patients` | Create patient resource |
| `GET` | `/api/local/patients` | Read local patient cache |
| `GET` | `/api/practitioners` | Search SatuSehat practitioner data |
| `GET` | `/api/local/practitioners` | Read local practitioner cache |
| `GET` | `/api/locations` | Search SatuSehat location data |
| `POST` | `/api/locations` | Create location resource |
| `GET` | `/api/local/locations` | Read local location cache |
| `GET` | `/api/encounters/{id}` | Read encounter detail |
| `POST` | `/api/encounters` | Create encounter resource |
| `PUT` | `/api/encounters/{id}` | Update encounter status |
| `GET` | `/api/local/encounters` | Read local encounter cache |

## Operational Notes

- Keep SatuSehat credentials out of source control.
- Use staging endpoints for development and operator training.
- Verify `ORG_ID` before creating resources, because several payloads reference the configured organization.
- The dashboard is designed for internal hospital operations teams, not for public patient access.
