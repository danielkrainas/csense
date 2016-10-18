package client

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/danielkrainas/csense/api/v1"
)

type HooksAPI interface {
	CreateHook(r *v1.NewHookRequest) error
}

type hooksAPI struct {
	*Client
}

func (c *Client) Hooks() HooksAPI {
	return &hooksAPI{c}
}

func (c *hooksAPI) CreateHook(hr *v1.NewHookRequest) error {
	body, err := json.Marshal(hr)
	if err != nil {
		return err
	}

	urlStr, err := c.urls().BuildHooks()
	if err != nil {
		return err
	}

	r, err := http.NewRequest(http.MethodPut, urlStr, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	resp, err := c.do(r)
	if err != nil {
		return err
	}

	resp.Body.Close()
	if err != nil {
		return err
	}

	return err
}
