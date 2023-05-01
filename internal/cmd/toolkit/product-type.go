package toolkit

import (
	"log"
	"path/filepath"

	"github.com/rafaelbeecker/mwskit/internal/mws"
	"github.com/spf13/cobra"
)

// NewRootCmd
func newProductTypeSchemaDownloader() *cobra.Command {
	cmd := &cobra.Command{
		Use: "download",
		RunE: func(cmd *cobra.Command, args []string) error {
			product, _ := cmd.Flags().GetString("product-type")
			target, _ := cmd.Flags().GetString("download-target")
			marketplace, _ := cmd.Flags().GetString("marketplace")

			s := mws.BrowseNodeService{}
			link, err := s.GetProductTypeDefSchemaUrl(marketplace, product)
			if err != nil {
				return err
			}

			dest := filepath.Join(target, product+".json")
			if err := s.DownloadProductTypeDef(dest, link); err != nil {
				return err
			}
			log.Printf("schema downloaded at %s\n", dest)
			return nil
		},
	}

	cmd.Flags().String("marketplace", "A3BLNL8STV6IGJ", "marketplace")
	cmd.Flags().String("product-type", "", "product type")
	cmd.Flags().String("download-target", "", "download target")

	cmd.MarkFlagRequired("product-type")
	cmd.MarkFlagRequired("download-target")
	return cmd
}

// NewRootCmd
func newProductTypeCmd() *cobra.Command {
	root := &cobra.Command{
		Use: "product-type",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
	root.AddCommand(newProductTypeSchemaDownloader())
	return root
}
