package satusehat

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	TokenType   string `json:"token_type"`
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
