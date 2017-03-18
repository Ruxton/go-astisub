package astisub

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// Constants
const (
	stringRegexpTime = "[\\d]+\\:[\\d]+\\:[\\d]+\\.[\\d]+"
)

// Vars
var (
	regexpIndex          = regexp.MustCompile("^[\\d]+$")
	regexpTimeBoundaries = regexp.MustCompile(fmt.Sprintf("^%s%s%s", stringRegexpTime, TimeBoundariesSeparator, stringRegexpTime))
)

// FromReaderVTT parses a .vtt content
func FromReaderVTT(i io.Reader) (o *Subtitles, err error) {
	// Init
	var scanner = bufio.NewScanner(i)
	var line string

	// Remove the header
	for scanner.Scan() {
		line = scanner.Text()
		if line != "" && line != "WEBVTT" {
			break
		}
	}

	// Recreate an .srt compliant content
	var c []byte
	var isCommentBlock, hasIndex bool
	var index int
	for scanner.Scan() {
		// Fetch line
		line = scanner.Text()
		if strings.HasPrefix(line, "NOTE ") {
			// This is the start of a comment block
			isCommentBlock = true
		} else if isCommentBlock && line == "" {
			// This is the end of a comment block
			isCommentBlock = false
		} else if !isCommentBlock {
			// This is not a comment block
			var bLine = bytes.TrimSpace([]byte(line))

			// Line contains time boundaries
			var match = regexpTimeBoundaries.Find(bLine)
			if len(match) > 0 {
				// Replace . with ,
				bLine = bytes.Replace(match, bytesPeriod, bytesComma, -1)

				// Previous line is not an index, we need to add it
				index++
				if !hasIndex {
					bLine = append([]byte(strconv.Itoa(index)+"\n"), bLine...)
				}
			}

			// Append content
			c = append(c, bLine...)
			c = append(c, BytesLineSeparator...)

			// Check if line is an index so that if next line is a time boundaries we know whether we need
			// to add an index or not
			hasIndex = regexpIndex.Match(bLine)
		}
	}

	// Create the .srt
	o, err = FromReaderSRT(bytes.NewReader(c))
	return
}
