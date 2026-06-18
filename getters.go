package main

import (
	"fmt"
	"log/slog"
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

	slog.Info("Fetching URL", "url", fullURL)

	var result Organization

	if err := getJSON(fullURL, &result); err != nil {
		return nil, err
	}

	if result.Name == "" {
		return nil, fmt.Errorf("%w: %s", ErrOrganizationNoResult, zipcode)
	}

	slog.Info("Organization detail", "name", result.Name)

	return &result, nil
}

func getStreetID(zipcode, street string) (string, error) {
	v := url.Values{}

	v.Set("q", street)
	v.Set("zipcodes", zipcode)

	slog.Info("Fetching streets", "url", streetsURL, "values", v)

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

	slog.Info("Street ID details", "id", result.Items[0].ID)

	return result.Items[0].ID, nil
}

func getZipcodeID(zipcode int) (string, error) {
	v := url.Values{}

	v.Set("q", strconv.Itoa(zipcode))

	slog.Info("Fetching zipcodes", "url", zipcodesURL, "values", v)

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

	slog.Info("Zipcode ID details", "id", result.Items[0].ID)

	return result.Items[0].ID, nil
}
