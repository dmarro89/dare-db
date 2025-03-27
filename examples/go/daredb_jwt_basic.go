package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	BaseURL  string `env:"BASE_URL" env-default:"https://127.0.0.1:2605"`
	Username string `env:"DB_USERNAME" env-default:"admin"`
	Password string `env:"DB_PASSWORD" env-default:"password"`
}

type DareDBPyClientBase struct {
	Username   string
	Password   string
	BaseURL    string
	JWTToken   string
	AuthURL    string
	HTTPClient *http.Client
}

func loadConfig() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}
	return &cfg, nil
}

func NewDareDBPyClientBase(username, password, baseURL string) *DareDBPyClientBase {
	client := &DareDBPyClientBase{
		Username:   username,
		Password:   password,
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Transport: &http.Transport{
			// InsecureSkipVerify: true, // Be cautious using this in production
		}},
	}
	client.AuthURL = client.BaseURL + "/login"
	client.JWTToken = client.getJWTToken()
	log.Printf("JWT token: %s", client.JWTToken)

	return client
}

func (c *DareDBPyClientBase) getJWTToken() string {

	if c.JWTToken != "" {
		return c.JWTToken
	}

	log.Printf("URL to get JWT token: %s", c.AuthURL)

	req, err := http.NewRequest("POST", c.AuthURL, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
		return ""
	}

	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error getting JWT token: status code %d", resp.StatusCode)
		return ""

	}

	var tokenData map[string]string
	err = json.NewDecoder(resp.Body).Decode(&tokenData)
	if err != nil {
		log.Fatalf("Error decoding response: %v", err)
		return ""
	}

	c.JWTToken = tokenData["token"]

	return c.JWTToken
}

func (c *DareDBPyClientBase) buildHeadersWithJWT(jwtToken string) map[string]string {
	token := jwtToken
	if token == "" {
		token = c.JWTToken
	}
	return map[string]string{
		"Authorization": token,
		"Content-Type":  "application/json",
	}
}

func (c *DareDBPyClientBase) sendPostWithJWT(urlStr string, data map[string]interface{}) (*http.Response, error) {
	headers := c.buildHeadersWithJWT("")
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	return resp, nil
}

func (c *DareDBPyClientBase) sendGetWithJWT(urlStr string) (*http.Response, error) {
	headers := c.buildHeadersWithJWT("")

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	return resp, nil
}

func (c *DareDBPyClientBase) sendDeleteWithJWT(urlStr string) (*http.Response, error) {
	headers := c.buildHeadersWithJWT("")

	req, err := http.NewRequest("DELETE", urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	return resp, nil
}

func (c *DareDBPyClientBase) logResponse(resp *http.Response) {
	body := new(bytes.Buffer)
	body.ReadFrom(resp.Body)
	log.Printf("HTTP Code: %d; content: %s", resp.StatusCode, body.String())
}

func main() {

	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	dareDBClient := NewDareDBPyClientBase(cfg.Username, cfg.Password, cfg.BaseURL)

	exampleURL, _ := url.Parse(dareDBClient.BaseURL + "/set")
	exampleData := map[string]interface{}{"keyToSave": "valueToSave"}
	log.Printf("URL to set value: %s\n", exampleURL)

	postResp, postErr := dareDBClient.sendPostWithJWT(exampleURL.String(), exampleData)
	defer postResp.Body.Close()

	if postErr != nil {
		log.Printf("POST error: %v", postErr)
	} else {
		log.Println("POST request results ->")
		dareDBClient.logResponse(postResp)
	}

	exampleURLGet, _ := url.Parse(dareDBClient.BaseURL + "/get/" + "keyToSave")
	getResp, getErr := dareDBClient.sendGetWithJWT(exampleURLGet.String())
	defer getResp.Body.Close()
	if getErr != nil {
		log.Printf("GET error: %v", getErr)
	} else {
		log.Println("GET request results ->")
		dareDBClient.logResponse(getResp)

	}
	exampleURLDelete, _ := url.Parse(dareDBClient.BaseURL + "/delete/" + "keyToSave")
	deleteResp, deleteErr := dareDBClient.sendDeleteWithJWT(exampleURLDelete.String())
	defer deleteResp.Body.Close()

	if deleteErr != nil {
		log.Printf("DELETE error: %v", deleteErr)
	} else {
		log.Println("DELETE request results ->")
		dareDBClient.logResponse(deleteResp)
	}

}
