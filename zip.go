package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
)

type ZipResponse struct {
	Items []struct {
		ID string `json:"id"`
	} `json:"items"`
	Total int `json:"total"`
	Pages int `json:"pages"`
	Page  int `json:"page"`
	Size  int `json:"size"`
}

func getZipcodeID(zipcode int, token string) (string, error) {
	v := url.Values{}

	v.Set("q", strconv.Itoa(zipcode))

	log.Printf("Fetching %#v with: %v", zipcodesURL, v)

	fullURL := zipcodesURL + "?" + v.Encode()

	var result ZipResponse

	if err := getJSON(fullURL, token, &result); err != nil {
		return "", err
	}

	if len(result.Items) == 0 {
		return "", fmt.Errorf("%w: %d", ErrZipcodeNoResult, zipcode)
	}

	if result.Items[0].ID == "" {
		return "", fmt.Errorf("%w: %d", ErrZipcodeNoResult, zipcode)
	}

	log.Printf("Zipcode ID is: %#v", result.Items[0].ID)

	return result.Items[0].ID, nil
}
