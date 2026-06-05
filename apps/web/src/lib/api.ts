export const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8083/api";

export class ApiError extends Error {
  constructor(message: string, public status?: number, public payload?: unknown) {
    super(message);
    this.name = "ApiError";
  }
}

export type PersonResponse = {
  id: string;
  nik: string;
  ihs_number?: string;
  name: string;
  gender?: string;
  birth_date?: string;
  phone?: string;
  address?: string;
};

export type CreatePatientRequest = {
  nik: string;
  name: string;
  gender: string;
  birth_date: string;
  phone?: string;
  address?: string;
  city?: string;
  province_code?: string;
  city_code?: string;
  district_code?: string;
  village_code?: string;
  rt?: string;
  rw?: string;
  postal_code?: string;
};

export type TokenResponse = {
  token: string;
};

export type LocationResponse = {
  id: string;
  identifier_value: string;
  name: string;
  description: string;
  phone?: string;
};

export type CreateLocationRequest = {
  identifier_value: string;
  name: string;
  description: string;
  phone?: string;
};

export type EncounterResponse = {
  id: string;
  identifier_value: string;
  status: string;
  subject_id: string;
  location_id: string;
  start_time: string;
};

export type CreateEncounterRequest = {
  identifier_value: string;
  subject_id: string;
  location_id: string;
  practitioner_id: string;
  start_time: string;
};

export type UpdateEncounterRequest = {
  status: string;
};

type QueryValue = string | number | undefined | null;

function cleanBody<T extends Record<string, unknown>>(body: T): T {
  return Object.fromEntries(
    Object.entries(body).filter(([, value]) => value !== undefined && value !== null && value !== ""),
  ) as T;
}

export function buildQuery(params: Record<string, QueryValue>) {
  const query = new URLSearchParams();

  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== "") {
      query.set(key, String(value));
    }
  });

  const value = query.toString();
  return value ? `?${value}` : "";
}

export async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...init?.headers,
    },
  });

  const contentType = response.headers.get("content-type") ?? "";
  const payload = contentType.includes("application/json") ? await response.json() : await response.text();

  if (!response.ok) {
    const message =
      typeof payload === "object" && payload !== null && "error" in payload
        ? String((payload as { error: unknown }).error)
        : `Request failed with status ${response.status}`;
    throw new ApiError(message, response.status, payload);
  }

  return payload as T;
}

export function getToken() {
  return request<TokenResponse>("/token", { method: "POST" });
}

export function getLocalPatients() {
  return request<PersonResponse[]>("/local/patients");
}

export function searchPatient(params: { nik: string; name?: string }) {
  return request<PersonResponse>(`/patients${buildQuery(params)}`);
}

export function createPatient(body: CreatePatientRequest) {
  return request<PersonResponse>("/patients", { method: "POST", body: JSON.stringify(cleanBody(body)) });
}

export function getLocalPractitioners() {
  return request<PersonResponse[]>("/local/practitioners");
}

export function searchPractitioners(params: {
  id?: string;
  nik?: string;
  name?: string;
  gender?: string;
  birthdate?: string;
  page?: number;
  limit?: number;
}) {
  return request<PersonResponse[]>(`/practitioners${buildQuery(params)}`);
}

export function getLocalLocations() {
  return request<LocationResponse[]>("/local/locations");
}

export function searchLocations(params: { id?: string; identifier?: string; page?: number; limit?: number }) {
  return request<LocationResponse[]>(`/locations${buildQuery(params)}`);
}

export function createLocation(body: CreateLocationRequest) {
  return request<LocationResponse>("/locations", { method: "POST", body: JSON.stringify(cleanBody(body)) });
}

export function getLocalEncounters() {
  return request<EncounterResponse[]>("/local/encounters");
}

export function getEncounterById(id: string) {
  return request<EncounterResponse>(`/encounters/${encodeURIComponent(id)}`);
}

export function createEncounter(body: CreateEncounterRequest) {
  return request<EncounterResponse>("/encounters", { method: "POST", body: JSON.stringify(cleanBody(body)) });
}

export function updateEncounterStatus(id: string, body: UpdateEncounterRequest) {
  return request<EncounterResponse>(`/encounters/${encodeURIComponent(id)}`, {
    method: "PUT",
    body: JSON.stringify(cleanBody(body)),
  });
}

export function formatApiError(error: unknown) {
  if (error instanceof ApiError) {
    return error.status ? `${error.status}: ${error.message}` : error.message;
  }
  if (error instanceof Error) {
    return error.message;
  }
  return "Unknown request error";
}
