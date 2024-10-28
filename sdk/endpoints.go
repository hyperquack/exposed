package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

type Object struct {
	Namespace string    `json:"namespace,omitempty"`
	Host      string    `json:"host,omitempty"`
	Value     string    `json:"value,omitempty"`
	Updated   time.Time `json:"updated,omitempty"`
}

type Response struct {
	Hits []Object `json:"hits,omitempty"`
}

func (c *CognitoClient) AddTarget(host string) (Response, error) {
	params := map[string]string{"host": host}
	return c.call("POST", "target", params, nil)
}

func (c *CognitoClient) GetTargets() (Response, error) {
	return c.call("GET", "target", nil, nil)
}

func (c *CognitoClient) RemoveTarget(host string) (Response, error) {
	params := map[string]string{"host": host}
	return c.call("DELETE", "target", params, nil)
}

func (c *CognitoClient) SetNotification(address string) (Response, error) {
	params := map[string]string{"url": address}
	return c.call("POST", "config", params, nil)
}

func (c *CognitoClient) Push(namespace, host, value string) (Response, error) {
	uri := fmt.Sprintf("feed/%s", namespace)
	body := Object{Host: host, Value: value}
	return c.call("POST", uri, nil, body)
}

func (c *CognitoClient) Pull(namespace, host string) (Response, error) {
	uri := fmt.Sprintf("feed/%s", namespace)
	params := map[string]string{"host": host}
	return c.call("GET", uri, params, nil)
}

func (c *CognitoClient) call(method, endpoint string, params map[string]string, body interface{}) (Response, error) {
	var result Response

	retry := func() (*http.Response, error) {
		baseURL, err := url.Parse(fmt.Sprintf("%s/%s", c.API, endpoint))
		if err != nil {
			return nil, err
		}

		if params != nil {
			p := url.Values{}
			for key, value := range params {
				p.Add(key, value)
			}
			baseURL.RawQuery = p.Encode()
		}

		var reqBody io.Reader
		if body != nil {
			jsonData, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			reqBody = bytes.NewBuffer(jsonData)
		}

		req, err := http.NewRequest(method, baseURL.String(), reqBody)
		if err != nil {
			return nil, err
		}

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
		if body != nil {
			req.Header.Add("Content-Type", "application/json")
		}

		client := &http.Client{}
		return client.Do(req)
	}

	resp, err := retry()
	if err != nil {
		slog.Error("exposed", "endpoint", endpoint, "error", err)
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		if err := c.Refresh(); err != nil {
			slog.Error("exposed", "endpoint", endpoint, "error", err)
			return result, err
		}

		resp, err = retry()
		if err != nil {
			slog.Error("exposed", "endpoint", endpoint, "error", err)
			return result, err
		}
		defer resp.Body.Close()
	}

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		slog.Error("exposed", "endpoint", endpoint, "error", string(respBody), "code", resp.StatusCode)
		return result, fmt.Errorf("request failed %d", resp.StatusCode)
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return result, err
	}

	return result, nil
}
