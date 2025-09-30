package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var Domain string

func makeRequest(method string, path string) (response *http.Response, err error) {
	var uri *url.URL
	if uri, err = url.Parse(Domain + path); err != nil {
		return nil, err
	}

	if response, err = (&http.Client{}).
		Do(
			&http.Request{
				Method: method,
				URL:    uri,
				Header: map[string][]string{
					"Accept": {"application/json"},
				},
			},
		); err != nil {
		return nil, err
	}

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("status code %s", response.Status)
	}

	return response, nil
}

func unmarshalToMap(bodyBytes []byte) (mapInst map[string]interface{}, err error) {
	err = json.Unmarshal(bodyBytes, &mapInst)
	return
}

func unmarshalToPlayer(bodyBytes []byte) (*PlayerSI, error) {
	var player PlayerSI
	err := json.Unmarshal(bodyBytes, &player)
	return &player, err
}
func GetMe(userID int64) (map[string]interface{}, error) {
	var response *http.Response
	var err error
	var bodyBytes []byte

	if response, err = makeRequest(
		"GET",
		fmt.Sprintf(
			"/public/spotify/api/%v/me",
			userID,
		),
	); err != nil {
		return nil, err
	}

	defer func(b io.ReadCloser) {
		b.Close()
	}(response.Body)

	if bodyBytes, err = io.ReadAll(response.Body); err != nil {
		return nil, err
	}

	return unmarshalToMap(bodyBytes)
}

func GetPlayer(userID int64) (*PlayerSI, error) {
	var response *http.Response
	var err error
	var bodyBytes []byte

	if response, err = makeRequest(
		"GET",
		fmt.Sprintf(
			"/public/spotify/api/%v/me/player",
			userID,
		),
	); err != nil || response.StatusCode == 204 {
		return nil, err
	}

	defer func(b io.ReadCloser) {
		b.Close()
	}(response.Body)

	if bodyBytes, err = io.ReadAll(response.Body); err != nil {
		return nil, err
	}

	return unmarshalToPlayer(bodyBytes)
}
