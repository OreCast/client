package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type DBSRecord map[string]any

// helper function to fetch data from DBS service
func getData(rurl string) []DBSRecord {
	var results []DBSRecord
	if verbose > 0 {
		fmt.Println("HTTP GET", rurl)
	}
	resp, err := http.Get(rurl)
	if err != nil {
		fmt.Println("ERROR:", err)
		return results
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&results); err != nil {
		fmt.Println("ERROR:", err)
		return results
	}
	return results
}

// helper function to print dbs record items
func printResults(rec DBSRecord) {
	fmt.Println("---")
	maxKey := 0
	for key, _ := range rec {
		if len(key) > maxKey {
			maxKey = len(key)
		}
	}
	for key, val := range rec {
		pad := strings.Repeat(" ", maxKey-len(key))
		fmt.Printf("%s%s\t%v\n", key, pad, val)
	}
}

// helper function to list dataset information
func dbsListRecord(args []string) {
	if len(args) == 1 {
		fmt.Println("WARNING: please provide dbs attribute")
		os.Exit(1)
	}
	if args[1] == "datasets" {
		rurl := fmt.Sprintf("%s/datasets", _oreConfig.Services.DataBookkeepingURL)
		for _, rec := range getData(rurl) {
			printResults(rec)
		}
	} else {
		fmt.Println("Not implemented yet")
	}
}

// helper function to add dataset information
func dbsAddRecord(args []string) {
}

// helper function to delete dataset information
func dbsDeleteRecord(args []string) {
}

// helper function to provide usage of dbs option
func dbsUsage() {
	fmt.Println("orecast dbs <ls|add|rm> [value]")
	fmt.Println("Examples:")
	fmt.Println("\n# list all dbs records:")
	fmt.Println("orecast dbs ls <dataset|site|file>")
	fmt.Println("\n# remove dbs-data record:")
	fmt.Println("orecast dbs rm <dataset|site|file>")
	fmt.Println("\n# add dbs-data record:")
	fmt.Println("orecast dbs add <dataset|site|file>")
}
func dbsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dbs",
		Short: "OreCast dbs command",
		Long: `OreCast data-bookkeeping system command
                Complete documentation is available at https://orecast.com/documentation/`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				dbsUsage()
			} else if args[0] == "ls" {
				dbsListRecord(args)
			} else if args[0] == "add" {
				dbsAddRecord(args)
			} else if args[0] == "rm" {
				dbsDeleteRecord(args)
			} else {
				fmt.Printf("WARNING: unsupported option(s) %+v", args)
			}
		},
	}
	cmd.SetUsageFunc(func(*cobra.Command) error {
		dbsUsage()
		return nil
	})
	return cmd
}
