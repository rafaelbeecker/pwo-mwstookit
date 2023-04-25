package mws

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//BrowseNodeService
type BrowseNodeService struct{}

// Read browse node tree report from path
func (b *BrowseNodeService) Read(p string) (*BrowseList, error) {
	file, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("get-node-list: %w", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read-node-list: %w", err)
	}

	var l BrowseList
	if err := xml.Unmarshal(data, &l); err != nil {
		return nil, fmt.Errorf("uml-node-list: %w", err)
	}
	return &l, nil
}

// Flat browse node list into individual files
func (b *BrowseNodeService) Flat(l *BrowseList, target string) error {

	var d = make(map[string]BrowseList)
	for _, v := range l.Result {
		s := strings.Split(v.BrowsePathById, ",")
		if _, ok := d[s[0]]; !ok {
			d[s[0]] = BrowseList{Result: []BrowseNode{v}}
		} else if ok {
			d[s[0]] = BrowseList{Result: append(d[s[0]].Result, v)}
		}
	}

	p, err := os.OpenFile(
		filepath.Join(target, "nodes.csv"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return err
	}
	defer p.Close()

	for k, v := range d {
		log.Printf("Writting %s...\n", k)
		d, err := xml.MarshalIndent(v, "", "  ")
		if err != nil {
			return err
		}

		if err := os.WriteFile(
			filepath.Join(target, k+".xml"),
			d,
			0644,
		); err != nil {
			return fmt.Errorf("write-node-list: %w", err)
		}

		s := k + ";" + v.Result[0].BrowseNodeName + ""
		if _, err := p.WriteString(s + "\n"); err != nil {
			return fmt.Errorf("write-parent-list: %w", err)
		}
	}
	return nil
}
