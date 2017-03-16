package subtitles_test

import (
	"testing"

	"github.com/asticode/go-subtitles"
	"github.com/stretchr/testify/assert"
)

func TestVTT(t *testing.T) {
	// Init
	var s *subtitles.Subtitles
	var err error
	var path = "./testdata/example.vtt"

	// From reader
	t.Run("FromReaderVTT", func(t *testing.T) {
		s, err = subtitles.Open(path)
		assert.NoError(t, err)
		assertSubtitles(*s, t)
	})
}
