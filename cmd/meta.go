package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// Response represences response from MetaData service
type Response struct {
	Status string `json:"status"`
	Error  any    `json:"error,omitempty"`
}

// helper function to get metadata
// MetaData represents MetaData object returned from discovery service
type MetaData struct {
	ID          string   `json:"id"`
	Site        string   `json:"site"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Bucket      string   `json:"bucket"`
}

// MetaDataRecord represents MetaData record returned by discovery service
type MetaDataRecord struct {
	Status string     `json:"status"`
	Data   []MetaData `json:"data"`
}

// helper function to fetch sites info from discovery service
func metadata(site string) MetaDataRecord {
	var results MetaDataRecord
	rurl := fmt.Sprintf("%s/meta/%s", _oreConfig.Services.MetaDataURL, site)
	if verbose > 0 {
		fmt.Println("HTTP GET", rurl)
	}
	resp, err := http.Get(rurl)
	if err != nil {
		log.Println("ERROR:", err)
		return results
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&results); err != nil {
		log.Println("ERROR:", err)
		return results
	}
	return results
}

func getMeta(site string) ([]MetaData, error) {
	var records []MetaData
	sites, err := getSites()
	if err != nil {
		return records, err
	}
	for _, sobj := range sites {
		if site == sobj.Name || site == "" {
			if verbose > 0 {
				fmt.Printf("processing %+v\n", sobj)
			}
			rec := metadata(site)
			if rec.Status == "ok" {
				for _, r := range rec.Data {
					records = append(records, r)
				}
			} else {
				fmt.Printf("WARNING: failed metadata record %+v\n", rec)
			}
		}
	}
	return records, nil
}

// helper function to provide usage of meta option
func metaUsage() {
	fmt.Println("orecast meta <ls|add|rm> [value]")
}

// helper function to add meta data record
func metaAddRecord(args []string) {
	fmt.Printf("addRecord with %+v", args)
}

// helper function to delete meta-data record
func metaDeleteRecord(args []string) {
	if len(args) != 2 {
		metaUsage()
		os.Exit(1)
	}
	mid := args[1]
	token, err := accessToken()
	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
	rurl := fmt.Sprintf("%s/meta/%s", _oreConfig.Services.MetaDataURL, mid)
	req, err := http.NewRequest("DELETE", rurl, nil)
	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("ERROR", err, "response body", string(body))
		os.Exit(1)
	}
	if response.Status == "ok" {
		fmt.Printf("SUCCESS: record %s was successfully removed\n", mid)
	} else {
		fmt.Printf("WARNING: record %s failed to be removed\n", mid)
	}

}

// helper funtion to list meta-data records
func metaListRecord(site string) {
	records, err := getMeta(site)
	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
	for _, r := range records {
		fmt.Println("---")
		fmt.Printf("ID         : %s\n", r.ID)
		fmt.Printf("Tags       : %v\n", r.Tags)
		fmt.Printf("Bucket     : %v\n", r.Bucket)
		fmt.Printf("Description: %s\n", r.Description)
	}
}

func metaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "meta",
		Short: "OreCast meta command",
		Long: `OreCast meta command
	Complete documentation is available at https://orecast.com/documentation/`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				metaUsage()
			} else if args[0] == "ls" {
				if len(args) == 2 {
					metaListRecord(args[1])
				} else {
					metaListRecord("")
				}
			} else if args[0] == "add" {
				metaAddRecord(args)
			} else if args[0] == "rm" {
				metaDeleteRecord(args)
			} else {
				fmt.Printf("WARNING: unsupported option(s) %+v", args)
			}
		},
	}
	cmd.SetUsageFunc(func(*cobra.Command) error {
		metaUsage()
		return nil
	})
	return cmd
}
