package subtitles

import (
	"errors"
	"os"
	"path/filepath"
	"time"
)

// Vars
var (
	BytesBOM              = []byte{239, 187, 191}
	bytesColon            = []byte(":")
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
	switch filepath.Ext(name) {
	case ".srt":
		s, err = FromReaderSRT(f)
	case ".ttml":
		s, err = FromReaderTTML(f)
	case ".vtt":
		s, err = FromReaderVTT(f)
	default:
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

// Empty returns whether the subtitles are empty
func (s Subtitles) Empty() bool {
	return len(s) == 0
}

// ForceDuration updates the subtitles duration
// If input duration is bigger, then we create a dummy item
// If input duration is smaller, then we remove useless items and we cut the last item
func (s *Subtitles) ForceDuration(d time.Duration) {
	// Input duration is the same as the subtitles'one
	if s.Duration() == d {
		return
	}

	// Input duration is bigger than subtitles'one
	if s.Duration() < d {
		// Add dummy item
		*s = append(*s, &Subtitle{EndAt: d, StartAt: d, Text: []string{"..."}})
	} else {
		// Find last item before input duration and update end at
		var lastIndex = -1
		for index, i := range *s {
			// Start at is bigger than input duration, we've found the last item
			if i.StartAt >= d {
				lastIndex = index
				break
			} else if i.EndAt > d {
				(*s)[index].EndAt = d
			}
		}

		// Last index has been found
		if lastIndex != -1 {
			(*s) = (*s)[:lastIndex]
		}
	}
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

// Merge merges subtitles i into subtitles s
func (s *Subtitles) Merge(i Subtitles) {
	// Loop through input subtitles
	for _, subInput := range i {
		var lastIndex int
		var inserted bool
		// Loop through parent subtitles
		for index, subParent := range *s {
			// Input sub is after parent sub
			if subInput.StartAt < subParent.StartAt {
				*s = append((*s)[:lastIndex+1], append([]*Subtitle{subInput}, (*s)[lastIndex+1:]...)...)
				inserted = true
				break
			}
			lastIndex = index
		}
		if !inserted {
			*s = append(*s, subInput)
		}
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
	switch filepath.Ext(name) {
	case ".srt":
		err = s.ToWriterSRT(f)
	case ".tml":
		err = s.ToWriterTTML(f)
	default:
		err = ErrInvalidExtension
	}
	return
}
