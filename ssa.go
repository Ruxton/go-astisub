package astisub

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// https://www.matroska.org/technical/specs/subtitles/ssa.html
// http://moodub.free.fr/video/ass-specs.doc
// https://en.wikipedia.org/wiki/SubStation_Alpha

// SSA alignment
const (
	ssaAlignmentCentered              = 2
	ssaAlignmentLeft                  = 1
	ssaAlignmentLeftJustifiedTopTitle = 5
	ssaAlignmentMidTitle              = 8
	ssaAlignmentRight                 = 3
	ssaAlignmentTopTitle              = 4
)

// SSA border styles
const (
	ssaBorderStyleOpaqueBox            = 3
	ssaBorderStyleOutlineAndDropShadow = 1
)

// SSA collisions
const (
	ssaCollisionsNormal  = "Normal"
	ssaCollisionsReverse = "Reverse"
)

// SSA event category
const (
	ssaEventCategoryCommand  = "command"
	ssaEventCategoryComment  = "comment"
	ssaEventCategoryDialogue = "dialogue"
	ssaEventCategoryMovie    = "movie"
	ssaEventCategoryPicture  = "picture"
	ssaEventCategorySound    = "sound"
)

// SSA section names
const (
	ssaSectionNameEvents     = "events"
	ssaSectionNameScriptInfo = "script.info"
	ssaSectionNameStyles     = "styles"
)

// SSA wrap style
const (
	ssaWrapStyleEndOfLineWordWrapping                   = "1"
	ssaWrapStyleNoWordWrapping                          = "2"
	ssaWrapStyleSmartWrapping                           = "0"
	ssaWrapStyleSmartWrappingWithLowerLinesGettingWider = "3"
)

// SSA regexp
var ssaRegexpEffect = regexp.MustCompile("\\{[^\\{]+\\}")

// ReadFromSSA parses an .ssa content
func ReadFromSSA(i io.Reader) (o *Subtitles, err error) {
	// Init
	o = NewSubtitles()
	var scanner = bufio.NewScanner(i)
	var si = &ssaScriptInfo{}
	var ss = []*ssaStyle{}
	var es = []*ssaEvent{}

	// Scan
	var line, sectionName string
	var format map[int]string
	for scanner.Scan() {
		// Fetch line
		line = strings.TrimSpace(scanner.Text())

		// Empty line
		if len(line) == 0 {
			continue
		}

		// Section name
		switch strings.ToLower(line) {
		case "[events]":
			sectionName = ssaSectionNameEvents
			format = make(map[int]string)
			continue
		case "[script info]":
			sectionName = ssaSectionNameScriptInfo
			continue
		case "[v4 styles]", "[v4+ styles]", "[v4 styles+]":
			sectionName = ssaSectionNameStyles
			format = make(map[int]string)
			continue
		}

		// Comment
		if len(line) > 0 && line[0] == ';' {
			si.comments = append(si.comments, strings.TrimSpace(line[1:]))
			continue
		}

		// Split on ":"
		var split = strings.Split(line, ":")
		if len(split) < 2 {
			err = fmt.Errorf("line '%s' should contain at least one ':'", line)
			return
		}
		var header = strings.ToLower(strings.TrimSpace(split[0]))
		var content = strings.TrimSpace(strings.Join(split[1:], ":"))

		// Switch on section name
		switch sectionName {
		case ssaSectionNameScriptInfo:
			if err = si.parse(header, content); err != nil {
				err = errors.Wrap(err, "parsing script info block failed")
				return
			}
		case ssaSectionNameEvents, ssaSectionNameStyles:
			// Parse format
			if header == "format" {
				for idx, item := range strings.Split(content, ",") {
					format[idx] = strings.TrimSpace(strings.ToLower(item))
				}
			} else {
				// No format provided
				if len(format) == 0 {
					err = fmt.Errorf("no %s format provided", sectionName)
					return
				}

				// Switch on section name
				switch sectionName {
				case ssaSectionNameEvents:
					var e *ssaEvent
					if e, err = newSSAEvent(header, content, format); err != nil {
						err = errors.Wrap(err, "building new ssa event failed")
						return
					}
					es = append(es, e)
				case ssaSectionNameStyles:
					var s *ssaStyle
					if s, err = newSSAStyle(content, format); err != nil {
						err = errors.Wrap(err, "building new ssa style failed")
						return
					}
					ss = append(ss, s)
				}
			}
		}
	}

	// Set metadata
	o.Metadata = &Metadata{
		Comments:  si.comments,
		Copyright: si.originalEditing,
		Title:     si.title,
	}

	// Loop through styles
	for _, s := range ss {
		var st = s.style()
		o.Styles[st.ID] = st
	}

	// Loop through events
	for _, e := range es {
		// Only process dialogues
		if e.category == ssaEventCategoryDialogue {
			// Init item
			var item = &Item{
				EndAt: e.end,
				InlineStyle: &StyleAttributes{
					SSAEffect:         e.effect,
					SSAMarginLeft:     e.marginLeft,
					SSAMarginRight:    e.marginRight,
					SSAMarginVertical: e.marginVertical,
				},
				StartAt: e.start,
			}

			// Set style
			if len(e.style) > 0 {
				var ok bool
				if item.Style, ok = o.Styles[e.style]; !ok {
					err = fmt.Errorf("style %s not found", e.style)
					return
				}
			}

			// Loop through lines
			for _, s := range strings.Split(e.text, "\\n") {
				// Init
				s = strings.TrimSpace(s)
				var l = Line{}

				// Extract effects
				var matches = ssaRegexpEffect.FindAllStringIndex(s, -1)
				if len(matches) > 0 {
					// Loop through matches
					var lineItem *LineItem
					var previousEffectEndOffset int
					for _, idxs := range matches {
						if lineItem != nil {
							lineItem.Text = s[previousEffectEndOffset:idxs[0]]
							l = append(l, *lineItem)
						}
						previousEffectEndOffset = idxs[1]
						lineItem = &LineItem{InlineStyle: &StyleAttributes{SSAEffect: s[idxs[0]:idxs[1]]}}
					}
					lineItem.Text = s[previousEffectEndOffset:]
					l = append(l, *lineItem)
				} else {
					l = append(l, LineItem{Text: s})
				}

				// Add line
				item.Lines = append(item.Lines, l)
			}

			// Add item
			o.Items = append(o.Items, item)
		}
	}
	return
}

// ssaScriptInfo represents an SSA script info block
type ssaScriptInfo struct {
	collisions          string
	comments            []string
	originalEditing     string
	originalScript      string
	originalTiming      string
	originalTranslation string
	playDepth           string
	playResX, playResY  int
	scriptType          string
	scriptUpdatedBy     string
	synchPoint          string
	timer               float64
	title               string
	updateDetails       string
	wrapStyle           string
}

// parse parses a script info header/content
func (b *ssaScriptInfo) parse(header, content string) (err error) {
	switch header {
	case "collisions":
		b.collisions = content
	case "original editing":
		b.originalEditing = content
	case "original script":
		b.originalScript = content
	case "original timing":
		b.originalTiming = content
	case "original translation":
		b.originalTranslation = content
	case "playdepth":
		b.playDepth = content
	case "playresx":
		if b.playResX, err = strconv.Atoi(content); err != nil {
			err = errors.Wrapf(err, "atoi of %s failed", content)
		}
	case "playresy":
		if b.playResY, err = strconv.Atoi(content); err != nil {
			err = errors.Wrapf(err, "atoi of %s failed", content)
		}
	case "scripttype":
		b.scriptType = content
	case "script updated by":
		b.scriptUpdatedBy = content
	case "synch point":
		b.synchPoint = content
	case "timer":
		if b.timer, err = strconv.ParseFloat(strings.Replace(content, ",", ".", -1), 64); err != nil {
			err = errors.Wrapf(err, "parseFloat of %s failed", content)
		}
	case "title":
		b.title = content
	case "update details":
		b.updateDetails = content
	case "wrapstyle":
		b.wrapStyle = content
	}
	return
}

// ssaStyle represents an SSA style
type ssaStyle struct {
	alignment       int
	alphaLevel      float64
	angle           float64 // degrees
	backColour      *Color
	bold            bool
	borderStyle     int
	encoding        int
	fontName        string
	fontSize        float64
	italic          bool
	outline         int // pixels
	outlineColour   *Color
	marginLeft      int // pixels
	marginRight     int // pixels
	marginVertical  int // pixels
	name            string
	primaryColour   *Color
	scaleX          float64 // %
	scaleY          float64 // %
	secondaryColour *Color
	shadow          int // pixels
	spacing         int // pixels
	strikeout       bool
	underline       bool
}

// newSSAStyle builds a new SSA style based on an input string and a format
func newSSAStyle(content string, format map[int]string) (s *ssaStyle, err error) {
	// Split content
	var items = strings.Split(content, ",")

	// Not enough items
	if len(items) < len(format) {
		err = fmt.Errorf("content has %d items whereas style format has %d items", len(items), len(format))
		return
	}

	// Loop through items
	s = &ssaStyle{}
	for idx, item := range items {
		// Index not found in format
		var attr string
		var ok bool
		if attr, ok = format[idx]; !ok {
			err = fmt.Errorf("index %d not found in style format %+v", idx, format)
			return
		}

		// Switch on attribute name
		switch attr {
		// Bool
		case "bold", "italic", "strikeout", "underline":
			var b = item == "-1"
			switch attr {
			case "bold":
				s.bold = b
			case "italic":
				s.italic = b
			case "strikeout":
				s.strikeout = b
			case "underline":
				s.underline = b
			}
		// Color
		case "primarycolour", "secondarycolour", "tertiarycolour", "outlinecolour", "backcolour":
			// Build color
			var c *Color
			if c, err = newColorFromSSAColor(item); err != nil {
				err = errors.Wrapf(err, "building new %s from ssa color %s failed", attr, item)
				return
			}

			// Set color
			switch attr {
			case "backcolour":
				s.backColour = c
			case "primarycolour":
				s.primaryColour = c
			case "secondarycolour":
				s.secondaryColour = c
			case "tertiarycolour", "outlinecolour":
				s.outlineColour = c
			}
		// Float
		case "alphalevel", "angle", "fontsize", "scalex", "scaley":
			// Parse float
			var f float64
			if f, err = strconv.ParseFloat(item, 64); err != nil {
				err = errors.Wrapf(err, "parsing float %s failed", item)
				return
			}

			// Set float
			switch attr {
			case "alphalevel":
				s.alphaLevel = f
			case "angle":
				s.angle = f
			case "fontsize":
				s.fontSize = f
			case "scalex":
				s.scaleX = f
			case "scaley":
				s.scaleY = f
			}
		// Int
		case "alignment", "borderstyle", "encoding", "marginl", "marginr", "marginv", "outline", "shadow", "spacing":
			// Parse int
			var i int
			if i, err = strconv.Atoi(item); err != nil {
				err = errors.Wrapf(err, "atoi of %s failed", item)
				return
			}

			// Set int
			switch attr {
			case "alignment":
				s.alignment = i
			case "borderstyle":
				s.borderStyle = i
			case "encoding":
				s.encoding = i
			case "marginl":
				s.marginLeft = i
			case "marginr":
				s.marginRight = i
			case "marginv":
				s.marginVertical = i
			case "outline":
				s.outline = i
			case "shadow":
				s.shadow = i
			case "spacing":
				s.spacing = i
			}
		// String
		case "fontname", "name":
			switch attr {
			case "fontname":
				s.fontName = item
			case "name":
				s.name = item
			}
		}
	}
	return
}

// style converts ssaStyle to Style
func (s *ssaStyle) style() *Style {
	return &Style{
		ID: s.name,
		InlineStyle: &StyleAttributes{
			SSAAlignment:       s.alignment,
			SSAAlphaLevel:      s.alphaLevel,
			SSAAngle:           s.angle,
			SSABackColour:      s.backColour,
			SSABold:            s.bold,
			SSABorderStyle:     s.borderStyle,
			SSAEncoding:        s.encoding,
			SSAFontName:        s.fontName,
			SSAFontSize:        s.fontSize,
			SSAItalic:          s.italic,
			SSAOutline:         s.outline,
			SSAOutlineColour:   s.outlineColour,
			SSAMarginLeft:      s.marginLeft,
			SSAMarginRight:     s.marginRight,
			SSAMarginVertical:  s.marginVertical,
			SSAPrimaryColour:   s.primaryColour,
			SSAScaleX:          s.scaleX,
			SSAScaleY:          s.scaleY,
			SSASecondaryColour: s.secondaryColour,
			SSAShadow:          s.shadow,
			SSASpacing:         s.spacing,
			SSAStrikeout:       s.strikeout,
			SSAUnderline:       s.underline,
		},
	}
}

// ssaEvent represents an SSA event
type ssaEvent struct {
	category       string
	effect         string
	end            time.Duration
	layer          int
	marked         string
	marginLeft     int // pixels
	marginRight    int // pixels
	marginVertical int // pixels
	name           string
	start          time.Duration
	style          string
	text           string
}

// newSSAEvent builds a new SSA event based on an input string and a format
func newSSAEvent(header, content string, format map[int]string) (e *ssaEvent, err error) {
	// Split content
	var items = strings.Split(content, ",")

	// Not enough items
	if len(items) < len(format) {
		err = fmt.Errorf("content has %d items whereas style format has %d items", len(items), len(format))
		return
	}

	// Last item may contain commas, therefore we need to fix it
	items[len(format)-1] = strings.Join(items[len(format)-1:], ",")

	// Loop through items
	e = &ssaEvent{category: header}
	for idx, item := range items {
		// Index not found in format
		var attr string
		var ok bool
		if attr, ok = format[idx]; !ok {
			err = fmt.Errorf("index %d not found in event format %+v", idx, format)
			return
		}

		// Switch on attribute name
		switch attr {
		// Duration
		case "start", "end":
			// Parse duration
			var d time.Duration
			if d, err = parseDurationSSA(item); err != nil {
				err = errors.Wrapf(err, "parsing ssa duration %s failed", item)
				return
			}

			// Set duration
			switch attr {
			case "end":
				e.end = d
			case "start":
				e.start = d
			}
		// Int
		case "layer", "marginl", "marginr", "marginv":
			// Parse int
			var i int
			if i, err = strconv.Atoi(item); err != nil {
				err = errors.Wrapf(err, "atoi of %s failed", item)
				return
			}

			// Set int
			switch attr {
			case "layer":
				e.layer = i
			case "marginl":
				e.marginLeft = i
			case "marginr":
				e.marginRight = i
			case "marginv":
				e.marginVertical = i
			}
		// String
		case "effect", "marked", "name", "style", "text":
			switch attr {
			case "effect":
				e.effect = item
			case "marked":
				e.marked = item
			case "name":
				e.name = item
			case "style":
				e.style = item
			case "text":
				e.text = item
			}
		}
	}
	return
}

// newColorFromSSAColor builds a new color based on an SSA color
func newColorFromSSAColor(i string) (c *Color, err error) {
	// Empty
	if len(i) == 0 {
		return
	}

	// Check whether input is decimal or hexadecimal
	var s = i
	var base = 10
	if strings.HasPrefix(i, "&H") {
		s = i[2:]
		base = 16
	}
	return newColorFromString(s, base)
}

// parseDurationSSA parses an .ssa duration
func parseDurationSSA(i string) (time.Duration, error) {
	return parseDuration(i, ".", 3)
}
