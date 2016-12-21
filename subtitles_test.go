package subtitles_test

import (
	"testing"
	"time"

	"github.com/asticode/go-subtitles"
	"github.com/stretchr/testify/assert"
)

func TestSubtitles_Duration(t *testing.T) {
	assert.Equal(t, time.Duration(0), subtitles.Subtitles{}.Duration())
	assert.Equal(t, 7*time.Second, subtitles.Subtitles{&subtitles.Subtitle{EndAt: 3 * time.Second, StartAt: time.Second}, &subtitles.Subtitle{EndAt: 7 * time.Second, StartAt: 3 * time.Second}}.Duration())
}

func TestSubtitles_Add(t *testing.T) {
	var s = subtitles.Subtitles{&subtitles.Subtitle{EndAt: 3 * time.Second, StartAt: time.Second}, &subtitles.Subtitle{EndAt: 7 * time.Second, StartAt: 3 * time.Second}}
	s.Add(time.Second)
	assert.Len(t, s, 2)
	assert.Equal(t, &subtitles.Subtitle{EndAt: 4 * time.Second, StartAt: 2 * time.Second}, s[0])
	assert.Equal(t, &subtitles.Subtitle{EndAt: 8 * time.Second, StartAt: 4 * time.Second}, s[1])
}

func TestSubtitles_SimulateDuration(t *testing.T) {
	var s = subtitles.Subtitles{&subtitles.Subtitle{EndAt: 3 * time.Second, StartAt: time.Second}, &subtitles.Subtitle{EndAt: 7 * time.Second, StartAt: 3 * time.Second}}
	s.SimulateDuration(10 * time.Second)
	assert.Len(t, s, 3)
	assert.Equal(t, 10*time.Second, s[2].EndAt)
	assert.Equal(t, 10*time.Second, s[2].StartAt)
	assert.Equal(t, "Thank you", s[2].Text[0])
}

func TestSubtitles_Fragment(t *testing.T) {
	var s = subtitles.Subtitles{&subtitles.Subtitle{EndAt: 3 * time.Second, StartAt: time.Second, Text: []string{"subtitle-1"}}, &subtitles.Subtitle{EndAt: 7 * time.Second, StartAt: 3 * time.Second, Text: []string{"subtitle-2"}}}
	s.Fragment(2 * time.Second)
	assert.Len(t, s, 5)
	assert.Equal(t, subtitles.Subtitle{EndAt: 2 * time.Second, StartAt: time.Second, Text: []string{"subtitle-1"}}, *s[0])
	assert.Equal(t, subtitles.Subtitle{EndAt: 3 * time.Second, StartAt: 2 * time.Second, Text: []string{"subtitle-1"}}, *s[1])
	assert.Equal(t, subtitles.Subtitle{EndAt: 4 * time.Second, StartAt: 3 * time.Second, Text: []string{"subtitle-2"}}, *s[2])
	assert.Equal(t, subtitles.Subtitle{EndAt: 6 * time.Second, StartAt: 4 * time.Second, Text: []string{"subtitle-2"}}, *s[3])
	assert.Equal(t, subtitles.Subtitle{EndAt: 7 * time.Second, StartAt: 6 * time.Second, Text: []string{"subtitle-2"}}, *s[4])
}
