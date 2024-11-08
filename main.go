package main

import (
	"fmt"
	"os"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/spf13/cobra"
)

func main() {
	clientFunc := func(hostname string) (*api.RESTClient, error) {
		// If user has overridden the hostname, then use that.
		if hostname != "" {
			fmt.Println("Using hostname override")
			return api.NewRESTClient(api.ClientOptions{Host: hostname})
		}

		// If we're not within a GitHub repo, then use the default host
		repo, err := repository.Current()
		if err != nil {
			fmt.Println("Using default GitHub host")
			return api.DefaultRESTClient()
		}

		// Otherwise, use the repository host
		fmt.Println("Using repository host")
		return api.NewRESTClient(api.ClientOptions{Host: repo.Host})
	}

	cmd := newRootCmd(clientFunc)

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newRootCmd(clientFunc func(string) (*api.RESTClient, error)) *cobra.Command {
	var hostnameOverride string

	cmd := &cobra.Command{
		Use:   "sandbox",
		Short: "Retrieve /user information",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFunc(hostnameOverride)
			if err != nil {
				return err
			}

			response := &struct {
				Login string `json:"login"`
				URL   string `json:"url"`
			}{}
			client.Get("user", response)

			fmt.Printf("Login: %s\nURL: %s\n\n", response.Login, response.URL)

			return nil
		},
	}

	cmd.Flags().StringVarP(&hostnameOverride, "hostname", "H", "", "override the GitHub hostname to use")

	return cmd
}
