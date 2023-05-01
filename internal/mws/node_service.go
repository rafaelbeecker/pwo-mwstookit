package mws

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/rafaelbeecker/mwskit/internal/mws/signer"
	"golang.org/x/sync/errgroup"
)

type BrowseNodeService struct{}

func (b *BrowseNodeService) Read(p string) (*BrowseList, error) {
	file, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("ls: %w", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("rd: %w", err)
	}

	var l BrowseList
	if err := xml.Unmarshal(data, &l); err != nil {
		return nil, fmt.Errorf("um: %w", err)
	}
	return &l, nil
}

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

	var eg errgroup.Group
	eg.SetLimit(len(d))

	for i, v := range d {
		eg.Go(func(k string, s BrowseList) func() error {
			return func() error {
				log.Printf("Writing %s...\n", k)

				d, err := xml.MarshalIndent(s, "", "  ")
				if err != nil {
					return fmt.Errorf("xml:%w", err)
				}

				if err := os.WriteFile(
					filepath.Join(target, k+".xml"),
					d,
					0644,
				); err != nil {
					return fmt.Errorf("xml:%w", err)
				}

				t := s.Result[0].BrowseNodeName + " (" + k + ")"
				if _, err := p.WriteString(t + "\n"); err != nil {
					return fmt.Errorf("xml:%w", err)
				}
				return nil
			}
		}(i, v))
	}
	return eg.Wait()
}

func (s *BrowseNodeService) GetProductTypeDefSchemaUrl(sellerId string, productType string) (string, error) {
	url := `https://sellingpartnerapi-na.amazon.com/definitions/2020-09-01/productTypes/` + productType
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("GetProductTypeDefSchemaUrl: %w", err)
	}

	q := req.URL.Query()
	q.Add("sellerId", sellerId)
	q.Add("marketplaceIds", "A2Q3Y263D00KWC")
	q.Add("requirements", "LISTING")
	q.Add("locale", "pt_BR")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("host", "sellingpartnerapi-na.amazon.com")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-amz-access-token", os.Getenv("AWS_ACCESS_TOKEN"))
	req.Header.Set("user-agent", "App 1.0 (Language=Golang/1.18);")

	req2 := signer.Sign4(req, signer.Credentials{
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		SecurityToken:   os.Getenv("AWS_SESSION_TOKEN"),
		Region:          "us-east-1",
		Service:         "execute-api",
	})

	client := http.Client{}
	resp, err := client.Do(req2)
	if err != nil {
		return "", fmt.Errorf("GetProductTypeDefSchemaUrl: %w", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("GetProductTypeDefSchemaUrl: %w", err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("GetProductTypeDefSchemaUrl: %d", resp.StatusCode)
	}

	payload := ProductTypeDefinitions{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return "", fmt.Errorf("GetProductTypeDefSchemaUrl: %w", err)
	}
	return payload.Schema.Link.Resource, nil
}

func (s *BrowseNodeService) DownloadProductTypeDef(dest string, link string) error {
	request, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return fmt.Errorf("DownloadProductTypeDef: %w", err)
	}

	file, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("DownloadProductTypeDef: %w", err)
	}

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("DownloadProductTypeDef: %w", err)
	}
	defer resp.Body.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("DownloadProductTypeDef: %w", err)
	}
	return nil
}

// DownloadBatchTypeDef
func (s *BrowseNodeService) DownloadBatchTypeDef(marketplace string, productList string, target string) error {
	file, err := os.Open(productList)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		return err
	}

	var eg errgroup.Group
	eg.SetLimit(5)

	for _, v := range data {
		eg.Go(func(t string) func() error {
			return func() error {
				dest := filepath.Join(target, t+".json")
				f, err := os.Stat(dest)
				if f != nil {
					log.Printf("schema already exists %s\n", dest)
					return nil
				} else if !errors.Is(err, os.ErrNotExist) {
					return err
				}
				log.Printf("downloading schema %s\n", t)
				link, err := s.GetProductTypeDefSchemaUrl(marketplace, t)
				if err != nil {
					return err
				}
				if err := s.DownloadProductTypeDef(dest, link); err != nil {
					return err
				}
				log.Printf("schema downloaded at %s\n", dest)
				return nil
			}
		}(v[0]))
	}
	return eg.Wait()
}
