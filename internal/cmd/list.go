package cmd

import (
	"fmt"
	"log/slog"

	"github.com/nobe4/gh-not/internal/jq"
	"github.com/spf13/cobra"
)

var (
	filterFlag = ""
	jqFlag     = ""

	listCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List notifications",
		Example: `
  gh-not list
  gh-not list --filter ''.author.login | contains("4")'
`,
		RunE: runList,
	}
)

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&filterFlag, "filter", "f", "", "Filter with a jq expression passed into a select(...) call")
	listCmd.Flags().StringVarP(&jqFlag, "jq", "q", "", "jq expression to run on the notification list")
	listCmd.MarkFlagsMutuallyExclusive("filter", "jq")
}

func runList(cmd *cobra.Command, args []string) error {
	notifications, err := client.Notifications()
	if err != nil {
		slog.Error("Failed to list notifications", "err", err)
		return err
	}

	if filterFlag != "" {
		notificationsList, err := jq.Filter(filterFlag, notifications.ToSlice())
		if err != nil {
			return err
		}
		notifications = notificationsList.ToMap()
	}

	if jqFlag != "" {
		return fmt.Errorf("`gh-not list --jq` implementation needed")
	}

	out, err := notifications.ToTable()
	if err != nil {
		slog.Warn("Failed to generate a table, using toString", "err", err)
		out = notifications.ToString()
	}

	fmt.Println(out)

	return nil
}
