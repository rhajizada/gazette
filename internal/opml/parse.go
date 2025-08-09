package opml

import (
	"bytes"
	"encoding/xml"
	"io"

	"github.com/rhajizada/gazette/internal/repository"
)

// Parse reads an OPML document from r and returns it.
func Parse(r io.Reader) (*OPML, error) {
	dec := xml.NewDecoder(r)
	var doc OPML
	if err := dec.Decode(&doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

// FeedsToOPML creates an OPML document from a slice of feeds.
func FeedsToOPML(title string, feeds []repository.Feed) *OPML {
	outlines := make([]Outline, 0, len(feeds))
	for _, f := range feeds {
		label := ""
		if f.Title != nil && *f.Title != "" {
			label = *f.Title
		} else {
			// fallback to FeedLink or Link if title missing
			if f.Link != nil && *f.Link != "" {
				label = *f.Link
			} else {
				label = f.FeedLink
			}
		}

		o := Outline{
			Text:    label,
			Title:   label,
			Type:    "rss", // default; can override if FeedType set
			XMLURL:  f.FeedLink,
			HTMLURL: firstNonEmpty(f.Link, f.Links),
			Language: func() string {
				if f.Language != nil {
					return *f.Language
				}
				return ""
			}(),
			Description: func() string {
				if f.Description != nil {
					return *f.Description
				}
				return ""
			}(),
			Version: func() string {
				if f.FeedVersion != nil {
					return *f.FeedVersion
				}
				return ""
			}(),
		}
		if f.FeedType != nil && *f.FeedType != "" {
			o.Type = *f.FeedType
		}
		outlines = append(outlines, o)
	}

	return &OPML{
		Version: "1.0",
		Head:    Head{Title: title},
		Body:    Body{Outlines: outlines},
	}
}

// ToXML serializes the document using tab indentation by default.
func (o *OPML) ToXML() (string, error) {
	var buf bytes.Buffer
	buf.WriteString(xml.Header) // <?xml version="1.0" encoding="UTF-8"?>

	enc := xml.NewEncoder(&buf)
	enc.Indent("", "\t")

	if err := enc.Encode(o); err != nil {
		return "", err
	}
	if err := enc.Flush(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func firstNonEmpty(primary *string, alts []string) string {
	if primary != nil && *primary != "" {
		return *primary
	}
	for _, s := range alts {
		if s != "" {
			return s
		}
	}
	return ""
}
