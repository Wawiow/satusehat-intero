# SatuSehat Admin Frontend Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Next.js + Ant Design internal admin app in a minimal monorepo that covers every current SatuSehat Intero API route.

**Architecture:** Keep the Go backend at the repository root and add `D:\satusehat-intero\apps\web` as the frontend workspace. The frontend uses client-side React state, a typed fetch wrapper, Ant Design layout/components, and direct calls to `NEXT_PUBLIC_API_BASE_URL` with `http://localhost:8083/api` as the fallback. The only backend behavior change is registering the already implemented local locations handler.

**Tech Stack:** Go `net/http`, SQLite backend already present; Next.js App Router, React, TypeScript, Ant Design, pnpm workspace.

---

## File Structure

- Modify `D:\satusehat-intero\internal\server\server.go`: register `GET /api/local/locations`.
- Create `D:\satusehat-intero\package.json`: root monorepo scripts for web workspace.
- Create `D:\satusehat-intero\pnpm-workspace.yaml`: include `apps/*`.
- Create `D:\satusehat-intero\apps\web\package.json`: Next.js dependencies and scripts.
- Create `D:\satusehat-intero\apps\web\next.config.ts`: Next config.
- Create `D:\satusehat-intero\apps\web\tsconfig.json`: TypeScript config.
- Create `D:\satusehat-intero\apps\web\next-env.d.ts`: Next env typing.
- Create `D:\satusehat-intero\apps\web\src\app\layout.tsx`: root metadata and global CSS import.
- Create `D:\satusehat-intero\apps\web\src\app\page.tsx`: app entry rendering admin dashboard.
- Create `D:\satusehat-intero\apps\web\src\app\globals.css`: Ant Design friendly global styling and responsive layout fixes.
- Create `D:\satusehat-intero\apps\web\src\lib\api.ts`: typed API client, models, and route functions.
- Create `D:\satusehat-intero\apps\web\src\components\AdminShell.tsx`: sidebar/header layout and screen switching.
- Create `D:\satusehat-intero\apps\web\src\components\JsonPanel.tsx`: formatted JSON/error display.
- Create `D:\satusehat-intero\apps\web\src\components\DataToolbar.tsx`: reusable compact toolbar for reload/search actions.
- Create `D:\satusehat-intero\apps\web\src\features\overview\OverviewScreen.tsx`: connectivity, token, counts.
- Create `D:\satusehat-intero\apps\web\src\features\patients\PatientsScreen.tsx`: local list, search, create.
- Create `D:\satusehat-intero\apps\web\src\features\practitioners\PractitionersScreen.tsx`: local list and search.
- Create `D:\satusehat-intero\apps\web\src\features\locations\LocationsScreen.tsx`: local list, search, create.
- Create `D:\satusehat-intero\apps\web\src\features\encounters\EncountersScreen.tsx`: local list, get, create, update status.
- Create `D:\satusehat-intero\apps\web\src\features\tools\ApiToolsScreen.tsx`: raw supported endpoint runner.
- Create `D:\satusehat-intero\apps\web\.env.example`: API base URL example.

---

### Task 1: Backend Route Coverage

**Files:**
- Modify: `D:\satusehat-intero\internal\server\server.go`

- [ ] **Step 1: Confirm local locations route is absent**

Run: `Select-String -Path 'D:\satusehat-intero\internal\server\server.go' -Pattern 'local/locations'`

Expected: no output.

- [ ] **Step 2: Add route registration**

In `D:\satusehat-intero\internal\server\server.go`, add this line near the locations routes:

```go
mux.HandleFunc("GET /api/local/locations", handlers.GetAllLocalLocations)
```

Expected surrounding block:

```go
// OK
mux.HandleFunc("GET /api/locations", handlers.GetLocations)
mux.HandleFunc("GET /api/local/locations", handlers.GetAllLocalLocations)
mux.HandleFunc("POST /api/locations", handlers.CreateLocation)
```

- [ ] **Step 3: Verify backend tests compile**

Run: `go test ./...`

Expected: all packages pass or report `[no test files]`, with no compile errors.

- [ ] **Step 4: Commit backend route**

Run:

```powershell
git add D:\satusehat-intero\internal\server\server.go
git commit -m "fix: expose local locations endpoint"
```

Expected: one commit containing only `internal/server/server.go`.

---

### Task 2: Monorepo And Next.js Scaffold

**Files:**
- Create: `D:\satusehat-intero\package.json`
- Create: `D:\satusehat-intero\pnpm-workspace.yaml`
- Create: `D:\satusehat-intero\apps\web\package.json`
- Create: `D:\satusehat-intero\apps\web\next.config.ts`
- Create: `D:\satusehat-intero\apps\web\tsconfig.json`
- Create: `D:\satusehat-intero\apps\web\next-env.d.ts`
- Create: `D:\satusehat-intero\apps\web\.env.example`
- Create: `D:\satusehat-intero\apps\web\src\app\layout.tsx`
- Create: `D:\satusehat-intero\apps\web\src\app\page.tsx`
- Create: `D:\satusehat-intero\apps\web\src\app\globals.css`

- [ ] **Step 1: Create root workspace files**

Create `D:\satusehat-intero\package.json`:

```json
{
  "name": "satusehat-intero",
  "private": true,
  "scripts": {
    "web:dev": "pnpm --filter @satusehat-intero/web dev",
    "web:build": "pnpm --filter @satusehat-intero/web build",
    "web:lint": "pnpm --filter @satusehat-intero/web lint",
    "web:typecheck": "pnpm --filter @satusehat-intero/web typecheck"
  },
  "packageManager": "pnpm@10.11.0"
}
```

Create `D:\satusehat-intero\pnpm-workspace.yaml`:

```yaml
packages:
  - "apps/*"
```

- [ ] **Step 2: Create web package files**

Create `D:\satusehat-intero\apps\web\package.json`:

```json
{
  "name": "@satusehat-intero/web",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "next dev --hostname 0.0.0.0 --port 3000",
    "build": "next build",
    "lint": "next lint",
    "typecheck": "tsc --noEmit"
  },
  "dependencies": {
    "@ant-design/icons": "^5.6.1",
    "antd": "^5.26.0",
    "dayjs": "^1.11.13",
    "next": "^15.3.3",
    "react": "^19.1.0",
    "react-dom": "^19.1.0"
  },
  "devDependencies": {
    "@types/node": "^22.15.29",
    "@types/react": "^19.1.6",
    "@types/react-dom": "^19.1.5",
    "eslint": "^9.28.0",
    "eslint-config-next": "^15.3.3",
    "typescript": "^5.8.3"
  }
}
```

Create `D:\satusehat-intero\apps\web\next.config.ts`:

```ts
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  reactStrictMode: true,
};

export default nextConfig;
```

Create `D:\satusehat-intero\apps\web\tsconfig.json`:

```json
{
  "compilerOptions": {
    "target": "ES2017",
    "lib": ["dom", "dom.iterable", "esnext"],
    "allowJs": false,
    "skipLibCheck": true,
    "strict": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "preserve",
    "incremental": true,
    "plugins": [{ "name": "next" }],
    "paths": { "@/*": ["./src/*"] }
  },
  "include": ["next-env.d.ts", "**/*.ts", "**/*.tsx", ".next/types/**/*.ts"],
  "exclude": ["node_modules"]
}
```

Create `D:\satusehat-intero\apps\web\next-env.d.ts`:

```ts
/// <reference types="next" />
/// <reference types="next/image-types/global" />
```

Create `D:\satusehat-intero\apps\web\.env.example`:

```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8083/api
```

- [ ] **Step 3: Create minimal app entry**

Create `D:\satusehat-intero\apps\web\src\app\layout.tsx`:

```tsx
import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "SatuSehat Intero Admin",
  description: "Internal admin console for SatuSehat Intero API",
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
```

Create `D:\satusehat-intero\apps\web\src\app\page.tsx`:

```tsx
"use client";

import { Typography } from "antd";

export default function Home() {
  return <Typography.Title level={1}>SatuSehat Intero Admin</Typography.Title>;
}
```

Create `D:\satusehat-intero\apps\web\src\app\globals.css`:

```css
* {
  box-sizing: border-box;
}

html,
body {
  min-height: 100%;
  margin: 0;
  background: #f5f7fb;
  color: #172033;
}

body {
  font-family: Arial, Helvetica, sans-serif;
}

button,
input,
textarea,
select {
  font: inherit;
}
```

- [ ] **Step 4: Install dependencies and verify scaffold**

Run: `pnpm install`

Run: `pnpm web:typecheck`

Expected: TypeScript passes.

Run: `pnpm web:build`

Expected: Next.js production build succeeds.

- [ ] **Step 5: Commit scaffold**

Run:

```powershell
git add D:\satusehat-intero\package.json D:\satusehat-intero\pnpm-workspace.yaml D:\satusehat-intero\pnpm-lock.yaml D:\satusehat-intero\apps\web
git commit -m "feat: scaffold next admin workspace"
```

Expected: one commit with frontend scaffold and lockfile.

---

### Task 3: Typed API Client And Shared Components

**Files:**
- Create: `D:\satusehat-intero\apps\web\src\lib\api.ts`
- Create: `D:\satusehat-intero\apps\web\src\components\JsonPanel.tsx`
- Create: `D:\satusehat-intero\apps\web\src\components\DataToolbar.tsx`

- [ ] **Step 1: Create typed API client**

Create `D:\satusehat-intero\apps\web\src\lib\api.ts` with exported models matching `D:\satusehat-intero\internal\api\models.go` and functions for all supported routes. Include `ApiError`, `request<T>()`, `getToken`, `getLocalPatients`, `searchPatient`, `createPatient`, `getLocalPractitioners`, `searchPractitioners`, `getLocalLocations`, `searchLocations`, `createLocation`, `getLocalEncounters`, `getEncounterById`, `createEncounter`, and `updateEncounterStatus`.

Required behavior:

```ts
export const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8083/api";

export class ApiError extends Error {
  constructor(message: string, public status?: number, public payload?: unknown) {
    super(message);
    this.name = "ApiError";
  }
}
```

`request<T>()` must parse JSON when available, return typed data on `response.ok`, and throw `ApiError` with backend `error` payload text when not ok.

- [ ] **Step 2: Create JSON panel**

Create `D:\satusehat-intero\apps\web\src\components\JsonPanel.tsx`:

```tsx
import { Alert, Typography } from "antd";

type JsonPanelProps = {
  title: string;
  data: unknown;
  error?: string;
};

export function JsonPanel({ title, data, error }: JsonPanelProps) {
  return (
    <section className="json-panel">
      <Typography.Title level={5}>{title}</Typography.Title>
      {error ? <Alert type="error" showIcon message={error} /> : null}
      <pre>{JSON.stringify(data ?? null, null, 2)}</pre>
    </section>
  );
}
```

- [ ] **Step 3: Create toolbar component**

Create `D:\satusehat-intero\apps\web\src\components\DataToolbar.tsx`:

```tsx
import { Button, Space, Typography } from "antd";
import { ReloadOutlined } from "@ant-design/icons";

type DataToolbarProps = {
  title: string;
  description?: string;
  loading?: boolean;
  onReload?: () => void;
  extra?: React.ReactNode;
};

export function DataToolbar({ title, description, loading, onReload, extra }: DataToolbarProps) {
  return (
    <div className="data-toolbar">
      <div>
        <Typography.Title level={3}>{title}</Typography.Title>
        {description ? <Typography.Text type="secondary">{description}</Typography.Text> : null}
      </div>
      <Space wrap>
        {extra}
        {onReload ? (
          <Button icon={<ReloadOutlined />} loading={loading} onClick={onReload}>
            Reload
          </Button>
        ) : null}
      </Space>
    </div>
  );
}
```

- [ ] **Step 4: Typecheck shared code**

Run: `pnpm web:typecheck`

Expected: TypeScript passes.

- [ ] **Step 5: Commit API client and shared components**

Run:

```powershell
git add D:\satusehat-intero\apps\web\src\lib\api.ts D:\satusehat-intero\apps\web\src\components
git commit -m "feat: add typed admin api client"
```

Expected: one commit with shared API/UI foundation.

---

### Task 4: Admin Shell And Global Styling

**Files:**
- Modify: `D:\satusehat-intero\apps\web\src\app\page.tsx`
- Modify: `D:\satusehat-intero\apps\web\src\app\globals.css`
- Create: `D:\satusehat-intero\apps\web\src\components\AdminShell.tsx`

- [ ] **Step 1: Create admin shell**

Create `D:\satusehat-intero\apps\web\src\components\AdminShell.tsx` with Ant Design `Layout`, `Sider`, `Header`, `Content`, and menu keys: `overview`, `patients`, `practitioners`, `locations`, `encounters`, `tools`. Use local `useState` for active screen and render placeholder components passed by import names from feature screens created in later tasks.

Required first render fallback before feature screens exist:

```tsx
function EmptyScreen({ name }: { name: string }) {
  return <Typography.Title level={2}>{name}</Typography.Title>;
}
```

- [ ] **Step 2: Wire page to shell**

Replace `D:\satusehat-intero\apps\web\src\app\page.tsx` with:

```tsx
"use client";

import { AdminShell } from "@/components/AdminShell";

export default function Home() {
  return <AdminShell />;
}
```

- [ ] **Step 3: Extend global CSS**

Add classes for `.admin-layout`, `.admin-sider`, `.admin-header`, `.admin-content`, `.screen-stack`, `.data-toolbar`, `.json-panel`, `.json-panel pre`, `.table-wrap`, and responsive breakpoint `@media (max-width: 760px)` so sider becomes top nav style and tables scroll horizontally.

- [ ] **Step 4: Verify shell build**

Run: `pnpm web:typecheck`

Expected: TypeScript passes.

Run: `pnpm web:build`

Expected: Next.js build succeeds.

- [ ] **Step 5: Commit shell**

Run:

```powershell
git add D:\satusehat-intero\apps\web\src\app D:\satusehat-intero\apps\web\src\components\AdminShell.tsx
git commit -m "feat: add internal admin shell"
```

Expected: one commit with layout and responsive CSS.

---

### Task 5: Overview, Patients, And Practitioners Screens

**Files:**
- Create: `D:\satusehat-intero\apps\web\src\features\overview\OverviewScreen.tsx`
- Create: `D:\satusehat-intero\apps\web\src\features\patients\PatientsScreen.tsx`
- Create: `D:\satusehat-intero\apps\web\src\features\practitioners\PractitionersScreen.tsx`
- Modify: `D:\satusehat-intero\apps\web\src\components\AdminShell.tsx`

- [ ] **Step 1: Build overview screen**

Create `OverviewScreen.tsx` with four actions on load or button click: `getLocalPatients`, `getLocalPractitioners`, `getLocalLocations`, `getLocalEncounters`. Show Ant Design `Statistic` counts and `POST /api/token` button that renders token result in `JsonPanel`.

- [ ] **Step 2: Build patients screen**

Create `PatientsScreen.tsx` with:

- Local patients table using columns `nik`, `ihs_number`, `id`, `name`, `gender`, `birth_date`, `phone`, `address`.
- Search form fields `nik` and `name`; call `searchPatient({ nik, name })`.
- Create form fields matching `CreatePatientRequest`: `nik`, `name`, `gender`, `birth_date`, `phone`, `address`, `city`, `province_code`, `city_code`, `district_code`, `village_code`, `rt`, `rw`, `postal_code`.
- Required frontend validation for `nik`, `name`, `gender`, and `birth_date`.
- Raw response panel for search/create result.

- [ ] **Step 3: Build practitioners screen**

Create `PractitionersScreen.tsx` with:

- Local practitioners table using same person columns as patients.
- Search form fields `id`, `nik`, `name`, `gender`, `birthdate`, `page`, `limit`.
- Default `page=1`, `limit=10`.
- Raw response panel for search result.

- [ ] **Step 4: Wire screens into admin shell**

Import and render `OverviewScreen`, `PatientsScreen`, and `PractitionersScreen` for matching menu keys. Keep placeholders for locations, encounters, and tools until later tasks.

- [ ] **Step 5: Verify first resource screens**

Run: `pnpm web:typecheck`

Expected: TypeScript passes.

Run: `pnpm web:build`

Expected: Next.js build succeeds.

- [ ] **Step 6: Commit first resource screens**

Run:

```powershell
git add D:\satusehat-intero\apps\web\src\features\overview D:\satusehat-intero\apps\web\src\features\patients D:\satusehat-intero\apps\web\src\features\practitioners D:\satusehat-intero\apps\web\src\components\AdminShell.tsx
git commit -m "feat: add patient and practitioner admin screens"
```

Expected: one commit with overview, patients, practitioners.

---

### Task 6: Locations And Encounters Screens

**Files:**
- Create: `D:\satusehat-intero\apps\web\src\features\locations\LocationsScreen.tsx`
- Create: `D:\satusehat-intero\apps\web\src\features\encounters\EncountersScreen.tsx`
- Modify: `D:\satusehat-intero\apps\web\src\components\AdminShell.tsx`

- [ ] **Step 1: Build locations screen**

Create `LocationsScreen.tsx` with:

- Local locations table columns `id`, `identifier_value`, `name`, `description`, `phone`.
- Search form fields `id`, `identifier`, `page`, `limit`, default `page=1`, `limit=10`.
- Create form fields `identifier_value`, `name`, `description`, `phone`; require first three.
- Raw response panel for search/create.

- [ ] **Step 2: Build encounters screen**

Create `EncountersScreen.tsx` with:

- Local encounters table columns `id`, `identifier_value`, `status`, `subject_id`, `location_id`, `start_time`.
- Get-by-ID form field `id`; call `getEncounterById`.
- Create form fields `identifier_value`, `subject_id`, `location_id`, `practitioner_id`, `start_time`; require all.
- Update status form fields `id`, `status`; status options `arrived`, `in-progress`, `finished`, `cancelled`.
- Raw response panel for get/create/update.

- [ ] **Step 3: Wire screens into admin shell**

Import and render `LocationsScreen` and `EncountersScreen`. Keep `API Tools` placeholder until Task 7.

- [ ] **Step 4: Verify second resource screens**

Run: `pnpm web:typecheck`

Expected: TypeScript passes.

Run: `pnpm web:build`

Expected: Next.js build succeeds.

- [ ] **Step 5: Commit second resource screens**

Run:

```powershell
git add D:\satusehat-intero\apps\web\src\features\locations D:\satusehat-intero\apps\web\src\features\encounters D:\satusehat-intero\apps\web\src\components\AdminShell.tsx
git commit -m "feat: add locations and encounters admin screens"
```

Expected: one commit with locations and encounters.

---

### Task 7: API Tools Screen

**Files:**
- Create: `D:\satusehat-intero\apps\web\src\features\tools\ApiToolsScreen.tsx`
- Modify: `D:\satusehat-intero\apps\web\src\components\AdminShell.tsx`

- [ ] **Step 1: Build supported endpoint runner**

Create `ApiToolsScreen.tsx` with:

- Method select: `GET`, `POST`, `PUT`.
- Endpoint select with these exact values: `/token`, `/local/patients`, `/patients`, `/local/practitioners`, `/practitioners`, `/local/locations`, `/locations`, `/local/encounters`, `/encounters`, `/encounters/{id}`.
- Path ID input visible for `/encounters/{id}`.
- Query string input for GET requests, entered without `?`, appended to URL.
- JSON body text area for POST/PUT requests.
- Execute button using `fetch` against `API_BASE_URL`.
- Response panel showing status code and parsed JSON or text.

- [ ] **Step 2: Wire tools screen into admin shell**

Import and render `ApiToolsScreen` for menu key `tools`.

- [ ] **Step 3: Verify tools screen**

Run: `pnpm web:typecheck`

Expected: TypeScript passes.

Run: `pnpm web:build`

Expected: Next.js build succeeds.

- [ ] **Step 4: Commit tools screen**

Run:

```powershell
git add D:\satusehat-intero\apps\web\src\features\tools D:\satusehat-intero\apps\web\src\components\AdminShell.tsx
git commit -m "feat: add admin api tools screen"
```

Expected: one commit with API tools.

---

### Task 8: End-To-End Verification And Handoff

**Files:**
- Modify only if verification finds defects in files created by earlier tasks.

- [ ] **Step 1: Run backend verification**

Run: `go test ./...`

Expected: no compile errors; tests pass or packages report `[no test files]`.

- [ ] **Step 2: Run frontend verification**

Run: `pnpm web:typecheck`

Expected: TypeScript passes.

Run: `pnpm web:build`

Expected: Next.js build succeeds.

- [ ] **Step 3: Start backend**

Run: `go run .\cmd\myapp`

Expected: log includes `Server listening on :8083`.

- [ ] **Step 4: Start frontend**

Run: `pnpm web:dev`

Expected: Next.js listens on `http://localhost:3000` or next available port.

- [ ] **Step 5: Browser verify core UI**

Open `http://localhost:3000` with Browser plugin and verify:

- Sidebar has `Overview`, `Patients`, `Practitioners`, `Locations`, `Encounters`, `API Tools`.
- Overview can request token and local counts.
- Patients screen can load local patients and show validation errors on missing required create fields.
- Practitioners screen can submit a search and show raw response or backend error.
- Locations screen can load local locations through `GET /api/local/locations`.
- Encounters screen can load local encounters and show update status form.
- API Tools can call `GET /local/patients` and show status plus response body.

- [ ] **Step 6: Browser verify responsive layout**

Use desktop width and mobile width around `390px`. Verify no primary controls overflow, tables scroll horizontally, and page remains usable.

- [ ] **Step 7: Final status**

Run: `git status --short`

Expected: only intentional changes remain. If `D:\satusehat-intero\myapp.exe` is still modified from before, report it as pre-existing/unrelated and do not revert it.

---

## Self-Review Notes

- Spec coverage: monorepo scaffold is Task 2; backend missing route is Task 1; typed API/data flow is Task 3; shell/navigation is Task 4; Overview/Patients/Practitioners are Task 5; Locations/Encounters are Task 6; API Tools is Task 7; verification is Task 8.
- Placeholder scan: no unresolved placeholder text and no omitted route group.
- Type consistency: route names and payload fields match `D:\satusehat-intero\internal\api\models.go`; local locations route matches existing `GetAllLocalLocations` handler.
