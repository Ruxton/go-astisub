package subtitles

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// Constants
const (
	TimeBoundariesSeparator = " --> "
)

// Vars
var (
	BytesTimeBoundariesSeparator = []byte(TimeBoundariesSeparator)
)

// ParseDurationSRT parses an .srt duration
func ParseDurationSRT(i string) (o time.Duration, err error) {
	// Split milliseconds
	var parts = strings.Split(i, ",")
	if len(parts) != 2 {
		err = fmt.Errorf("No milliseconds detected in %s", i)
		return
	}
	if len(parts[1]) != 3 {
		err = fmt.Errorf("Invalid number of millisecond digits detected in %s", i)
		return
	}

	// Parse milliseconds
	var milliseconds int
	if milliseconds, err = strconv.Atoi(strings.TrimSpace(parts[1])); err != nil {
		return
	}

	// Split hours, minutes and seconds
	parts = strings.Split(strings.TrimSpace(parts[0]), ":")
	if len(parts) != 3 {
		err = fmt.Errorf("No hours, minutes or seconds detected in %s", i)
		return
	}

	// Parse seconds
	var seconds int
	if seconds, err = strconv.Atoi(strings.TrimSpace(parts[2])); err != nil {
		return
	}

	// Parse minutes
	var minutes int
	if minutes, err = strconv.Atoi(strings.TrimSpace(parts[1])); err != nil {
		return
	}

	// Parse hours
	var hours int
	if hours, err = strconv.Atoi(strings.TrimSpace(parts[0])); err != nil {
		return
	}

	// Generate output
	o = time.Duration(milliseconds)*time.Millisecond + time.Duration(seconds)*time.Second + time.Duration(minutes)*time.Minute + time.Duration(hours)*time.Hour
	return
}

// FromReaderSRT parses an .srt content
func FromReaderSRT(i io.Reader) (o *Subtitles, err error) {
	// Init
	o = &Subtitles{}
	var scanner = bufio.NewScanner(i)

	// Scan
	var line string
	var s = &Subtitle{}
	for scanner.Scan() {
		// Fetch line
		line = scanner.Text()

		// Line contains time boundaries
		if strings.Contains(line, TimeBoundariesSeparator) {
			// Remove last item of previous subtitle since it's the index
			s.Text = s.Text[:len(s.Text)-1]

			// Remove trailing empty lines
			if len(s.Text) > 0 {
				for i := len(s.Text) - 1; i > 0; i-- {
					if s.Text[i] == "" {
						s.Text = s.Text[:i]
					} else {
						break
					}
				}
			}

			// Init subtitle
			s = &Subtitle{}

			// Fetch time boundaries
			boundaries := strings.Split(line, TimeBoundariesSeparator)
			if s.StartAt, err = ParseDurationSRT(boundaries[0]); err != nil {
				return
			}
			if s.EndAt, err = ParseDurationSRT(boundaries[1]); err != nil {
				return
			}

			// Append subtitle
			*o = append(*o, s)
		} else {
			// Add text
			s.Text = append(s.Text, line)
		}
	}
	return
}

// FormatDurationSRT formats an .srt duration
func FormatDurationSRT(i time.Duration) (s string) {
	// Parse hours
	var hours = int(i / time.Hour)
	var n = i % time.Hour
	if hours < 10 {
		s += "0"
	}
	s += strconv.Itoa(hours) + ":"

	// Parse minutes
	var minutes = int(n / time.Minute)
	n = i % time.Minute
	if minutes < 10 {
		s += "0"
	}
	s += strconv.Itoa(minutes) + ":"

	// Parse seconds
	var seconds = int(n / time.Second)
	n = i % time.Second
	if seconds < 10 {
		s += "0"
	}
	s += strconv.Itoa(seconds) + ","

	// Parse milliseconds
	var milliseconds = int(n / time.Millisecond)
	if milliseconds < 10 {
		s += "00"
	} else if milliseconds < 100 {
		s += "0"
	}
	s += strconv.Itoa(milliseconds)
	return
}

// ToWriterSRT formats subtitles as .srt format into a writer
func ToWriterSRT(i Subtitles, o io.Writer) (err error) {
	// Init
	var c []byte

	// Do not write anything if no subtitles
	if len(i) == 0 {
		err = ErrNoSubtitlesToWrite
		return
	}

	// Add BOM header
	c = append(c, BytesBOM...)

	// Loop through subtitles
	for k, v := range i {
		// Init content
		c = append(c, []byte(strconv.Itoa(k+1))...)
		c = append(c, BytesLineSeparator...)
		c = append(c, []byte(FormatDurationSRT(v.StartAt))...)
		c = append(c, BytesTimeBoundariesSeparator...)
		c = append(c, []byte(FormatDurationSRT(v.EndAt))...)
		c = append(c, BytesLineSeparator...)

		// Add text
		for _, t := range v.Text {
			c = append(c, []byte(t)...)
			c = append(c, BytesLineSeparator...)
		}

		// Add new line
		c = append(c, BytesLineSeparator...)
	}

	// Remove last new line
	c = c[:len(c)-1]

	// Write
	if _, err = o.Write(c); err != nil {
		return
	}
	return
}
