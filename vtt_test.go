package astisub_test

import (
	"testing"

	astisub "github.com/asticode/go-astisub"
	"github.com/stretchr/testify/assert"
)

func TestVTT(t *testing.T) {
	// Init
	var s *astisub.Subtitles
	var err error
	var path = "./testdata/example.vtt"

	// From reader
	t.Run("FromReaderVTT", func(t *testing.T) {
		s, err = astisub.Open(path)
		assert.NoError(t, err)
		assertSubtitles(*s, t)
	})
}
