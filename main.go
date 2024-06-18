package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/PuerkitoBio/goquery"
)

var (
	inputDir  *string
	outputDir *string
)

func main() {
	inputDir = flag.String("i", "html", "dir path to read from")
	outputDir = flag.String("o", ".", "dir path to write to")
	flag.Parse()

	if ok, err := isDir(*inputDir); !ok {
		log.Fatal(err)
	}
	if ok, err := isDir(*outputDir); !ok {
		log.Fatal(err)
	}

	files, err := ListDirFiles(*inputDir)
	if err != nil {
		log.Fatalf("list dir files is wrong")
	}
	for _, file := range files {
		processFile(*inputDir + file)
	}
}
func isDir(path string) (bool, error) {
	dirinfo, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("dir:%s has a error:%w", path, err)
	}
	if !dirinfo.IsDir() {
		return false, fmt.Errorf("path:%s is not dir", path)
	}
	return true, nil
}

func processFile(file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(bufio.NewReader(f))
	if err != nil {
		log.Fatal(err)
	}

	content, err := parseHtml(doc)
	if err != nil {
		log.Fatalf("parse html error:%v", err)
	}

	err = buildMarkdown(*outputDir, content)
	if err != nil {
		log.Fatalf("build markdown error:%v", err)
	}
}

// ListDirFiles lists all the file or dir names in the specified directory.
// Note that ListDirFiles don't traverse recursively.
func ListDirFiles(dirname string) ([]string, error) {
	infos, err := os.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(infos))
	for i, info := range infos {
		if info.IsDir() {
			continue
		}
		names[i] = info.Name()
	}
	return names, nil
}
