package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"plassstic.tech/gopkg/golang-manager/internal/depend/logger"
)

func makeRequest(method string, path string) (*http.Response, error) {
	uri, err := url.Parse("https://api.plassstic.tech" + path)
	if err != nil {
		return nil, err
	}

	response, err := (&http.Client{}).Do(&http.Request{Method: method, URL: uri, Header: map[string][]string{
		"Accept": {"application/json"},
	}})

	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("status code %s", response.Status)
	}

	return response, nil
}

func GetMe(userID int64) (map[string]interface{}, error) {
	resp, err := makeRequest("GET", fmt.Sprintf("/public/spotify/api/%v/me", userID))
	logger.GetLogger("123").Infof("%#v", resp)
	if err != nil {
		return nil, err
	}

	var r map[string]interface{}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}

	return r, nil
}
