package appendedGo

import (
	"encoding/json"
	"errors"
	"fmt"
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

// Create a new Appended client
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

// CreateFolio creates a new folio
func (c *Client) CreateFolio(folioName string) error {
	_, err := c.makeRequest("POST", "/folios", map[string]string{"name": folioName})
	if err != nil {
		return err
	}

	return nil
}

// GetFolios will return the names of all folios
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

// GetNotes will return a slice of all notes' text for a given folio
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

// AddNote will add a note to a folio
func (c *Client) AddNote(folioName string, note string) error {
	_, err := c.makeRequest("POST", "/folios/"+folioName, map[string]string{"note": note})
	if err != nil {
		return err
	}

	return nil
}

// EditNote will overwrite the text of a note
func (c *Client) EditNote(folioName string, index int, note string) error {
	path := fmt.Sprintf("/folios/%v/%v", folioName, index)
	_, err := c.makeRequest("PUT", path, map[string]string{"note": note})
	if err != nil {
		return err
	}

	return nil
}

// ToggleDone will toggle the done property on a note
func (c *Client) ToggleDone(folioName string, index int) error {
	path := fmt.Sprintf("/folios/%v/%v/done", folioName, index)
	_, err := c.makeRequest("GET", path, nil)
	if err != nil {
		return err
	}

	return nil
}

// DeleteFolio will delete a folio permanently
func (c *Client) DeleteFolio(folioName string) error {
	_, err := c.makeRequest("DELETE", "/folios/"+folioName, nil)
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
