package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	md "github.com/nao1215/markdown"
)

func buildMarkdown(dir string, content Content) (err error) {
	meta := content.Meta

	name := fmt.Sprintf("%s/%s.md", dir, meta.Title)

	blocks, err := originFileBlocks(name)
	if err != nil {
		return err
	}

	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b := md.NewMarkdown(f).
		// b := md.NewMarkdown(os.Stdout).
		H1(meta.Title).
		H2("MetaData").
		BulletList("Author:"+meta.Authors, "Cication:"+meta.Citation)

	b.H2("Highlights").LF()

	sections := content.Sections
	for _, section := range sections {
		b.H3(section.Text).LF()
		subTitle := ""

		for i, note := range section.Notes {
			heading := note.Heading
			if subTitle != heading.SubTitle {
				subTitle = heading.SubTitle
				b.H4(subTitle).LF()
			}

			if heading.Type == MarkHeadingType {
				b.Blockquote(note.Text)
			} else {
				b.PlainTextf("> [!note]   %s", note.Text)
			}

			if blocks != nil {
				if block, ok := blocks[heading.Location]; ok {
					b.PlainText(block).LF()
				}
			}

			b.LF().PlainTextf("- location: [%d]()", heading.Location).LF()

			if !hasNote(i, section.Notes) {
				b.HorizontalRule().LF()
			}
		}
	}

	return b.Build()
}

func hasNote(i int, notes []Note) bool {
	heading := notes[i].Heading

	if heading.Type == NoteHeadingType {
		return false
	}

	if i >= len(notes)-1 {
		return false
	}

	if notes[i+1].Heading.Type == NoteHeadingType {

		return true
	}
	return false
}
func originFileBlocks(name string) (blocks map[int]string, err error) {
	exists, err := PathExists(name)
	if err != nil {
		return nil, fmt.Errorf("name, PathExists:%w", name, err)
	}
	if !exists {
		return nil, fmt.Errorf("name not exists", name)
	}

	f, err := os.Open(name)
	if err != nil {
		return
	}
	defer f.Close()

	br := bufio.NewReader(f)
	block := ""
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		line := string(a)

		// 判断是否存在
		if strings.HasPrefix(line, "^") {
			block = line
		}

		// 查找下一个location
		if block != "" {
			if strings.HasPrefix(line, "- location") {
				valid := regexp.MustCompile("[0-9]+")
				location := valid.FindAllStringSubmatch(line, -1)
				l, _ := strconv.Atoi(location[0][0])

				if blocks == nil {
					blocks = make(map[int]string)
				}

				blocks[l] = block
				block = ""
			}
		}
	}

	return
}
