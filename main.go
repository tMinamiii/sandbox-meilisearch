package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/meilisearch/meilisearch-go"
	"github.com/urfave/cli/v2"
)

// IndexName is meilisearch index name
const (
	IndexName = "itnews"
	Q         = "q"
	Limit     = "limit"
	Offset    = "offset"
)

type SearchParams struct {
	q      string
	limit  int64
	offset int64
}

var client = func() *meilisearch.Client {
	c := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   "http://localhost",
		APIKey: "MASTER_KEY",
	})
	return c
}()

func main() {
	app := &cli.App{
		Name: "meilisearch-cli",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  Q,
				Usage: "search query",
			},
			&cli.Int64Flag{
				Name:  Limit,
				Usage: "hits limit",
				Value: 20,
			},
			&cli.Int64Flag{
				Name:  Offset,
				Usage: "offset",
				Value: 0,
			},
		},
		Action: func(cCtx *cli.Context) error {
			q := cCtx.String(Q)
			limit := cCtx.Int64(Limit)
			offset := cCtx.Int64(Offset)
			if q == "" {
				return nil
			}
			req := &meilisearch.SearchRequest{
				Limit:  limit,
				Offset: offset,
			}
			return run(q, req)
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(q string, req *meilisearch.SearchRequest) error {
	resp, err := client.Index(IndexName).Search(q, req)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	bjson, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println(string(bjson))

	return nil

}
