package satusehat

type AccessTokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        string `json:"expires_in"`
	TokenType        string `json:"token_type"`
	DeveloperEmail   string `json:"developer.email"`
	OrganizationName string `json:"organization_name"`
}

type Bundle struct {
	ResourceType string  `json:"resourceType"`
	Total        int     `json:"total"`
	Entry        []Entry `json:"entry"`
}

type Entry struct {
	Resource Resource `json:"resource"`
}

type Resource struct {
	ResourceType string       `json:"resourceType"`
	ID           string       `json:"id,omitempty"`
	Identifier   []Identifier `json:"identifier,omitempty"`
	Active       bool         `json:"active,omitempty"`
	Name         []Name       `json:"name,omitempty"`
	Gender       string       `json:"gender,omitempty"`
	BirthDate    string       `json:"birthDate,omitempty"`
	Address      []Address    `json:"address,omitempty"`
	Telecom      []Telecom    `json:"telecom,omitempty"`
	Contact      []Contact    `json:"contact,omitempty"`
}

type Identifier struct {
	System string `json:"system"`
	Value  string `json:"value"`
	Use    string `json:"use,omitempty"`
}

type Name struct {
	Use  string `json:"use,omitempty"`
	Text string `json:"text,omitempty"`
}

type Address struct {
	Use        string      `json:"use,omitempty"`
	Line       []string    `json:"line,omitempty"`
	City       string      `json:"city,omitempty"`
	PostalCode string      `json:"postalCode,omitempty"`
	Country    string      `json:"country,omitempty"`
	Extension  []Extension `json:"extension,omitempty"`
}

type Extension struct {
	URL       string      `json:"url,omitempty"`
	Extension []Extension `json:"extension,omitempty"`
	ValueCode string      `json:"valueCode,omitempty"`
}

type Contact struct {
	Name    Name      `json:"name,omitempty"`
	Telecom []Telecom `json:"telecom,omitempty"`
}

type Telecom struct {
	System string `json:"system,omitempty"`
	Value  string `json:"value,omitempty"`
	Use    string `json:"use,omitempty"`
}

type PatientPayload struct {
	ResourceType         string          `json:"resourceType"`
	Meta                 Meta            `json:"meta"`
	Identifier           []Identifier    `json:"identifier"`
	Active               bool            `json:"active"`
	Name                 []Name          `json:"name"`
	Gender               string          `json:"gender"`
	BirthDate            string          `json:"birthDate"`
	DeceasedBoolean      bool            `json:"deceasedBoolean"`
	Address              []Address       `json:"address"`
	MaritalStatus        CodeableConcept `json:"maritalStatus"`
	MultipleBirthInteger int             `json:"multipleBirthInteger"`
	Contact              []Contact       `json:"contact"`
	Telecom              []Telecom       `json:"telecom,omitempty"`
	Communication        []Communication `json:"communication"`
}

type Meta struct {
	Profile []string `json:"profile"`
}

type CodeableConcept struct {
	Coding []Coding `json:"coding"`
	Text   string   `json:"text,omitempty"`
}

type Coding struct {
	System  string `json:"system"`
	Code    string `json:"code"`
	Display string `json:"display,omitempty"`
}

type Communication struct {
	Language  CodeableConcept `json:"language"`
	Preferred bool            `json:"preferred"`
}

// Location Models
type LocationBundle struct {
	Total int `json:"total"`
	Entry []struct {
		Resource LocationResource `json:"resource"`
	} `json:"entry"`
}

type LocationResource struct {
	ResourceType         string           `json:"resourceType"`
	ID                   string           `json:"id,omitempty"`
	Identifier           []Identifier     `json:"identifier,omitempty"`
	Status               string           `json:"status,omitempty"`
	Name                 string           `json:"name,omitempty"`
	Description          string           `json:"description,omitempty"`
	Mode                 string           `json:"mode,omitempty"`
	Telecom              []Telecom        `json:"telecom,omitempty"`
	Address              *Address         `json:"address,omitempty"`
	PhysicalType         *CodeableConcept `json:"physicalType,omitempty"`
	Position             *Position        `json:"position,omitempty"`
	ManagingOrganization *Reference       `json:"managingOrganization,omitempty"`
}

type Position struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Altitude  float64 `json:"altitude"`
}

// Encounter Models
type EncounterBundle struct {
	Total int `json:"total"`
	Entry []struct {
		Resource EncounterResource `json:"resource"`
	} `json:"entry"`
}

type EncounterResource struct {
	ResourceType    string              `json:"resourceType"`
	ID              string              `json:"id,omitempty"`
	Identifier      []Identifier        `json:"identifier,omitempty"`
	Status          string              `json:"status,omitempty"`
	Class           *Coding             `json:"class,omitempty"`
	Subject         *Reference          `json:"subject,omitempty"`
	Participant     []Participant       `json:"participant,omitempty"`
	Period          *Period             `json:"period,omitempty"`
	Location        []EncounterLocation `json:"location,omitempty"`
	StatusHistory   []StatusHistory     `json:"statusHistory,omitempty"`
	ServiceProvider *Reference          `json:"serviceProvider,omitempty"`
}

type Reference struct {
	Reference string `json:"reference,omitempty"`
	Display   string `json:"display,omitempty"`
}

type Participant struct {
	Type       []CodeableConcept `json:"type,omitempty"`
	Individual Reference         `json:"individual,omitempty"`
}

type Period struct {
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
}

type EncounterLocation struct {
	Location Reference `json:"location,omitempty"`
}

type StatusHistory struct {
	Status string `json:"status,omitempty"`
	Period Period `json:"period,omitempty"`
}
