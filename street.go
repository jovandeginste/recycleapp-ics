package main

import (
	"fmt"
	"log"
	"net/url"
)

type StreetResponse struct {
	Items []struct {
		ID string `json:"id"`
	} `json:"items"`
	Total int `json:"total"`
	Pages int `json:"pages"`
	Page  int `json:"page"`
	Size  int `json:"size"`
}

func getStreetID(zipcode string, street string, token string) (string, error) {
	v := url.Values{}

	v.Set("q", street)
	v.Set("zipcodes", zipcode)

	log.Printf("Fetching %#v with: %v", streetsURL, v)

	fullURL := streetsURL + "?" + v.Encode()

	var result StreetResponse

	if err := getJSON(fullURL, token, &result); err != nil {
		return "", err
	}

	if len(result.Items) == 0 {
		return "", fmt.Errorf("%w: zipcode=%#v street=%#v", ErrStreetNoResult, zipcode, street)
	}

	if result.Items[0].ID == "" {
		return "", fmt.Errorf("%w: zipcode=%#v street=%#v", ErrStreetNoResult, zipcode, street)
	}

	log.Printf("Street ID is: %#v", result.Items[0].ID)

	return result.Items[0].ID, nil
}
