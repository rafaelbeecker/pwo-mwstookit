package toolkit

import (
	"flag"
	"fmt"
	"os"

	"github.com/rafaelbeecker/mwskit/internal/mws"
)

// Run
func Run() error {
	var report = flag.String("report", "", "xml browse tree report file name")
	var target = flag.String("target", "", "flat node file output directory")

	flag.Parse()

	if _, err := os.Stat(string(*target)); err != nil {
		return fmt.Errorf("open-target: %w", err)
	}

	if _, err := os.Stat(string(*report)); err != nil {
		return fmt.Errorf("open-report: %w", err)
	}

	s := mws.BrowseNodeService{}
	l, err := s.Read(string(*report))
	if err != nil {
		return err
	}
	return s.Flat(l, string(*target))
}
