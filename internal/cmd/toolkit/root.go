package toolkit

import (
	"github.com/spf13/cobra"
)

// NewRootCmd
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use: "tookit",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
	root.AddCommand(newBrowseTreeCmd())
	root.AddCommand(newProductTypeCmd())
	return root
}

// // Run
// func Run() error {
// 	// var report = flag.String("report", "", "xml browse tree report file name")
// 	// var target = flag.String("target", "", "flat node file output directory")

// 	// flag.Parse()

// 	// if _, err := os.Stat(string(*target)); err != nil {
// 	// 	return fmt.Errorf("open-target: %w", err)
// 	// }

// 	// if _, err := os.Stat(string(*report)); err != nil {
// 	// 	return fmt.Errorf("open-report: %w", err)
// 	// }

// 	s := mws.BrowseNodeService{}
// 	// _, err := s.Read(string(*report))
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	link, err := s.GetProductTypeDefSchemaUrl("A3BLNL8STV6IGJ", "SHAMPOO")
// 	if err != nil {
// 		return nil
// 	}

// 	wd, _ := os.Getwd()
// 	dest := filepath.Join(wd, "resources", "schemas", "SHAMPOO.json")
// 	if err := s.DownloadProductTypeDef(dest, link); err != nil {
// 		return err
// 	}

// 	fmt.Println(link)
// 	return nil
// 	//return s.Flat(l, string(*target))
// }
