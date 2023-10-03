package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

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

func siteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "site",
		Short: "OreCast site command",
		Long: `OreCast site command
                Complete documentation is available at https://orecast.com/documentation/`,
		Run: func(cmd *cobra.Command, args []string) {
			if sites, err := getSites(); err == nil {
				for _, s := range sites {
					fmt.Println("---")
					fmt.Printf("Name       : %s\n", s.Name)
					fmt.Printf("URL        : %s\n", s.URL)
					fmt.Printf("Description: %s\n", s.Description)
				}
			} else {
				fmt.Println("ERROR", err)
			}
		},
	}
	return cmd
}
