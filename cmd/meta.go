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
		if site == sobj.Name {
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
func metaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "meta",
		Short: "OreCast meta command",
		Long: `OreCast meta command
	Complete documentation is available at https://orecast.com/documentation/`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			site := "Cornell"
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
		},
	}
	return cmd
}
