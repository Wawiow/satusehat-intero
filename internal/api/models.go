package api

type PersonResponse struct {
	ID        string `json:"id"`
	NIK       string `json:"nik"`
	IHSNumber string `json:"ihs_number,omitempty"`
	Name      string `json:"name"`
	Gender    string `json:"gender,omitempty"`
	BirthDate string `json:"birth_date,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Address   string `json:"address,omitempty"`
}

type CreatePatientRequest struct {
	NIK          string `json:"nik"`
	Name         string `json:"name"`
	Gender       string `json:"gender"`
	BirthDate    string `json:"birth_date"`
	Phone        string `json:"phone,omitempty"`
	Address      string `json:"address,omitempty"`
	City         string `json:"city,omitempty"`
	ProvinceCode string `json:"province_code,omitempty"`
	CityCode     string `json:"city_code,omitempty"`
	DistrictCode string `json:"district_code,omitempty"`
	VillageCode  string `json:"village_code,omitempty"`
	RT           string `json:"rt,omitempty"`
	RW           string `json:"rw,omitempty"`
	PostalCode   string `json:"postal_code,omitempty"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type LocationResponse struct {
	ID              string `json:"id"`
	IdentifierValue string `json:"identifier_value"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Phone           string `json:"phone,omitempty"`
}

type CreateLocationRequest struct {
	IdentifierValue string `json:"identifier_value"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Phone           string `json:"phone,omitempty"`
}

type EncounterResponse struct {
	ID              string `json:"id"`
	IdentifierValue string `json:"identifier_value"`
	Status          string `json:"status"`
	SubjectID       string `json:"subject_id"`
	LocationID      string `json:"location_id"`
	StartTime       string `json:"start_time"`
}

type CreateEncounterRequest struct {
	IdentifierValue string `json:"identifier_value"`
	SubjectID       string `json:"subject_id"`
	LocationID      string `json:"location_id"`
	PractitionerID  string `json:"practitioner_id"`
	StartTime       string `json:"start_time"`
}

type UpdateEncounterRequest struct {
	Status string `json:"status"`
}
