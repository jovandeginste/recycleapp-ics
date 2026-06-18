package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
)

func (r *Organization) URLForLanguage(lang string) string {
	if u, ok := r.URL[lang]; ok {
		return u
	}

	return "???"
}

func getOrganization(zipcode string) (*Organization, error) {
	fullURL := organizationsURL + zipcode

	log.Printf("Fetching %#v", fullURL)

	var result Organization

	if err := getJSON(fullURL, &result); err != nil {
		return nil, err
	}

	if result.Name == "" {
		return nil, fmt.Errorf("%w: %s", ErrOrganizationNoResult, zipcode)
	}

	log.Printf("Organization is: %#v", result.Name)

	return &result, nil
}

func getStreetID(zipcode, street string) (string, error) {
	v := url.Values{}

	v.Set("q", street)
	v.Set("zipcodes", zipcode)

	log.Printf("Fetching %#v with: %v", streetsURL, v)

	fullURL := streetsURL + "?" + v.Encode()

	var result StreetResponse

	if err := getJSON(fullURL, &result); err != nil {
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

func getZipcodeID(zipcode int) (string, error) {
	v := url.Values{}

	v.Set("q", strconv.Itoa(zipcode))

	log.Printf("Fetching %#v with: %v", zipcodesURL, v)

	fullURL := zipcodesURL + "?" + v.Encode()

	var result ZipResponse

	if err := getJSON(fullURL, &result); err != nil {
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
