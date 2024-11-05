package main

import (
	"fmt"
	"log"

	"github.com/hyperquack/exposed/sdk"
	"github.com/spf13/cobra"
)

var client *sdk.CognitoClient

func main() {
	rootCmd := &cobra.Command{
		Use:   "exposed",
		Short: "A command-line interface for the Exposed API",
	}

	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(addTargetCmd, removeTargetCmd, getTargetsCmd, readCmd, notifyCmd)

	if err := initializeClient(); err != nil {
		log.Fatalf("%v", err)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func initializeClient() error {
	var err error
	client, err = sdk.Authenticate()
	return err
}

var addTargetCmd = &cobra.Command{
	Use:     "start [host]",
	Short:   "Start monitoring a new target",
	Example: "exposed start example.com",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := client.AddTarget(args[0]); err == nil {
			fmt.Printf("Added target: %+v\n", args[0])
		}
	},
}

var removeTargetCmd = &cobra.Command{
	Use:     "stop [host]",
	Short:   "Stop monitoring an existing target",
	Example: "exposed stop example.com",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := client.RemoveTarget(args[0]); err == nil {
			fmt.Printf("Removed target: %+v\n", args[0])
		}
	},
}

var getTargetsCmd = &cobra.Command{
	Use:     "targets",
	Short:   "List all targets",
	Example: "exposed targets",
	Run: func(cmd *cobra.Command, args []string) {
		if resp, err := client.GetTargets(); err == nil {
			for _, hit := range resp.Hits {
				fmt.Println(hit.Host)
			}
		}
	},
}

var readCmd = &cobra.Command{
	Use:     "read [feed] [target]",
	Short:   "Read feed for a target",
	Long:    "Read data from any available feed: port, domain, login, cve",
	Example: "exposed read port example.com",
	Args:    cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		namespace, host := args[0], ""

		if len(args) > 1 {
			host = args[1]
		}

		if resp, err := client.Pull(namespace, host); err == nil {
			for _, hit := range resp.Hits {
				fmt.Printf("%s//%s\n", hit.Host, hit.Value)
			}
		}
	},
}

var notifyCmd = &cobra.Command{
	Use:     "notify [url]",
	Short:   "Set a callback address for notifications",
	Long:    "A callback URL can be any endpoint that receives a POST request, including Slack or Teams webhooks",
	Example: "exposed notify https://hooks.slack.com/services/T123/B456/7890",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client.SetNotification(args[0])
	},
}
