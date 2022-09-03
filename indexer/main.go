package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/meilisearch/meilisearch-go"
)

func main() {
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   "http://127.0.0.1",
		APIKey: "MASTER_KEY",
	})
	// An index is where the documents are stored.
	index := client.Index("itnews")

	var files []string
	root := "../news"

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() && filepath.Ext(path) == ".csv" {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
	id := int64(1)
	for _, file := range files {
		fp, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
			continue
		}
		defer fp.Close()

		reader := csv.NewReader(fp)
		docs := make([]map[string]any, 0, len(files))
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
				continue
			}
			doc := map[string]any{
				"id":             id,
				"category":       record[0],
				"title":          record[1],
				"manuscript_len": record[2],
				"manuscript":     record[3],
			}
			docs = append(docs, doc)
			id++
		}
		index.AddDocuments(docs)
	}
}
