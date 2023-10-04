package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// Site represents Site object returned from discovery service
type Site struct {
	Name         string `json:"name" form:"name" binding:"required"`
	URL          string `json:"url" form:"url" binding:"required"`
	Endpoint     string `json:"endpoint" form:"endpoint" binding:"required"`
	AccessKey    string `json:"access_key" form:"access_key" binding:"required"`
	AccessSecret string `json:"access_secret" form:"access_secret" binding:"required"`
	UseSSL       bool   `json:"use_ssl" form:"use_ssl"`
	Description  string `json:"description" form:"description"`
}

func getSites() ([]Site, error) {
	var out []Site
	rurl := fmt.Sprintf("%s/sites", _oreConfig.Services.DiscoveryURL)
	if verbose > 0 {
		fmt.Println("HTTP GET", rurl)
	}
	resp, err := http.Get(rurl)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	var results []Site
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&results); err != nil {
		return out, err
	}
	return results, nil
}

// helper function to provide usage of site option
func siteUsage() {
	fmt.Println("orecast site <ls|add|rm> [value]")
}

// helper function to add site data record
func siteAddRecord(args []string) {
	fmt.Printf("addRecord with %+v", args)
}

// helper function to delete site-data record
func siteDeleteRecord(args []string) {
	if len(args) != 2 {
		metaUsage()
		os.Exit(1)
	}
	site := args[1]
	token, err := accessToken()
	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
	rurl := fmt.Sprintf("%s/site/%s", _oreConfig.Services.DiscoveryURL, site)
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
		fmt.Printf("SUCCESS: site %s was successfully removed\n", site)
	} else {
		fmt.Printf("WARNING: site %s failed to be removed\n", site)
	}
}

// helper funciont to list site record(s)
func siteListRecord(site string) {
	if sites, err := getSites(); err == nil {
		for _, s := range sites {
			fmt.Println("---")
			fmt.Printf("Name       : %s\n", s.Name)
			fmt.Printf("URL        : %s\n", s.URL)
			fmt.Printf("Description: %s\n", s.Description)
		}
	}
}

func siteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "site",
		Short: "OreCast site command",
		Long: `OreCast site command
                Complete documentation is available at https://orecast.com/documentation/`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				siteUsage()
			} else if args[0] == "ls" {
				if len(args) == 2 {
					siteListRecord(args[1])
				} else {
					siteListRecord("")
				}
			} else if args[0] == "add" {
				siteAddRecord(args)
			} else if args[0] == "rm" {
				siteDeleteRecord(args)
			} else {
				fmt.Printf("WARNING: unsupported option(s) %+v", args)
			}
		},
	}
	cmd.SetUsageFunc(func(*cobra.Command) error {
		siteUsage()
		return nil
	})
	return cmd
}
