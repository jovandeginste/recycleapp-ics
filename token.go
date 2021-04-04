package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type RecycleToken struct {
	ExpiresAt   time.Time `json:"expiresAt"`
	AccessToken string    `json:"accessToken"`
}

func getToken() (*RecycleToken, error) {
	secret, err := getSecret()
	if err != nil {
		return nil, err
	}

	url := tokenURL

	log.Printf("Fetching access token from %#v", url)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("x-consumer", consumer)
	req.Header.Set("x-secret", secret)

	r, err := myClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var token RecycleToken

	if err = json.NewDecoder(r.Body).Decode(&token); err != nil {
		return nil, err
	}

	log.Printf("Received token: %#v", token.AccessToken)
	log.Printf("Valid until: %s", token.ExpiresAt.In(localLocation))

	return &token, nil
}

func getSecret() (string, error) {
	body, err := getMainPage()
	if err != nil {
		return "", err
	}

	res := jsRegexp.FindStringSubmatch(body)
	if len(res) < 2 {
		return "", ErrNoJSMatch
	}

	jsPath := res[1]

	body, err = getJSPage(jsPath)
	if err != nil {
		return "", err
	}

	res = secretRegexp.FindStringSubmatch(body)
	if len(res) < 2 {
		return "", ErrNoJSMatch
	}

	token := res[1]

	log.Printf("We got a secret: %#v", token)

	return token, nil
}

func getJSPage(path string) (string, error) {
	url := baseURL + path

	log.Printf("Fetching JS page: %#v", url)

	return getPage(url)
}

func getMainPage() (string, error) {
	url := calendarURL

	log.Printf("Fetching main page: %#v", url)

	return getPage(url)
}

func getPage(url string) (string, error) {
	req, _ := http.NewRequest("GET", url, nil)

	r, err := myClient.Do(req)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
