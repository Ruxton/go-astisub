package subtitles

import (
	"encoding/xml"
	"io"
	"strings"
	"time"
)

// TTML represents a TTML
type TTML struct {
	Lang      string         `xml:"lang,attr,omitempty"`
	Regions   []TTMLRegion   `xml:"head>layout>region,omitempty"`
	Styles    []TTMLStyle    `xml:"head>styling>style,omitempty"`
	Subtitles []TTMLSubtitle `xml:"body>div>p"`
	XMLName   xml.Name       `xml:"tt"`
}

// TTMLRegions represents a TTML region
type TTMLRegion struct {
	Extent string `xml:"extent,attr,omitempty"`
	ID     string `xml:"id,attr,omitempty"`
	Origin string `xml:"origin,attr,omitempty"`
	Style  string `xml:"style,attr,omitempty"`
	ZIndex string `xml:"zIndex,attr,omitempty"`
}

// TTML Style represents a TTML style
type TTMLStyle struct {
	BackgroundColor string `xml:"backgroundColor,attr,omitempty"`
	Color           string `xml:"color,attr,omitempty"`
	DisplayAlign    string `xml:"displayAlign,attr,omitempty"`
	Extent          string `xml:"extent,attr,omitempty"`
	FontFamily      string `xml:"fontFamily,attr,omitempty"`
	FontSize        string `xml:"fontSize,attr,omitempty"`
	ID              string `xml:"id,attr,omitempty"`
	Origin          string `xml:"origin,attr,omitempty"`
	Style           string `xml:"style,attr,omitempty"`
	TextAlign       string `xml:"textAlign,attr,omitempty"`
}

// TTMLDuration represents a TTML duration
type TTMLDuration struct {
	time.Duration
}

// MarshalText allows TTMLDuration to implement the TextMarshaler interface
func (t *TTMLDuration) MarshalText() (text []byte, err error) {
	text = []byte(strings.Replace(FormatDurationSRT(t.Duration), ",", ".", -1))
	return
}

// UnmarshalText allows TTMLDuration to implement the TextUnmarshaler interface
func (t *TTMLDuration) UnmarshalText(text []byte) (err error) {
	t.Duration, err = ParseDurationSRT(strings.Replace(string(text), ".", ",", -1))
	return
}

// TTMLSubtitle represents a TTML subtitle
type TTMLSubtitle struct {
	Begin  TTMLDuration `xml:"begin,attr"`
	End    TTMLDuration `xml:"end,attr"`
	ID     string       `xml:"id,attr,omitempty"`
	Region string       `xml:"region,attr,omitempty"`
	Text   []TTMLText   `xml:"span"`
}

// TTMLText represents a TTML text
type TTMLText struct {
	Style    string `xml:"style,attr,omitempty"`
	Sentence string `xml:",chardata"`
}

// FromReaderTTML parses a .ttml content
func FromReaderTTML(i io.Reader) (o *Subtitles, err error) {
	// Init
	o = &Subtitles{}

	// Unmarshal XML
	var ttml TTML
	if err = xml.NewDecoder(i).Decode(&ttml); err != nil {
		return
	}

	// Loop through subtitles
	for _, s := range ttml.Subtitles {
		// Get text
		var text []string
		for _, t := range s.Text {
			text = append(text, t.Sentence)
		}

		// Append subtitle
		*o = append(*o, &Subtitle{
			EndAt:   s.End.Duration,
			StartAt: s.Begin.Duration,
			Text:    text,
		})
	}
	return
}

// ToWriterTTML formats subtitles as .ttml format into a writer
func (s Subtitles) ToWriterTTML(o io.Writer) (err error) {
	// Do not write anything if no subtitles
	if len(s) == 0 {
		err = ErrNoSubtitlesToWrite
		return
	}

	// Init TTML
	var ttml = TTML{}
	for _, sub := range s {
		// Init TTML text
		var text = []TTMLText{}
		for _, t := range sub.Text {
			text = append(text, TTMLText{Sentence: t})
		}

		// Append subtitle
		ttml.Subtitles = append(ttml.Subtitles, TTMLSubtitle{
			Begin: TTMLDuration{Duration: sub.StartAt},
			End:   TTMLDuration{Duration: sub.EndAt},
			Text:  text,
		})
	}

	// Marshal XML
	err = xml.NewEncoder(o).Encode(ttml)
	return
}