package goSDK

import (
	"fmt"
	"net/http"
)

type Client struct {
	token  string
	client *http.Client
}

func New(token string) *Client {

	client := Client{}
	client.token = token
	client.client = &http.Client{}

	return &client
}

func (c *Client) GetNotes(folioName string) ([]string, error) {
	req, err := http.NewRequest("GET", "localhost:8081", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.token)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	fmt.Println("hello!")
	fmt.Println(resp)

	return nil, nil
}
