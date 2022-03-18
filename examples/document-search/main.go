package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/vench/word_index"
)

func main() {
	pathDirectory := flag.String("path", "./examples/document-search/documents", "Path directory of documents")
	flag.Parse()

	files, err := ioutil.ReadDir(*pathDirectory)
	if err != nil {
		log.Fatalf("failed to open directory: %v", err)
	}

	documents := make([]string, len(files))
	for i := range files {
		f := files[i]
		log.Printf("read file: %s\n", f.Name())

		filePath := path.Join(*pathDirectory, f.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("failed to read file[%s]: %v", filePath, err)
		}
		documents[i] = string(data)
	}

	s := word_index.NewDocumentSearch(documents...)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		log.Printf("search text: %s\n", text)
		query := word_index.SequenceToFeature(text)
		log.Printf("query: %v \n", query)
		result := s.Find(query...)

		log.Println(result)
	}
}
