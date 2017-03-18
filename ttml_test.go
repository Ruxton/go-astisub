package astisub_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	astisub "github.com/asticode/go-astisub"
	"github.com/stretchr/testify/assert"
)

func TestTTML(t *testing.T) {
	// Init
	var s *astisub.Subtitles
	var err error

	// From reader
	t.Run("FromReader", func(t *testing.T) {
		// Open example file
		var file *os.File
		file, err = os.Open("./testdata/example-in.ttml")
		assert.NoError(t, err)
		defer file.Close()

		// Test
		s, err = astisub.FromReaderTTML(file)
		assert.NoError(t, err)
		assertSubtitles(*s, t)
	})

	// To writer
	t.Run("ToWriter", func(t *testing.T) {
		// No subtitles
		var w = &bytes.Buffer{}
		err = astisub.Subtitles{}.ToWriterTTML(w)
		assert.EqualError(t, err, astisub.ErrNoSubtitlesToWrite.Error())

		// Get example file content
		var c []byte
		c, err = ioutil.ReadFile("./testdata/example-out.ttml")
		assert.NoError(t, err)

		// Test
		err = (*s).ToWriterTTML(w)
		assert.NoError(t, err)
		assert.Equal(t, string(c), w.String())
	})
}
