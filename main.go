package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/jordic/goics"
	"github.com/spf13/cobra"
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
		zipcode, houseNumber int
		street               string
		year                 int
	)

	rootCmd := &cobra.Command{
		Use:   "recycleapp-ics",
		Short: "Generate iCalendar (ICS) files for the recycleapp.be garbage collection schedule",
		RunE: func(cmd *cobra.Command, args []string) error {
			fromDate := fmt.Sprintf("%d-01-01", year)
			untilDate := fmt.Sprintf("%d-12-31", year)
			size := "200"

			zipcodeID, err := getZipcodeID(zipcode)
			if err != nil {
				return err
			}

			org, err := getOrganization(zipcodeID)
			if err != nil {
				return err
			}

			streetID, err := getStreetID(zipcodeID, street)
			if err != nil {
				return err
			}

			v := url.Values{}
			v.Set("zipcodeId", zipcodeID)
			v.Set("streetId", streetID)
			v.Set("houseNumber", strconv.Itoa(houseNumber))
			v.Set("fromDate", fromDate)
			v.Set("untilDate", untilDate)
			v.Set("size", size)

			slog.Info("Fetching collections", "url", collectionsURL, "values", v)

			fullURL := collectionsURL + "?" + v.Encode()

			var result RecycleInfo
			if err := getJSON(fullURL, &result); err != nil {
				return err
			}

			result.Org = org

			b := bytes.Buffer{}
			goics.NewICalEncode(&b).Encode(result)

			fmt.Println(b.String())
			return nil
		},
	}

	rootCmd.Flags().StringVar(&lang, "lang", "nl", "your language (nl, fr, en, de)")
	rootCmd.Flags().IntVar(&zipcode, "zipcode", 0, "your zip code")
	rootCmd.Flags().StringVar(&street, "street", "", "your street name")
	rootCmd.Flags().IntVar(&houseNumber, "house", 0, "your house number (digits only)")
	rootCmd.Flags().IntVar(&year, "year", time.Now().Year(), "the year")

	// Make mandatory flags required (if appropriate, otherwise keep optional)
	_ = rootCmd.MarkFlagRequired("zipcode")
	_ = rootCmd.MarkFlagRequired("street")
	_ = rootCmd.MarkFlagRequired("house")

	if err := rootCmd.Execute(); err != nil {
		slog.Error("Execution failed", "error", err)
		os.Exit(1)
	}
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
