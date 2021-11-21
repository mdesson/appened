package appendedGo

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	token  string
	client *http.Client
	url    string
}

func New(token string, baseURL string) *Client {

	client := Client{}
	client.token = token
	client.client = &http.Client{}

	if baseURL[len(baseURL)-1:] == "/" {
		client.url = baseURL[:len(baseURL)-1]
	} else {
		client.url = baseURL
	}

	return &client
}

func (c *Client) CreateFolio(folioName string) error {
	_, err := c.makeRequest("POST", "/folios", map[string]string{"name": folioName})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetFolios() ([]string, error) {
	body, err := c.makeRequest("GET", "/folios", nil)
	if err != nil {
		return nil, err
	}

	folioNames := []string{}

	err = json.Unmarshal(body, &folioNames)
	if err != nil {
		return nil, err
	}

	return folioNames, nil
}

func (c *Client) GetNotes(folioName string) ([]string, error) {
	body, err := c.makeRequest("GET", "/folios/"+folioName, nil)
	if err != nil {
		return nil, err
	}

	notes := []string{}

	err = json.Unmarshal(body, &notes)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (c *Client) AddNote(folioName string, note string) error {
	_, err := c.makeRequest("POST", "/folios/"+folioName, map[string]string{"note": note})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) makeRequest(method string, route string, data map[string]string) ([]byte, error) {
	postData := url.Values{}
	for key, val := range data {
		postData.Set(key, val)
	}

	req, err := http.NewRequest(method, c.url+route, strings.NewReader(postData.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.token)
	if data != nil {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 201 {
		msg := http.StatusText(resp.StatusCode)
		return nil, errors.New(msg)

	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}
