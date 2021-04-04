package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

const (
	consumer         = "recycleapp.be"
	baseURL          = "https://" + consumer + "/"
	calendarURL      = baseURL + "calendar"
	APIURL           = baseURL + "api/app/v1/"
	tokenURL         = APIURL + "access-token"
	collectionsURL   = APIURL + "collections"
	streetsURL       = APIURL + "streets"
	zipcodesURL      = APIURL + "zipcodes"
	organisationsURL = APIURL + "organisations/"
)

var (
	myClient = &http.Client{Timeout: 10 * time.Second}

	lang          string
	localLocation *time.Location
	jsRegexp      = regexp.MustCompile(`src="(/static/js/main.[[:alnum:]]*\.chunk\.js)"`)
	secretRegexp  = regexp.MustCompile(`n="([^"]+)",c="/api/v1/assets/"`)

	ErrNoJSMatch            = errors.New("main page did not contain the expected main js url")
	ErrZipcodeNoResult      = errors.New("zipcode query returned nothing")
	ErrStreetNoResult       = errors.New("street query returned nothing")
	ErrOrganizationNoResult = errors.New("organization query returned nothing")
)

func main() {
	var (
		err                  error
		zipcode, houseNumber int
		street               string
	)

	year := time.Now().Year()

	fromDate := fmt.Sprintf("%d-01-01", year)
	untilDate := fmt.Sprintf("%d-12-31", year)
	size := "200"

	flag.StringVar(&lang, "lang", "nl", "your language (nl, fr, en, de)")
	flag.IntVar(&zipcode, "zipcode", 0, "your zip code")
	flag.StringVar(&street, "street", "", "your street name")
	flag.IntVar(&houseNumber, "house", 0, "your house number (digits only)")
	flag.Parse()

	localLocation, err = time.LoadLocation("Local")
	if err != nil {
		log.Fatal(`Failed to load location "Local"`)
	}

	authorizationToken, err := getToken()
	if err != nil {
		log.Fatal(err)
	}

	token := authorizationToken.AccessToken

	zipcodeID, err := getZipcodeID(zipcode, token)
	if err != nil {
		log.Fatal(err)
	}

	org, err := getOrganization(zipcodeID, token)
	if err != nil {
		log.Fatal(err)
	}

	streetID, err := getStreetID(zipcodeID, street, token)
	if err != nil {
		log.Fatal(err)
	}

	v := url.Values{}

	v.Set("zipcodeId", zipcodeID)
	v.Set("streetId", streetID)
	v.Set("houseNumber", fmt.Sprintf("%d", houseNumber))
	v.Set("fromDate", fromDate)
	v.Set("untilDate", untilDate)
	v.Set("size", size)

	log.Printf("Fetching %#v with:\n%v", collectionsURL, v)

	fullURL := collectionsURL + "?" + v.Encode()

	var result RecycleInfo

	if err := getJSON(fullURL, token, &result); err != nil {
		log.Fatal(err)
	}

	cal := result.ToCalendar(org)

	fmt.Println(cal.Serialize())
}

func getJSON(url string, token string, target interface{}) error {
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("x-consumer", consumer)
	req.Header.Set("Authorization", token)

	r, err := myClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}