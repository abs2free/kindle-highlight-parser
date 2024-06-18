package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type MetaData struct {
	Title    string
	Authors  string
	Citation string
}

type HeadingType int

const (
	MarkHeadingType HeadingType = iota
	NoteHeadingType
)

type Content struct {
	Meta     MetaData
	Sections []Section
}

type Section struct {
	Text  string
	Notes []Note
}

type Note struct {
	Heading Heading
	Text    string
}

type Heading struct {
	Type     HeadingType
	SubTitle string
	Location int
}

func parseHtml(doc *goquery.Document) (content Content, err error) {
	content.Meta = parseMetaData(doc)
	content.Sections, err = parseContent(doc)
	return
}

func parseMetaData(doc *goquery.Document) MetaData {
	meta := MetaData{}

	meta.Title = strings.TrimSpace(doc.Find(".bookTitle").Text())
	meta.Authors = strings.TrimSpace(doc.Find(".authors").Text())
	meta.Citation = strings.TrimSpace(doc.Find(".citation").Text())
	return meta
}

func parseContent(doc *goquery.Document) (sections []Section, err error) {

	doc.Find(".sectionHeading").Each(func(i int, h *goquery.Selection) {

		var notes []Note
		var heading Heading

		h.NextAll().Each(func(i int, s *goquery.Selection) {
			// noteHeading
			class, exists := s.Attr("class")
			if !exists {
				return
			}
			if class == "noteHeading" {
				var e error
				heading, e = parseHeading(s.Text())
				if e != nil {
					err = fmt.Errorf("can not parse Heading:%w,from text:%s", e, s.Text())
					return
				}
			} else { // noteText
				note := Note{
					Heading: heading,
					Text:    strings.TrimSpace(s.Text()),
				}
				heading = Heading{}
				notes = append(notes, note)
			}
		})

		section := Section{
			Text:  strings.TrimSpace(h.Text()),
			Notes: notes,
		}
		sections = append(sections, section)
	})

	return
}

func parseHeading(h string) (heading Heading, err error) {
	h = strings.TrimSpace(h)

	s := strings.Split(h, "-")
	st := strings.TrimSpace(s[0])
	if strings.HasPrefix(st, "标注") {
		heading.Type = MarkHeadingType
	} else {
		heading.Type = NoteHeadingType
	}

	ss := strings.Split(s[1], ">")
	heading.SubTitle = strings.TrimSpace(ss[0])

	location := strings.TrimSpace(ss[1])
	location = strings.TrimLeft(location, "位置 ")
	l, err := strconv.Atoi(location)
	if err != nil {
		return Heading{}, err
	}
	heading.Location = l

	return
}
