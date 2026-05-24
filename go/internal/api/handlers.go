package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/bukunya/intero-go/internal/satusehat"
)

type Handlers struct {
	ssClient *satusehat.Client
}

func NewHandlers(client *satusehat.Client) *Handlers {
	return &Handlers{ssClient: client}
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{Error: message})
}

func mapResourceToPerson(res *satusehat.Resource) PersonResponse {
	person := PersonResponse{
		Gender:    res.Gender,
		BirthDate: res.BirthDate,
	}

	for _, id := range res.Identifier {
		if strings.Contains(id.System, "id/nik") {
			person.NIK = id.Value
		} else if strings.Contains(id.System, "id/ihs-number") {
			person.IHSNumber = id.Value
		}
	}

	if len(res.Name) > 0 {
		person.Name = res.Name[0].Text
	}

	if len(res.Address) > 0 {
		if len(res.Address[0].Line) > 0 {
			person.Address = res.Address[0].Line[0]
		}
	}

	for _, tel := range res.Telecom {
		if tel.System == "phone" {
			person.Phone = tel.Value
			break
		}
	}
	if person.Phone == "" {
		for _, contact := range res.Contact {
			for _, tel := range contact.Telecom {
				if tel.System == "phone" {
					person.Phone = tel.Value
					break
				}
			}
			if person.Phone != "" {
				break
			}
		}
	}

	return person
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1]
	}
	return ""
}

func (h *Handlers) GetToken(w http.ResponseWriter, r *http.Request) {
	token, err := h.ssClient.GetToken()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get token from SatuSehat: "+err.Error())
		return
	}
	respondJSON(w, http.StatusOK, TokenResponse{Token: token})
}

func (h *Handlers) GetPatient(w http.ResponseWriter, r *http.Request) {
	nik := r.URL.Query().Get("nik")
	if nik == "" {
		respondError(w, http.StatusBadRequest, "nik query parameter is required")
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		respondError(w, http.StatusBadRequest, "name query parameter is required")
		return
	}

	bundle, err := h.ssClient.GetPatient(nik, name)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get patient from SatuSehat")
		return
	}

	if bundle.Total == 0 || len(bundle.Entry) == 0 {
		respondError(w, http.StatusNotFound, "patient not found")
		return
	}

	person := mapResourceToPerson(&bundle.Entry[0].Resource)
	respondJSON(w, http.StatusOK, person)
}

func (h *Handlers) GetPractitioner(w http.ResponseWriter, r *http.Request) {
	nik := r.URL.Query().Get("nik")
	if nik == "" {
		respondError(w, http.StatusBadRequest, "nik query parameter is required")
		return
	}

	bundle, err := h.ssClient.GetPractitionerByNIK(nik)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get practitioner from SatuSehat")
		return
	}

	if bundle.Total == 0 || len(bundle.Entry) == 0 {
		respondError(w, http.StatusNotFound, "practitioner not found")
		return
	}

	person := mapResourceToPerson(&bundle.Entry[0].Resource)
	respondJSON(w, http.StatusOK, person)
}

func (h *Handlers) CreatePatient(w http.ResponseWriter, r *http.Request) {
	var req CreatePatientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid json payload")
		return
	}

	fhirPayload := &satusehat.PatientPayload{
		ResourceType: "Patient",
		Meta: satusehat.Meta{
			Profile: []string{"https://fhir.kemkes.go.id/r4/StructureDefinition/Patient"},
		},
		Identifier: []satusehat.Identifier{
			{Use: "official", System: "https://fhir.kemkes.go.id/id/nik", Value: req.NIK},
		},
		Active: true,
		Name: []satusehat.Name{
			{Use: "official", Text: req.Name},
		},
		Gender:    req.Gender,
		BirthDate: req.BirthDate,
	}

	if req.Address != "" {
		addr := satusehat.Address{
			Use:        "home",
			Line:       []string{req.Address},
			City:       req.City,
			PostalCode: req.PostalCode,
			Country:    "ID",
		}

		if req.ProvinceCode != "" {
			ext := satusehat.Extension{
				URL: "https://fhir.kemkes.go.id/r4/StructureDefinition/administrativeCode",
				Extension: []satusehat.Extension{
					{URL: "province", ValueCode: req.ProvinceCode},
				},
			}
			if req.CityCode != "" {
				ext.Extension = append(ext.Extension, satusehat.Extension{URL: "city", ValueCode: req.CityCode})
			}
			if req.DistrictCode != "" {
				ext.Extension = append(ext.Extension, satusehat.Extension{URL: "district", ValueCode: req.DistrictCode})
			}
			if req.VillageCode != "" {
				ext.Extension = append(ext.Extension, satusehat.Extension{URL: "village", ValueCode: req.VillageCode})
			}
			if req.RT != "" {
				ext.Extension = append(ext.Extension, satusehat.Extension{URL: "rt", ValueCode: req.RT})
			}
			if req.RW != "" {
				ext.Extension = append(ext.Extension, satusehat.Extension{URL: "rw", ValueCode: req.RW})
			}

			addr.Extension = append(addr.Extension, ext)
		}

		fhirPayload.Address = append(fhirPayload.Address, addr)
	}

	if req.Phone != "" {
		fhirPayload.Telecom = append(fhirPayload.Telecom, satusehat.Telecom{
			System: "phone",
			Value:  req.Phone,
			Use:    "mobile",
		})
	}

	res, err := h.ssClient.CreatePatient(fhirPayload)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create patient: "+err.Error())
		return
	}

	person := mapResourceToPerson(res)
	fmt.Println(person)
	respondJSON(w, http.StatusCreated, person)
}
