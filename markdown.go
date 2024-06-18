package main

import (
	"fmt"
	"os"

	md "github.com/nao1215/markdown"
)

func buildMarkdown(dir string, content Content) error {
	meta := content.Meta
	name := fmt.Sprintf("%s/%s.md", dir, meta.Title)
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

			b.LF().PlainTextf("- location: [%d]()", heading.Location).LF()

			if !nextNote(i, section.Notes) {
				b.HorizontalRule().LF()
			}
		}
	}

	return b.Build()
}

func nextNote(i int, notes []Note) bool {
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
