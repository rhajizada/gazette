package opml

import (
	"encoding/xml"
)

// OPML is the root element.
type OPML struct {
	XMLName xml.Name `xml:"opml"`
	Version string   `xml:"version,attr"`
	Head    Head     `xml:"head"`
	Body    Body     `xml:"body"`
}

// Head holds meta information about the document.
type Head struct {
	Title           string `xml:"title"`
	DateCreated     string `xml:"dateCreated,omitempty"`
	DateModified    string `xml:"dateModified,omitempty"`
	OwnerName       string `xml:"ownerName,omitempty"`
	OwnerEmail      string `xml:"ownerEmail,omitempty"`
	OwnerID         string `xml:"ownerId,omitempty"`
	Docs            string `xml:"docs,omitempty"`
	ExpansionState  string `xml:"expansionState,omitempty"`
	VertScrollState string `xml:"vertScrollState,omitempty"`
	WindowTop       string `xml:"windowTop,omitempty"`
	WindowBottom    string `xml:"windowBottom,omitempty"`
	WindowLeft      string `xml:"windowLeft,omitempty"`
	WindowRight     string `xml:"windowRight,omitempty"`
}

// Body contains all top-level outlines.
type Body struct {
	Outlines []Outline `xml:"outline"`
}

// Outline can be either a folder (with child outlines) or a feed.
type Outline struct {
	Text        string    `xml:"text,attr,omitempty"`
	Title       string    `xml:"title,attr,omitempty"`
	Type        string    `xml:"type,attr,omitempty"`
	XMLURL      string    `xml:"xmlUrl,attr,omitempty"`
	HTMLURL     string    `xml:"htmlUrl,attr,omitempty"`
	Category    string    `xml:"category,attr,omitempty"`
	Description string    `xml:"description,attr,omitempty"`
	Language    string    `xml:"language,attr,omitempty"`
	Version     string    `xml:"version,attr,omitempty"`
	Outlines    []Outline `xml:"outline,omitempty"`
}
