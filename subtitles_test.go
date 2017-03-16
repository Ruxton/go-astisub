package subtitles_test

import (
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

func TestSubtitles_Add(t *testing.T) {
	var s = subtitles.Subtitles{&subtitles.Subtitle{EndAt: 3 * time.Second, StartAt: time.Second}, &subtitles.Subtitle{EndAt: 7 * time.Second, StartAt: 3 * time.Second}}
	s.Add(time.Second)
	assert.Len(t, s, 2)
	assert.Equal(t, &subtitles.Subtitle{EndAt: 4 * time.Second, StartAt: 2 * time.Second}, s[0])
	assert.Equal(t, &subtitles.Subtitle{EndAt: 8 * time.Second, StartAt: 4 * time.Second}, s[1])
}

func TestSubtitles_Duration(t *testing.T) {
	assert.Equal(t, time.Duration(0), subtitles.Subtitles{}.Duration())
	assert.Equal(t, 7*time.Second, subtitles.Subtitles{&subtitles.Subtitle{EndAt: 3 * time.Second, StartAt: time.Second}, &subtitles.Subtitle{EndAt: 7 * time.Second, StartAt: 3 * time.Second}}.Duration())
}

func TestSubtitles_Empty(t *testing.T) {
	assert.True(t, subtitles.Subtitles{}.Empty())
	assert.False(t, subtitles.Subtitles{&subtitles.Subtitle{EndAt: 3 * time.Second, StartAt: time.Second}, &subtitles.Subtitle{EndAt: 7 * time.Second, StartAt: 3 * time.Second}}.Empty())
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

func TestSubtitles_Merge(t *testing.T) {
	var s1 = subtitles.Subtitles{&subtitles.Subtitle{EndAt: 3 * time.Second, StartAt: time.Second}, &subtitles.Subtitle{EndAt: 8 * time.Second, StartAt: 5 * time.Second}, &subtitles.Subtitle{EndAt: 12 * time.Second, StartAt: 10 * time.Second}}
	var s2 = subtitles.Subtitles{&subtitles.Subtitle{EndAt: 4 * time.Second, StartAt: 2 * time.Second}, &subtitles.Subtitle{EndAt: 7 * time.Second, StartAt: 6 * time.Second}, &subtitles.Subtitle{EndAt: 11 * time.Second, StartAt: 9 * time.Second}, &subtitles.Subtitle{EndAt: 14 * time.Second, StartAt: 13 * time.Second}}
	s1.Merge(s2)
	assert.Len(t, s1, 7)
	assert.Equal(t, &subtitles.Subtitle{EndAt: 3 * time.Second, StartAt: time.Second}, s1[0])
	assert.Equal(t, &subtitles.Subtitle{EndAt: 4 * time.Second, StartAt: 2 * time.Second}, s1[1])
	assert.Equal(t, &subtitles.Subtitle{EndAt: 8 * time.Second, StartAt: 5 * time.Second}, s1[2])
	assert.Equal(t, &subtitles.Subtitle{EndAt: 7 * time.Second, StartAt: 6 * time.Second}, s1[3])
	assert.Equal(t, &subtitles.Subtitle{EndAt: 11 * time.Second, StartAt: 9 * time.Second}, s1[4])
	assert.Equal(t, &subtitles.Subtitle{EndAt: 12 * time.Second, StartAt: 10 * time.Second}, s1[5])
	assert.Equal(t, &subtitles.Subtitle{EndAt: 14 * time.Second, StartAt: 13 * time.Second}, s1[6])
}

func TestSubtitles_ForceDuration(t *testing.T) {
	var s = subtitles.Subtitles{&subtitles.Subtitle{EndAt: 3 * time.Second, StartAt: time.Second}, &subtitles.Subtitle{EndAt: 7 * time.Second, StartAt: 3 * time.Second}}
	s.ForceDuration(10 * time.Second)
	assert.Len(t, s, 3)
	assert.Equal(t, 10*time.Second, s[2].EndAt)
	assert.Equal(t, 10*time.Second, s[2].StartAt)
	assert.Equal(t, "...", s[2].Text[0])
	s[2].StartAt = 7 * time.Second
	s[2].EndAt = 12 * time.Second
	s.ForceDuration(10 * time.Second)
	assert.Len(t, s, 3)
	assert.Equal(t, 10*time.Second, s[2].EndAt)
	assert.Equal(t, 7*time.Second, s[2].StartAt)
}
