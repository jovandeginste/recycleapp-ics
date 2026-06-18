package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/jordic/goics"
	"github.com/spf13/cobra"

	"github.com/jovandeginste/recycleapp-ics/pkg/recycleapp"
)

var (
	myClient = func() *http.Client {
		client := retryablehttp.NewClient()
		client.RetryMax = 4
		client.Logger = slogLeveledLogger{}
		return client.StandardClient()
	}()
)

func main() {
	var (
		zipcode, houseNumber int
		street               string
		year                 int
		format               string
		lang                 string
	)

	rootCmd := &cobra.Command{
		Use:   "recycleapp-ics",
		Short: "Generate iCalendar (ICS) files for the recycleapp.be garbage collection schedule",
		RunE: func(cmd *cobra.Command, args []string) error {
			if format != "json" && format != "ics" {
				return fmt.Errorf("invalid format %q, must be either 'json' or 'ics'", format)
			}

			fromDate := fmt.Sprintf("%d-01-01", year)
			untilDate := fmt.Sprintf("%d-12-31", year)

			ctx := cmd.Context()
			client := recycleapp.NewClient(myClient)

			zipcodeID, err := client.GetZipcodeID(ctx, zipcode)
			if err != nil {
				return err
			}

			org, err := client.GetOrganization(ctx, zipcodeID)
			if err != nil {
				return err
			}

			streetID, err := client.GetStreetID(ctx, zipcodeID, street)
			if err != nil {
				return err
			}

			result, err := client.GetCollections(ctx, recycleapp.CollectionsParams{
				ZipcodeID:   zipcodeID,
				StreetID:    streetID,
				HouseNumber: houseNumber,
				FromDate:    fromDate,
				UntilDate:   untilDate,
				Lang:        lang,
			})
			if err != nil {
				return err
			}

			result.Org = org

			if format == "json" {
				jsonData, err := json.MarshalIndent(result.ToJSONEvents(), "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(jsonData))
			} else {
				b := bytes.Buffer{}
				goics.NewICalEncode(&b).Encode(result)
				fmt.Println(b.String())
			}
			return nil
		},
	}

	rootCmd.Flags().StringVar(&lang, "lang", "nl", "your language (nl, fr, en, de)")
	rootCmd.Flags().IntVar(&zipcode, "zipcode", 0, "your zip code")
	rootCmd.Flags().StringVar(&street, "street", "", "your street name")
	rootCmd.Flags().IntVar(&houseNumber, "house", 0, "your house number (digits only)")
	rootCmd.Flags().IntVar(&year, "year", time.Now().Year(), "the year")
	rootCmd.Flags().StringVar(&format, "format", "json", "output format (json, ics)")

	// Make mandatory flags required (if appropriate, otherwise keep optional)
	_ = rootCmd.MarkFlagRequired("zipcode")
	_ = rootCmd.MarkFlagRequired("street")
	_ = rootCmd.MarkFlagRequired("house")

	if err := rootCmd.Execute(); err != nil {
		slog.Error("Execution failed", "error", err)
		os.Exit(1)
	}
}

type slogLeveledLogger struct{}

func (s slogLeveledLogger) Error(msg string, keysAndValues ...any) {
	slog.Error(msg, keysAndValues...)
}

func (s slogLeveledLogger) Info(msg string, keysAndValues ...any) {
	slog.Info(msg, keysAndValues...)
}

func (s slogLeveledLogger) Debug(msg string, keysAndValues ...any) {
	slog.Debug(msg, keysAndValues...)
}

func (s slogLeveledLogger) Warn(msg string, keysAndValues ...any) {
	slog.Warn(msg, keysAndValues...)
}
