package satusehat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bukunya/intero-go/internal/config"
)

type Client struct {
	config     *config.Config
	httpClient *http.Client
	token      string
	tokenExp   time.Time
	mu         sync.Mutex
}

func (c *Client) OrgID() string {
	return c.config.OrgID
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetToken() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.token != "" && time.Now().Before(c.tokenExp.Add(-1*time.Minute)) {
		return c.token, nil
	}

	data := url.Values{}
	data.Set("client_id", c.config.ClientID)
	data.Set("client_secret", c.config.ClientSecret)

	reqURL := fmt.Sprintf("%s/accesstoken?grant_type=client_credentials", c.config.AuthURL)
	req, err := http.NewRequest("POST", reqURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("failed to get token: %s", string(body))
	}

	var tokenResp AccessTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	expSeconds, _ := strconv.Atoi(tokenResp.ExpiresIn)
	if expSeconds == 0 {
		expSeconds = 3600
	}

	c.token = tokenResp.AccessToken
	c.tokenExp = time.Now().Add(time.Duration(expSeconds) * time.Second)

	return c.token, nil
}

func (c *Client) Request(method, path string, body interface{}) ([]byte, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, fmt.Errorf("auth error: %w", err)
	}

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(b)
	}

	reqURL := fmt.Sprintf("%s%s", c.config.BaseURL, path)
	req, err := http.NewRequest(method, reqURL, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (c *Client) GetPatient(nik, name string) (*Bundle, error) {
	path := fmt.Sprintf("/Patient?identifier=https://fhir.kemkes.go.id/id/nik|%s&name=%s", nik, url.QueryEscape(name))
	resp, err := c.Request("GET", path, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(resp))

	var bundle Bundle
	if err := json.Unmarshal(resp, &bundle); err != nil {
		return nil, err
	}

	return &bundle, nil
}

func (c *Client) GetPractitionerByNIK(nik string) (*Bundle, error) {
	path := fmt.Sprintf("/Practitioner?identifier=https://fhir.kemkes.go.id/id/nik|%s", nik)
	resp, err := c.Request("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var bundle Bundle
	if err := json.Unmarshal(resp, &bundle); err != nil {
		return nil, err
	}

	return &bundle, nil
}

func (c *Client) SearchPractitioners(name, gender, birthdate, identifier string) (*Bundle, error) {
	params := url.Values{}
	if name != "" {
		params.Set("name", name)
	}
	if gender != "" {
		params.Set("gender", gender)
	}
	if birthdate != "" {
		params.Set("birthdate", birthdate)
	}
	if identifier != "" {
		params.Set("identifier", identifier)
	}

	path := "/Practitioner"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	resp, err := c.Request("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var bundle Bundle
	if err := json.Unmarshal(resp, &bundle); err != nil {
		return nil, err
	}

	return &bundle, nil
}

func (c *Client) GetPractitionerByID(id string) (*Resource, error) {
	path := fmt.Sprintf("/Practitioner/%s", url.QueryEscape(id))
	resp, err := c.Request("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var res Resource
	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) CreatePatient(payload *PatientPayload) (*Resource, error) {
	resp, err := c.Request("POST", "/Patient", payload)
	if err != nil {
		return nil, err
	}

	var resource Resource
	if err := json.Unmarshal(resp, &resource); err != nil {
		return nil, err
	}

	return &resource, nil
}

func (c *Client) GetLocationByIdentifier(orgID, identifier string) (*LocationBundle, error) {
	path := fmt.Sprintf("/Location?identifier=http://sys-ids.kemkes.go.id/location/%s|%s", orgID, url.QueryEscape(identifier))
	resp, err := c.Request("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var bundle LocationBundle
	if err := json.Unmarshal(resp, &bundle); err != nil {
		return nil, err
	}
	return &bundle, nil
}

func (c *Client) GetLocationById(id string) (*LocationResource, error) {
	path := fmt.Sprintf("/Location/%s", url.QueryEscape(id))
	resp, err := c.Request("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var res LocationResource
	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) CreateLocation(payload *LocationResource) (*LocationResource, error) {
	resp, err := c.Request("POST", "/Location", payload)
	if err != nil {
		return nil, err
	}

	var res LocationResource
	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetEncounterById(id string) (*EncounterResource, error) {
	path := fmt.Sprintf("/Encounter/%s", url.QueryEscape(id))
	resp, err := c.Request("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var res EncounterResource
	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) CreateEncounter(payload *EncounterResource) (*EncounterResource, error) {
	resp, err := c.Request("POST", "/Encounter", payload)
	if err != nil {
		return nil, err
	}

	var res EncounterResource
	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) UpdateEncounter(id string, payload *EncounterResource) (*EncounterResource, error) {
	path := fmt.Sprintf("/Encounter/%s", url.QueryEscape(id))
	resp, err := c.Request("PUT", path, payload)
	if err != nil {
		return nil, err
	}

	var res EncounterResource
	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
