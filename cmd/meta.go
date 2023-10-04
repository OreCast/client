package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

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

// helper function to add meta data record
func addRecord(args []string) {
	fmt.Printf("addRecord with %+v", args)
}

// helper function to delete meta-data record
func deleteRecord(args []string) {
	fmt.Printf("deleteRecord with %+v", args)
}

// helper funtion to list meta-data records
func listRecords(site string) {
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

// helper function to provide usage of meta option
func usage() {
	fmt.Println("orecast meta <ls|add|rm> [value]")
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
				usage()
			} else if args[0] == "ls" {
				if len(args) == 2 {
					listRecords(args[1])
				} else {
					listRecords("")
				}
			} else if args[0] == "add" {
				addRecord(args)
			} else if args[0] == "rm" {
				deleteRecord(args)
			} else {
				fmt.Printf("WARNING: unsupported option(s) %+v", args)
			}
		},
	}
	cmd.SetUsageFunc(func(*cobra.Command) error {
		usage()
		return nil
	})
	return cmd
}
