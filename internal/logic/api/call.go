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

func unmarshalToMap(bodyBytes []byte, mapInst map[string]interface{}) error {
	return json.Unmarshal(bodyBytes, &mapInst)
}

func unmarshalToPlayer(bodyBytes []byte, playerSI *PlayerSI) error {
	return json.Unmarshal(bodyBytes, playerSI)
}
func GetMe(userID int64) (map[string]interface{}, error) {
	var response *http.Response
	var err error
	var data map[string]interface{}
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
		_ = b.Close()
	}(response.Body)

	if bodyBytes, err = io.ReadAll(response.Body); err != nil {
		return nil, err
	}

	if err = unmarshalToMap(bodyBytes, data); err != nil {
		return nil, err
	}

	return data, nil
}

func GetPlayer(userID int64) (*PlayerSI, error) {
	var response *http.Response
	var err error
	var playerSI PlayerSI
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

	if response.StatusCode == 204 {
		return nil, nil // No Content from Spotify API, player is inactive
	}

	defer func(b io.ReadCloser) {
		b.Close()
	}(response.Body)

	if bodyBytes, err = io.ReadAll(response.Body); err != nil {
		return nil, err
	}

	if err = unmarshalToPlayer(bodyBytes, &playerSI); err != nil {
		return nil, err
	}

	return &playerSI, nil
}
