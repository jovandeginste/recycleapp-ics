package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/jordic/goics"
)

const (
	consumer = "recycleapp.be"
	baseURL  = "https://api.fostplus.be/recyclecms/public/v1"
)

//nolint:errcheck
var (
	collectionsURL, _   = url.JoinPath(baseURL, "collections")
	streetsURL, _       = url.JoinPath(baseURL, "streets")
	zipcodesURL, _      = url.JoinPath(baseURL, "zipcodes")
	organizationsURL, _ = url.JoinPath(baseURL, "organizations/")
)

var (
	myClient = &http.Client{Timeout: 10 * time.Second}

	lang string

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
		year                 int
	)

	flag.StringVar(&lang, "lang", "nl", "your language (nl, fr, en, de)")
	flag.IntVar(&zipcode, "zipcode", 0, "your zip code")
	flag.StringVar(&street, "street", "", "your street name")
	flag.IntVar(&houseNumber, "house", 0, "your house number (digits only)")
	flag.IntVar(&year, "year", time.Now().Year(), "the year")
	flag.Parse()

	fromDate := fmt.Sprintf("%d-01-01", year)
	untilDate := fmt.Sprintf("%d-12-31", year)
	size := "200"

	zipcodeID, err := getZipcodeID(zipcode)
	if err != nil {
		log.Fatal(err)
	}

	org, err := getOrganization(zipcodeID)
	if err != nil {
		log.Fatal(err)
	}

	streetID, err := getStreetID(zipcodeID, street)
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

	if err := getJSON(fullURL, &result); err != nil {
		log.Fatal(err)
	}

	result.Org = org

	b := bytes.Buffer{}

	goics.NewICalEncode(&b).Encode(result)

	fmt.Println(b.String())
}

func getJSON(fullURL string, target any) error {
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("x-consumer", consumer)

	r, err := myClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
