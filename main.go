package main

import (
	"fmt"
	"log"

	"exposed/sdk"
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
	rootCmd.AddCommand(addTargetCmd, removeTargetCmd, getTargetsCmd, pushCmd, pullCmd, notifyCmd)

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

var pushCmd = &cobra.Command{
	Use:     "push [feed] [target] [value]",
	Short:   "Push data into a feed",
	Long:    "Push data into any available feed: port, domain, login, cve",
	Example: "exposed push port example.com 80",
	Args:    cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		namespace, host, value := args[0], args[1], args[2]
		if _, err := client.Push(namespace, host, value); err == nil {
			fmt.Printf("Saved data")
		}
	},
}

var pullCmd = &cobra.Command{
	Use:     "pull [feed] [target]",
	Short:   "Pull feed for a target",
	Long:    "Pull data from any available feed: port, domain, login, cve",
	Example: "exposed pull port example.com",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		namespace, host := args[0], args[1]
		if resp, err := client.Pull(namespace, host); err == nil {
			for _, hit := range resp.Hits {
				fmt.Println(hit.Value)
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
