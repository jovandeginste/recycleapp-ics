package main

import (
	"fmt"
	"log"
)

type OrganizationsResponse struct {
	Name        string `json:"name"`
	Description struct {
		Nl string `json:"nl"`
		Fr string `json:"fr"`
		En string `json:"en"`
		De string `json:"de"`
	} `json:"description"`
	URL struct {
		Nl string `json:"nl"`
		Fr string `json:"fr"`
		En string `json:"en"`
		De string `json:"de"`
	} `json:"url"`
}

// https://recycleapp.be/api/app/v1/organisations/3110-24094

func getOrganization(zipcode string, token string) (string, error) {
	fullURL := organisationsURL + zipcode

	log.Printf("Fetching %#v", fullURL)

	var result OrganizationsResponse

	if err := getJSON(fullURL, token, &result); err != nil {
		return "", err
	}

	if result.Name == "" {
		return "", fmt.Errorf("%w: %s", ErrOrganizationNoResult, zipcode)
	}

	log.Printf("Organization is: %#v", result.Name)

	return result.Name, nil
}
