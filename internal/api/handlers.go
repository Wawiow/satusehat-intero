package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bukunya/intero-go/internal/satusehat"
)

type Handlers struct {
	ssClient *satusehat.Client
	db       *sql.DB
}

func NewHandlers(client *satusehat.Client, db *sql.DB) *Handlers {
	return &Handlers{ssClient: client, db: db}
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
		ID:        res.ID,
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

// GetToken godoc
// @Summary      Get SatuSehat OAuth Token
// @Description  Get a fresh OAuth2 bearer token from SatuSehat
// @Tags         Auth
// @Produce      json
// @Success      200  {object}  TokenResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/token [post]
func (h *Handlers) GetToken(w http.ResponseWriter, r *http.Request) {
	token, err := h.ssClient.GetToken()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get token from SatuSehat: "+err.Error())
		return
	}
	respondJSON(w, http.StatusOK, TokenResponse{Token: token})
}

// GetAllLocalPatients godoc
// @Summary      Get All Local Patients
// @Description  Retrieve all patients stored in the local SQLite database
// @Tags         Patients
// @Produce      json
// @Success      200  {array}   PersonResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/local/patients [get]
func (h *Handlers) GetAllLocalPatients(w http.ResponseWriter, r *http.Request) {
	patients, err := getAllPatientsLocal(h.db)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch local patients")
		return
	}
	respondJSON(w, http.StatusOK, patients)
}

// GetAllLocalLocations godoc
// @Summary      Get All Local Locations
// @Description  Retrieve all locations stored in the local SQLite database
// @Tags         Locations
// @Produce      json
// @Success      200  {array}   LocationResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/local/locations [get]
func (h *Handlers) GetAllLocalLocations(w http.ResponseWriter, r *http.Request) {
	locs, err := getAllLocationsLocal(h.db)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch local locations")
		return
	}
	respondJSON(w, http.StatusOK, locs)
}

// GetAllLocalEncounters godoc
// @Summary      Get All Local Encounters
// @Description  Retrieve all encounters stored in the local SQLite database
// @Tags         Encounters
// @Produce      json
// @Success      200  {array}   EncounterResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/local/encounters [get]
func (h *Handlers) GetAllLocalEncounters(w http.ResponseWriter, r *http.Request) {
	encs, err := getAllEncountersLocal(h.db)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch local encounters")
		return
	}
	respondJSON(w, http.StatusOK, encs)
}

// GetPatient godoc
// @Summary      Get Patient
// @Description  Get patient by NIK. If not found locally, searches SatuSehat (name required for SatuSehat search).
// @Tags         Patients
// @Produce      json
// @Param        nik   query     string  true   "NIK of the patient"
// @Param        name  query     string  false  "Name of the patient (required if searching SatuSehat)"
// @Success      200   {object}  PersonResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      404   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/patients [get]
func (h *Handlers) GetPatient(w http.ResponseWriter, r *http.Request) {
	nik := r.URL.Query().Get("nik")
	if nik == "" {
		respondError(w, http.StatusBadRequest, "NIK query parameter is required")
		return
	}

	localPatient, err := getPatientLocal(h.db, nik)
	if err == nil && localPatient != nil {
		respondJSON(w, http.StatusOK, localPatient)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		respondError(w, http.StatusBadRequest, "Local patient is not found. Name query parameter is required to search SatuSehat")
		return
	}

	bundle, err := h.ssClient.GetPatient(nik, name)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get patient from SatuSehat: "+err.Error())
		return
	}

	if bundle.Total == 0 || len(bundle.Entry) == 0 {
		respondError(w, http.StatusNotFound, "patient not found")
		return
	}

	person := mapResourceToPerson(&bundle.Entry[0].Resource)

	// if err := savePatientLocal(h.db, person); err != nil {
	// 	fmt.Printf("Warning: failed to save patient locally: %v\n", err)
	// }

	respondJSON(w, http.StatusOK, person)
}

// GetPractitioners godoc
// @Summary      Get Practitioners
// @Description  Get practitioner details from local DB or search SatuSehat by NIK and/or ID
// @Tags         Practitioners
// @Produce      json
// @Param        nik   query     string  false  "NIK of the practitioner"
// @Param        id    query     string  false  "ID/IHS number of the practitioner"
// @Param        page  query     int     false  "Page number"
// @Param        limit query     int     false  "Number of results per page"
// @Success      200   {array}   PersonResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      404   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/practitioners [get]
func (h *Handlers) GetPractitioners(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	nik := r.URL.Query().Get("nik")
	name := r.URL.Query().Get("name")
	gender := r.URL.Query().Get("gender")
	birthdate := r.URL.Query().Get("birthdate")
	if birthdate == "" {
		birthdate = r.URL.Query().Get("birth_date")
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			if l < 0 {
				limit = 0
			} else {
				limit = l
			}
		}
	}

	localPracs, err := searchPractitionersLocal(h.db, id, nik, name, gender, birthdate, page, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to query local practitioners: "+err.Error())
		return
	}

	if len(localPracs) == 0 {
		var ssFetched bool
		if id != "" {
			res, err := h.ssClient.GetPractitionerByID(id)
			if err != nil {
				respondError(w, http.StatusInternalServerError, "failed to get practitioner from SatuSehat by ID: "+err.Error())
				return
			}
			if res != nil {
				person := mapResourceToPerson(res)
				savePractitionerLocal(h.db, person)
				ssFetched = true
			}
		} else {
			var bundle *satusehat.Bundle
			var err error
			if nik != "" {
				identifier := "https://fhir.kemkes.go.id/id/nik|" + nik
				bundle, err = h.ssClient.SearchPractitioners(name, gender, birthdate, identifier)
			} else if name != "" || gender != "" || birthdate != "" {
				bundle, err = h.ssClient.SearchPractitioners(name, gender, birthdate, "")
			}

			if err != nil {
				respondError(w, http.StatusInternalServerError, "failed to search practitioners on SatuSehat: "+err.Error())
				return
			}

			if bundle != nil && bundle.Total > 0 && len(bundle.Entry) > 0 {
				for _, entry := range bundle.Entry {
					person := mapResourceToPerson(&entry.Resource)
					
					// Auto unmask if name contains *
					if strings.Contains(person.Name, "*") && person.ID != "" {
						res, err := h.ssClient.GetPractitionerByID(person.ID)
						if err == nil && res != nil {
							unmasked := mapResourceToPerson(res)
							if unmasked.Name != "" && !strings.Contains(unmasked.Name, "*") {
								person.Name = unmasked.Name
							}
						}
					}

					savePractitionerLocal(h.db, person)
				}
				ssFetched = true
			}
		}

		if ssFetched {
			localPracs, err = searchPractitionersLocal(h.db, id, nik, name, gender, birthdate, page, limit)
			if err != nil {
				respondError(w, http.StatusInternalServerError, "failed to query local practitioners after sync: "+err.Error())
				return
			}
		}
	} else {
		// Auto unmask any existing cached records on-the-fly if they contain "*"
		for i := range localPracs {
			if strings.Contains(localPracs[i].Name, "*") && localPracs[i].ID != "" {
				res, err := h.ssClient.GetPractitionerByID(localPracs[i].ID)
				if err == nil && res != nil {
					unmasked := mapResourceToPerson(res)
					if unmasked.Name != "" && !strings.Contains(unmasked.Name, "*") {
						localPracs[i].Name = unmasked.Name
						savePractitionerLocal(h.db, localPracs[i])
					}
				}
			}
		}
	}

	if localPracs == nil {
		localPracs = []PersonResponse{}
	}

	respondJSON(w, http.StatusOK, localPracs)
}

// GetAllLocalPractitioners godoc
// @Summary      Get All Local Practitioners
// @Description  Get all practitioners stored in the local SQLite database
// @Tags         Practitioners
// @Produce      json
// @Success      200  {array}   PersonResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/local/practitioners [get]
func (h *Handlers) GetAllLocalPractitioners(w http.ResponseWriter, r *http.Request) {
	pracs, err := getAllPractitionersLocal(h.db)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch local practitioners")
		return
	}
	respondJSON(w, http.StatusOK, pracs)
}

// CreatePatient godoc
// @Summary      Create Patient
// @Description  Create a new patient in SatuSehat and save locally, or retrieve existing by NIK
// @Tags         Patients
// @Accept       json
// @Produce      json
// @Param        body  body      CreatePatientRequest  true  "Create Patient Request Body"
// @Success      201   {object}  PersonResponse
// @Success      200   {object}  PersonResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/patients [post]
func (h *Handlers) CreatePatient(w http.ResponseWriter, r *http.Request) {
	var req CreatePatientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid json payload")
		return
	}

	if req.NIK == "" {
		respondError(w, http.StatusBadRequest, "nik is required")
		return
	}

	localPatient, err := getPatientLocal(h.db, req.NIK)
	if err == nil && localPatient != nil {
		respondJSON(w, http.StatusOK, localPatient)
		return
	}

	if req.Name != "" {
		bundle, err := h.ssClient.GetPatient(req.NIK, req.Name)
		if err == nil && bundle != nil && bundle.Total > 0 && len(bundle.Entry) > 0 {
			person := mapResourceToPerson(&bundle.Entry[0].Resource)
			if err := savePatientLocal(h.db, person); err != nil {
				fmt.Printf("Warning: failed to save existing patient locally: %v\n", err)
			}
			respondJSON(w, http.StatusOK, person)
			return
		}
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

	if err := savePatientLocal(h.db, person); err != nil {
		fmt.Printf("Warning: failed to save created patient locally: %v\n", err)
	}

	fmt.Println(person)
	respondJSON(w, http.StatusCreated, person)
}

// GetLocations godoc
// @Summary      Get Locations
// @Description  Search locations locally or sync from SatuSehat by id or identifier
// @Tags         Locations
// @Produce      json
// @Param        id          query     string  false  "Location ID"
// @Param        identifier  query     string  false  "Location Identifier"
// @Param        page        query     int     false  "Page number (default: 1)"
// @Param        limit       query     int     false  "Limit (default: 10)"
// @Success      200         {array}   LocationResponse
// @Failure      500         {object}  ErrorResponse
// @Router       /api/locations [get]
func (h *Handlers) GetLocations(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	identifier := r.URL.Query().Get("identifier")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, limit := 1, 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			if l < 0 {
				limit = 0
			} else {
				limit = l
			}
		}
	}

	localLocs, err := searchLocationsLocal(h.db, id, identifier, page, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to query local locations: "+err.Error())
		return
	}

	if len(localLocs) == 0 {
		var ssFetched bool
		if id != "" {
			res, err := h.ssClient.GetLocationById(id)
			if err != nil {
				respondError(w, http.StatusInternalServerError, "failed to get location from SatuSehat by ID: "+err.Error())
				return
			}
			if res != nil {
				loc := LocationResponse{
					ID:          res.ID,
					Name:        res.Name,
					Description: res.Description,
				}
				if len(res.Identifier) > 0 {
					loc.IdentifierValue = res.Identifier[0].Value
				}
				if len(res.Telecom) > 0 {
					loc.Phone = res.Telecom[0].Value
				}
				saveLocationLocal(h.db, loc)
				ssFetched = true
			}
		} else if identifier != "" {
			bundle, err := h.ssClient.GetLocationByIdentifier(h.ssClient.OrgID(), identifier)
			if err != nil {
				respondError(w, http.StatusInternalServerError, "failed to get location from SatuSehat by Identifier: "+err.Error())
				return
			}
			if bundle != nil && bundle.Total > 0 && len(bundle.Entry) > 0 {
				for _, entry := range bundle.Entry {
					res := entry.Resource
					loc := LocationResponse{
						ID:          res.ID,
						Name:        res.Name,
						Description: res.Description,
					}
					if len(res.Identifier) > 0 {
						loc.IdentifierValue = res.Identifier[0].Value
					}
					if len(res.Telecom) > 0 {
						loc.Phone = res.Telecom[0].Value
					}
					saveLocationLocal(h.db, loc)
				}
				ssFetched = true
			}
		}

		if ssFetched {
			localLocs, err = searchLocationsLocal(h.db, id, identifier, page, limit)
			if err != nil {
				respondError(w, http.StatusInternalServerError, "failed to query local locations after sync: "+err.Error())
				return
			}
		}
	}

	if localLocs == nil {
		localLocs = []LocationResponse{}
	}

	respondJSON(w, http.StatusOK, localLocs)
}

// CreateLocation godoc
// @Summary      Create Location
// @Description  Create a new location in SatuSehat and save locally
// @Tags         Locations
// @Accept       json
// @Produce      json
// @Param        body  body      CreateLocationRequest  true  "Create Location Request Body"
// @Success      201   {object}  LocationResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/locations [post]
func (h *Handlers) CreateLocation(w http.ResponseWriter, r *http.Request) {
	var req CreateLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid json payload")
		return
	}

	payload := &satusehat.LocationResource{
		ResourceType: "Location",
		Status:       "active",
		Name:         req.Name,
		Description:  req.Description,
		Mode:         "instance",
		Identifier: []satusehat.Identifier{
			{System: "http://sys-ids.kemkes.go.id/location/" + h.ssClient.OrgID(), Value: req.IdentifierValue},
		},
		ManagingOrganization: &satusehat.Reference{Reference: "Organization/" + h.ssClient.OrgID()},
	}
	if req.Phone != "" {
		payload.Telecom = append(payload.Telecom, satusehat.Telecom{System: "phone", Value: req.Phone, Use: "work"})
	}

	res, err := h.ssClient.CreateLocation(payload)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create location")
		return
	}

	loc := LocationResponse{
		ID:              res.ID,
		IdentifierValue: req.IdentifierValue,
		Name:            res.Name,
		Description:     res.Description,
		Phone:           req.Phone,
	}
	saveLocationLocal(h.db, loc)
	respondJSON(w, http.StatusCreated, loc)
}

// GetEncounterById godoc
// @Summary      Get Encounter by ID
// @Description  Get encounter by ID from local database or SatuSehat
// @Tags         Encounters
// @Produce      json
// @Param        id   path      string  true  "Encounter ID"
// @Success      200  {object}  EncounterResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/encounters/{id} [get]
func (h *Handlers) GetEncounterById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "id path parameter is required")
		return
	}

	localEnc, err := getEncounterLocal(h.db, id)
	if err == nil && localEnc != nil {
		respondJSON(w, http.StatusOK, localEnc)
		return
	}

	res, err := h.ssClient.GetEncounterById(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get encounter from SatuSehat")
		return
	}

	enc := EncounterResponse{
		ID:     res.ID,
		Status: res.Status,
	}
	if len(res.Identifier) > 0 {
		enc.IdentifierValue = res.Identifier[0].Value
	}
	if res.Subject != nil {
		enc.SubjectID = res.Subject.Reference
	}
	if len(res.Location) > 0 {
		enc.LocationID = res.Location[0].Location.Reference
	}
	if res.Period != nil {
		enc.StartTime = res.Period.Start
	}

	saveEncounterLocal(h.db, enc)
	respondJSON(w, http.StatusOK, enc)
}

// CreateEncounter godoc
// @Summary      Create Encounter
// @Description  Create a new encounter in SatuSehat and save locally
// @Tags         Encounters
// @Accept       json
// @Produce      json
// @Param        body  body      CreateEncounterRequest  true  "Create Encounter Request Body"
// @Success      201   {object}  EncounterResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/encounters [post]
func (h *Handlers) CreateEncounter(w http.ResponseWriter, r *http.Request) {
	var req CreateEncounterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid json payload")
		return
	}

	payload := &satusehat.EncounterResource{
		ResourceType: "Encounter",
		Status:       "arrived",
		Class: &satusehat.Coding{
			System:  "http://terminology.hl7.org/CodeSystem/v3-ActCode",
			Code:    "AMB",
			Display: "ambulatory",
		},
		Subject: &satusehat.Reference{Reference: req.SubjectID},
		Participant: []satusehat.Participant{
			{
				Type: []satusehat.CodeableConcept{{
					Coding: []satusehat.Coding{{System: "http://terminology.hl7.org/CodeSystem/v3-ParticipationType", Code: "ATND"}},
				}},
				Individual: satusehat.Reference{Reference: req.PractitionerID},
			},
		},
		Period:          &satusehat.Period{Start: req.StartTime},
		Location:        []satusehat.EncounterLocation{{Location: satusehat.Reference{Reference: req.LocationID}}},
		ServiceProvider: &satusehat.Reference{Reference: "Organization/" + h.ssClient.OrgID()},
		Identifier: []satusehat.Identifier{
			{System: "http://sys-ids.kemkes.go.id/encounter/" + h.ssClient.OrgID(), Value: req.IdentifierValue},
		},
		StatusHistory: []satusehat.StatusHistory{
			{Status: "arrived", Period: satusehat.Period{Start: req.StartTime}},
		},
	}

	res, err := h.ssClient.CreateEncounter(payload)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create encounter")
		return
	}

	enc := EncounterResponse{
		ID:              res.ID,
		IdentifierValue: req.IdentifierValue,
		Status:          res.Status,
		SubjectID:       req.SubjectID,
		LocationID:      req.LocationID,
		StartTime:       req.StartTime,
	}
	saveEncounterLocal(h.db, enc)
	respondJSON(w, http.StatusCreated, enc)
}

// UpdateEncounterStatus godoc
// @Summary      Update Encounter Status
// @Description  Update the status of an encounter in SatuSehat and locally
// @Tags         Encounters
// @Accept       json
// @Produce      json
// @Param        id    path      string                  true  "Encounter ID"
// @Param        body  body      UpdateEncounterRequest  true  "Update Encounter Status Request Body"
// @Success      200   {object}  EncounterResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/encounters/{id} [put]
func (h *Handlers) UpdateEncounterStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "id path parameter is required")
		return
	}

	var req UpdateEncounterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid json payload")
		return
	}

	// We first need the existing encounter to update it
	res, err := h.ssClient.GetEncounterById(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch existing encounter for update")
		return
	}

	res.Status = req.Status
	// Depending on status, you may want to append to StatusHistory
	// res.StatusHistory = append(res.StatusHistory, satusehat.StatusHistory{Status: req.Status})

	updated, err := h.ssClient.UpdateEncounter(id, res)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update encounter in SatuSehat")
		return
	}

	updateEncounterStatusLocal(h.db, id, req.Status)

	enc := EncounterResponse{
		ID:     updated.ID,
		Status: updated.Status,
	}
	if len(updated.Identifier) > 0 {
		enc.IdentifierValue = updated.Identifier[0].Value
	}
	if updated.Subject != nil {
		enc.SubjectID = updated.Subject.Reference
	}
	if len(updated.Location) > 0 {
		enc.LocationID = updated.Location[0].Location.Reference
	}
	if updated.Period != nil {
		enc.StartTime = updated.Period.Start
	}

	respondJSON(w, http.StatusOK, enc)
}
