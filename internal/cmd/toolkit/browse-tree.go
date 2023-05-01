package toolkit

import (
	"flag"
	"fmt"
	"os"

	"github.com/rafaelbeecker/mwskit/internal/mws"
	"github.com/spf13/cobra"
)

func newXmlBrowseTreeReportFlatter() *cobra.Command {
	cmd := &cobra.Command{
		Use: "flat",
		RunE: func(cmd *cobra.Command, args []string) error {
			report, _ := cmd.Flags().GetString("report-file")
			target, _ := cmd.Flags().GetString("download-target")

			flag.Parse()

			if _, err := os.Stat(target); err != nil {
				return fmt.Errorf("open-target: %w", err)
			}

			if _, err := os.Stat(report); err != nil {
				return fmt.Errorf("open-report: %w", err)
			}

			s := mws.BrowseNodeService{}
			l, err := s.Read(report)
			if err != nil {
				return err
			}
			return s.Flat(l, target)
		},
	}

	cmd.Flags().String("report-file", "", "product type")
	cmd.Flags().String("download-target", "", "download target")

	cmd.MarkFlagRequired("report-file")
	cmd.MarkFlagRequired("download-target")
	return cmd
}

func newBrowseTreeCmd() *cobra.Command {
	root := &cobra.Command{
		Use: "browse-tree",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
	root.AddCommand(newXmlBrowseTreeReportFlatter())
	return root
}
