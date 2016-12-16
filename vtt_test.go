package subtitles_test

import (
	"os"
	"testing"

	"github.com/asticode/go-subtitles"
	"github.com/stretchr/testify/assert"
)

func TestVTT(t *testing.T) {
	// Init
	var s *subtitles.Subtitles
	var err error
	var pathVTT = "./tests/example.vtt"

	// From reader
	t.Run("FromReaderVTT", func(t *testing.T) {
		// Open example file
		var fileVTT *os.File
		fileVTT, err = os.Open(pathVTT)
		assert.NoError(t, err)
		defer fileVTT.Close()

		// Test
		s, err = subtitles.FromReaderVTT(fileVTT)
		assert.NoError(t, err)
		assertSubtitles(*s, t)
	})
}
