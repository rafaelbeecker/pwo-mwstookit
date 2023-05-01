package toolkit

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use: "tookit",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
	root.SetHelpCommand(&cobra.Command{Hidden: true})
	root.AddCommand(newBrowseTreeCmd())
	root.AddCommand(newProductTypeCmd())
	return root
}
