package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type MetaData struct {
	Title    string
	Authors  string
	Citation string
}

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

type HeadingType int

const (
	HeadingTypeHighlight HeadingType = iota
	HeadingTypeNote
)

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

	var st, sub, location string
	if strings.Contains(h, ">") {
		s := strings.Split(h, "-")
		st = strings.TrimSpace(s[0])
		ss := strings.Split(s[1], ">")
		sub = ss[0]
		location = ss[1]
	} else {
		s := strings.Split(h, "-")
		st = strings.TrimSpace(s[0])
		location = s[1]
	}

	if strings.HasPrefix(st, "标注") || strings.HasPrefix(st, "Highlight") {
		heading.Type = HeadingTypeHighlight
	} else {
		heading.Type = HeadingTypeNote
	}

	heading.SubTitle = strings.TrimSpace(sub)
	l, err := extractNumber(strings.TrimSpace(location))
	if err != nil {
		err = fmt.Errorf("extractNumber:%s has error:%w", location, err)
		return
	}
	heading.Location = l

	return
}
