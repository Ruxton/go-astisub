package subtitles

import (
	"errors"
	"os"
	"strings"
	"time"
)

// Vars
var (
	BytesBOM              = []byte{239, 187, 191}
	bytesComma            = []byte(",")
	BytesLineSeparator    = []byte("\n")
	bytesPeriod           = []byte(".")
	ErrInvalidExtension   = errors.New("Invalid extension")
	ErrNoSubtitlesToWrite = errors.New("No subtitles to write")
)

// OSOpen allows testing functions using it
var OSOpen = func(name string) (*os.File, error) {
	return os.Open(name)
}

// Open opens a subtitle file
func Open(name string) (s *Subtitles, err error) {
	// Open the file
	var f *os.File
	if f, err = OSOpen(name); err != nil {
		return
	}
	defer f.Close()

	// Parse the content
	if strings.HasSuffix(name, ".vtt") {
		s, err = FromReaderVTT(f)
	} else if strings.HasSuffix(name, ".srt") {
		s, err = FromReaderSRT(f)
	} else {
		err = ErrInvalidExtension
	}
	return
}

// Subtitles represents an ordered list of subtitles
type Subtitles []*Subtitle

// Duration returns the subtitles duration
func (s Subtitles) Duration() time.Duration {
	if len(s) == 0 {
		return time.Duration(0)
	}
	return s[len(s)-1].EndAt
}

// Subtitle represents a text to show between 2 time boundaries
type Subtitle struct {
	EndAt   time.Duration
	StartAt time.Duration
	Text    []string
}

// Add adds a duration to each time boundaries. As in the time package, duration can be negative.
func (s *Subtitles) Add(d time.Duration) {
	for _, v := range *s {
		v.EndAt += d
		v.StartAt += d
	}
}

// SimulateDuration makes sure the last item is at least ending at the requested duration
func (s *Subtitles) SimulateDuration(d time.Duration) {
	// Subtitles duration is bigger than requested duration
	if s.Duration() >= d {
		return
	}

	// Add dummy item
	*s = append(*s, &Subtitle{EndAt: d, StartAt: d, Text: []string{"Thank you"}})
}

// Fragment fragments subtitles with a specific fragment duration
func (s *Subtitles) Fragment(f time.Duration) {
	// Nothing to fragment
	if len(*s) == 0 {
		return
	}

	// Here we want to simulate fragments of duration f until there are no subtitles left in that period of time
	var fragmentStartAt, fragmentEndAt = time.Duration(0), f
	for fragmentStartAt < (*s)[len(*s)-1].EndAt {
		// We loop through subtitles and process the ones that either contain the fragment start at,
		// or contain the fragment end at
		//
		// It's useless processing subtitles contained between fragment start at and end at
		//             |____________________|             <- subtitle
		//           |                        |
		//   fragment start at        fragment end at
		for i, sub := range *s {
			// A switch is more readable here
			switch {
			// Subtitle contains fragment start at
			// |____________________|                         <- subtitle
			//           |                        |
			//   fragment start at        fragment end at
			case sub.StartAt < fragmentStartAt && sub.EndAt > fragmentStartAt:
				// Init
				var newSub = &Subtitle{}
				*newSub = *sub

				// Update boundaries
				sub.StartAt = fragmentStartAt
				newSub.EndAt = fragmentStartAt

				// Insert new sub
				(*s) = append((*s)[:i], append(Subtitles{newSub}, (*s)[i:]...)...)
			// Subtitle contains fragment end at
			//                         |____________________| <- subtitle
			//           |                        |
			//   fragment start at        fragment end at
			case sub.StartAt < fragmentEndAt && sub.EndAt > fragmentEndAt:
				// Init
				var newSub = &Subtitle{}
				*newSub = *sub

				// Update boundaries
				sub.StartAt = fragmentEndAt
				newSub.EndAt = fragmentEndAt

				// Append new sub
				(*s) = append((*s)[:i], append(Subtitles{newSub}, (*s)[i:]...)...)
			}
		}

		// Update fragments boundaries
		fragmentStartAt += f
		fragmentEndAt += f
	}
}

// OSCreate allows testing functions using it
var OSCreate = func(name string) (*os.File, error) {
	return os.Create(name)
}

// Write writes subtitles to a file
func (s Subtitles) Write(name string) (err error) {
	// Create the file
	var f *os.File
	if f, err = OSCreate(name); err != nil {
		return
	}
	defer f.Close()

	// Write the content
	if strings.HasSuffix(name, ".srt") {
		err = s.ToWriterSRT(f)
	} else {
		err = ErrInvalidExtension
	}
	return
}
