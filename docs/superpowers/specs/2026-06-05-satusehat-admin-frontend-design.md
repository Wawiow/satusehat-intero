# SatuSehat Intero Admin Frontend Design

Date: 2026-06-05

## Goal

Build an internal admin frontend for the existing SatuSehat Intero Go API. The frontend must cover the API surface currently implemented in the backend and make operator workflows practical: search, create, list local records, inspect errors, and update encounter status.

## Scope

In scope:

- Convert the repository into a minimal monorepo without moving the current Go backend.
- Add a Next.js app with Ant Design at `apps/web`.
- Implement internal admin screens only. No public page and no frontend login.
- Implement UI coverage for token, patients, practitioners, locations, and encounters.
- Add one backend route registration for the existing `GetAllLocalLocations` handler: `GET /api/local/locations`.
- Use a configurable API base URL, defaulting to `http://localhost:8083/api`.

Out of scope:

- Backend authentication or user management.
- Moving Go backend files into `apps/api`.
- Changing SatuSehat request semantics beyond the missing local locations route.
- Replacing SQLite schema or adding migrations beyond what already exists.

## Architecture

The repository will become a minimal monorepo:

- Root remains the Go backend build root.
- `apps/web` contains the Next.js admin app.
- Root `package.json` and `pnpm-workspace.yaml` manage frontend workspace scripts.
- Frontend calls backend directly through `NEXT_PUBLIC_API_BASE_URL`.

This preserves the existing Go module path and build command while adding a frontend workspace cleanly.

## Frontend Structure

The Next.js app will use App Router, TypeScript, Ant Design, and client-side data fetching.

Primary navigation:

- `Overview`: backend connectivity, token request, and local dataset counts.
- `Patients`: local patient table, search by NIK/name, create patient form.
- `Practitioners`: local practitioner table and search by NIK, ID, name, gender, birthdate, page, and limit.
- `Locations`: local locations table, search by ID or identifier, create location form.
- `Encounters`: local encounters table, get by ID, create encounter, update status.
- `API Tools`: raw request runner for debugging supported endpoints.

## Data Flow

The frontend will use a typed API client:

- `request<T>()` wraps `fetch`, parses JSON, captures status code, and throws a structured API error.
- API modules expose functions for each backend route.
- Pages keep request state local with React hooks and Ant Design feedback components.
- Create/update forms refresh relevant local lists after success.

Supported backend routes:

- `POST /api/token`
- `GET /api/patients?nik=&name=`
- `GET /api/local/patients`
- `POST /api/patients`
- `GET /api/practitioners?...`
- `GET /api/local/practitioners`
- `GET /api/locations?...`
- `GET /api/local/locations`
- `POST /api/locations`
- `GET /api/encounters/{id}`
- `GET /api/local/encounters`
- `POST /api/encounters`
- `PUT /api/encounters/{id}`

## UI Behavior

The UI is a dense internal operations dashboard, not a marketing page.

- Use Ant Design `Layout`, `Menu`, `Table`, `Form`, `Input`, `Select`, `DatePicker`, `Button`, `Tabs`, `Alert`, `Modal`, `Descriptions`, `Statistic`, and `Tag`.
- Keep tables central for repeated data. Do not convert records into card grids.
- Keep forms compact and grouped by resource.
- Use visible status messages for backend errors, validation failures, empty responses, and SatuSehat failures.
- Display raw JSON response panels for create/search actions so operators can verify exact API output.

## Error Handling

Each API interaction must surface:

- Loading state.
- HTTP status code when available.
- Backend JSON error message when available.
- Network error message when backend is offline or CORS/proxy fails.
- Empty state when arrays return no records.

The frontend must not hide unstable backend behavior behind generic success text.

## Backend Adjustment

Add this route to `internal/server/server.go`:

```go
mux.HandleFunc("GET /api/local/locations", handlers.GetAllLocalLocations)
```

No other backend behavior changes are required for the frontend scope.

## Testing And Verification

Backend verification:

- Run `go test ./...`.

Frontend verification:

- Install workspace dependencies.
- Run lint/typecheck/build through available package scripts.
- Start Go backend on `:8083`.
- Start Next.js dev server.
- Open the app and verify core flows: token, local lists, patient search/create form validation, practitioner search, location search/create, encounter get/create/status update UI.

Visual verification:

- Use browser testing after implementation to check desktop and mobile layouts.
- Ensure tables, forms, and buttons do not overflow on mobile.
- Confirm visible UI matches the internal admin design: restrained, table-first, Ant Design based, no public landing page.

## Acceptance Criteria

- Repository has a working `apps/web` Next.js app managed through root monorepo scripts.
- Go backend still builds and tests from the root.
- Frontend covers every backend route listed in this spec.
- Missing local locations route is registered.
- API base URL can be changed with `NEXT_PUBLIC_API_BASE_URL`.
- All main screens show loading, error, empty, and success states.
- Build/test verification commands are run and reported.
