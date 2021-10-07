package main

import (
	"fmt"
	"log"
)

type Organization struct {
	Name        string            `json:"name"`
	URL         map[string]string `json:"url"`
	Description map[string]string `json:"description"`
}

func (r *Organization) URLForLanguage(lang string) string {
	if u, ok := r.URL[lang]; ok {
		return u
	}

	return "???"
}

func getOrganization(zipcode string, token string) (*Organization, error) {
	fullURL := organizationsURL + zipcode

	log.Printf("Fetching %#v", fullURL)

	var result Organization

	if err := getJSON(fullURL, token, &result); err != nil {
		return nil, err
	}

	if result.Name == "" {
		return nil, fmt.Errorf("%w: %s", ErrOrganizationNoResult, zipcode)
	}

	log.Printf("Organization is: %#v", result.Name)

	return &result, nil
}
