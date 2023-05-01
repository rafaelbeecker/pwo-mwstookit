package toolkit

import (
	"log"
	"path/filepath"

	"github.com/rafaelbeecker/mwskit/internal/mws"
	"github.com/spf13/cobra"
)

func newProductTypeSchemaDownloader() *cobra.Command {
	cmd := &cobra.Command{
		Use: "download",
		RunE: func(cmd *cobra.Command, args []string) error {
			productType, _ := cmd.Flags().GetString("product-type")
			productList, _ := cmd.Flags().GetString("product-list")
			target, _ := cmd.Flags().GetString("download-target")
			marketplace, _ := cmd.Flags().GetString("marketplace")

			if productType != "" {
				s := mws.BrowseNodeService{}
				link, err := s.GetProductTypeDefSchemaUrl(marketplace, productType)
				if err != nil {
					return err
				}
				dest := filepath.Join(target, productType+".json")
				if err := s.DownloadProductTypeDef(dest, link); err != nil {
					return err
				}
				log.Printf("schema downloaded at %s\n", dest)
			} else if productList != "" {
				s := mws.BrowseNodeService{}
				if err := s.DownloadBatchTypeDef(
					marketplace,
					productList,
					target,
				); err != nil {
					return err
				}
				log.Println("batch downloaded successfuly")
			}
			return nil
		},
	}

	cmd.Flags().String("marketplace", "A3BLNL8STV6IGJ", "marketplace")
	cmd.Flags().String("product-type", "", "product type")
	cmd.Flags().String("download-target", "", "download target")
	cmd.Flags().String("product-list", "", "product type (csv)")
	cmd.MarkFlagRequired("download-target")
	return cmd
}

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
