package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// StorageRecord represents Storage record returned by datamanagement service
type StorageRecord struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
}

// helper function to provide usage of s3 option
func s3Usage() {
	fmt.Println("orecast s3 <ls|create|delete|upload> [value]")
	fmt.Println("Examples:")
	fmt.Println("\n# create mew nuicket:")
	fmt.Println("orecast s3 create Cornell/bucket")
	fmt.Println("\n# remove bucket or file:")
	fmt.Println("orecast s3 delete Cornell/bucket")
	fmt.Println("\n# upload new file to a bucket:")
	fmt.Println("orecast s3 upload Cornell/bucket file.txt")
	fmt.Println("\n# upload all files from given directory to a bucket:")
	fmt.Println("orecast s3 upload Cornell/bucket someDirectory")
	fmt.Println("\n# list content of s3 storage:")
	fmt.Println("orecast s3 ls Cornell")
	fmt.Println("\n# list specific bucket on s3 storage:")
	fmt.Println("orecast s3 ls Cornell/bucket")
}

// helper function to list content of a bucket on s3 storage
func s3List(args []string) {
	// args contains [ls bucket]
	if len(args) != 2 {
		fmt.Println("ERROR: wrong number of arguments")
		os.Exit(1)
	}
	if args[0] != "ls" {
		fmt.Println("ERROR: wrong action", args)
		os.Exit(1)
	}
	bucketName := args[1]
	fmt.Printf("INFO: list bucket %s", bucketName)

	var results StorageRecord
	rurl := fmt.Sprintf("%s/storage/%s", _oreConfig.Services.DataManagementURL, bucketName)
	if verbose > 0 {
		fmt.Println("HTTP GET", rurl)
	}
	resp, err := http.Get(rurl)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&results); err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
	fmt.Printf("results: %+v", results)
}

// helper function to create new bucket on s3 storage
func s3Create(args []string) {
	// args contains [create bucket]
	if len(args) != 2 {
		fmt.Println("ERROR: wrong number of arguments")
		os.Exit(1)
	}
	if args[0] != "create" {
		fmt.Println("ERROR: wrong action", args)
		os.Exit(1)
	}
	bucketName := args[1]
	fmt.Printf("INFO: create bucket %s", bucketName)
}

// helper function to upload file or directory to bucket on s3 storage
func s3Upload(args []string) {
	// args contains [upload bucket file|dir]
	if len(args) != 3 {
		fmt.Println("ERROR: wrong number of arguments")
		os.Exit(1)
	}
	if args[0] != "upload" {
		fmt.Println("ERROR: wrong action", args)
		os.Exit(1)
	}
	bucketName := args[1]
	fobj := args[2]
	fmt.Printf("INFO: upload %s to bucket %s", fobj, bucketName)
}

// helper function to delete bucket on s3 storage
func s3Delete(args []string) {
	// args contains [delete bucket]
	if len(args) != 2 {
		fmt.Println("ERROR: wrong number of arguments")
		os.Exit(1)
	}
	if args[0] != "delete" {
		fmt.Println("ERROR: wrong action", args)
		os.Exit(1)
	}
	bucketName := args[1]
	fmt.Printf("INFO: delete bucket %s", bucketName)
}

func s3Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "s3",
		Short: "OreCast s3 command",
		Long: `OreCast s3 command
	Complete documentation is available at https://orecast.com/documentation/`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				s3Usage()
			} else if args[0] == "ls" {
				s3List(args)
			} else if args[0] == "create" {
				s3Create(args)
			} else if args[0] == "delete" {
				s3Delete(args)
			} else if args[0] == "upload" {
				s3Upload(args)
			} else {
				fmt.Printf("WARNING: unsupported option(s) %+v", args)
			}
		},
	}
	cmd.SetUsageFunc(func(*cobra.Command) error {
		s3Usage()
		return nil
	})
	return cmd
}
