package subtitles_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/asticode/go-subtitles"
	"github.com/stretchr/testify/assert"
)

func assertSubtitles(i subtitles.Subtitles, t *testing.T) {
	assert.Len(t, i, 6)
	assert.Equal(t, time.Duration(101370000000), (i)[0].EndAt)
	assert.Equal(t, time.Duration(99000000000), (i)[0].StartAt)
	assert.Equal(t, []string{"(deep rumbling)"}, (i)[0].Text)
	assert.Equal(t, time.Duration(153225000000), (i)[5].EndAt)
	assert.Equal(t, time.Duration(151056000000), (i)[5].StartAt)
	assert.Equal(t, []string{"(computer playing", "electronic melody)"}, (i)[5].Text)
}

func TestParseDurationSRT(t *testing.T) {
	d, err := subtitles.ParseDurationSRT("12:34:56")
	assert.EqualError(t, err, "No milliseconds detected in 12:34:56")
	d, err = subtitles.ParseDurationSRT("12:34:56,1234")
	assert.EqualError(t, err, "Invalid number of millisecond digits detected in 12:34:56,1234")
	d, err = subtitles.ParseDurationSRT("12,123")
	assert.EqualError(t, err, "No hours, minutes or seconds detected in 12,123")
	d, err = subtitles.ParseDurationSRT("12:34,123")
	assert.EqualError(t, err, "No hours, minutes or seconds detected in 12:34,123")
	d, err = subtitles.ParseDurationSRT("12:34:56,123")
	assert.NoError(t, err, "")
	assert.Equal(t, time.Duration(45296123000000), d)
}

func TestFormatDurationSRT(t *testing.T) {
	s := subtitles.FormatDurationSRT(time.Duration(1234567))
	assert.Equal(t, "00:00:00,001", s)
	s = subtitles.FormatDurationSRT(time.Duration(10234567))
	assert.Equal(t, "00:00:00,010", s)
	s = subtitles.FormatDurationSRT(time.Duration(100234567))
	assert.Equal(t, "00:00:00,100", s)
	s = subtitles.FormatDurationSRT(time.Duration(1234567891))
	assert.Equal(t, "00:00:01,234", s)
	s = subtitles.FormatDurationSRT(time.Duration(12345678912))
	assert.Equal(t, "00:00:12,345", s)
	s = subtitles.FormatDurationSRT(time.Duration(123456789123))
	assert.Equal(t, "00:02:03,456", s)
	s = subtitles.FormatDurationSRT(time.Duration(1234567891234))
	assert.Equal(t, "00:20:34,567", s)
	s = subtitles.FormatDurationSRT(time.Duration(12345678912345))
	assert.Equal(t, "03:25:45,678", s)
	s = subtitles.FormatDurationSRT(time.Duration(123456789123456))
	assert.Equal(t, "34:17:36,789", s)
}

func TestSRT(t *testing.T) {
	// Init
	var s *subtitles.Subtitles
	var err error
	var path = "./tests/example.srt"

	// From reader
	t.Run("FromReader", func(t *testing.T) {
		// Open example file
		var file *os.File
		file, err = os.Open(path)
		assert.NoError(t, err)
		defer file.Close()

		// Test
		s, err = subtitles.FromReaderSRT(file)
		assert.NoError(t, err)
		assertSubtitles(*s, t)
	})

	// To writer
	t.Run("ToWriter", func(t *testing.T) {
		// No subtitles
		var w = &bytes.Buffer{}
		err = subtitles.ToWriterSRT(subtitles.Subtitles{}, w)
		assert.EqualError(t, err, subtitles.ErrNoSubtitlesToWrite.Error())

		// Get example file content
		var c []byte
		c, err = ioutil.ReadFile(path)
		assert.NoError(t, err)

		// Test
		err = subtitles.ToWriterSRT(*s, w)
		assert.NoError(t, err)
		assert.Equal(t, string(c), w.String())
	})
}
