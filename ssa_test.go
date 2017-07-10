package astisub_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/asticode/go-astisub"
	"github.com/stretchr/testify/assert"
)

func TestSSA(t *testing.T) {
	// Open
	s, err := astisub.OpenFile("./testdata/example-in.ssa")
	assert.NoError(t, err)
	assertSubtitleItems(t, s)
	// Metadata
	assert.Equal(t, &astisub.Metadata{Comments: []string{"Comment 1", "Comment 2"}, Copyright: "Copyright test", Title: "SSA test"}, s.Metadata)
	// Styles
	assert.Equal(t, 3, len(s.Styles))
	assert.Equal(t, astisub.Style{ID: "1", InlineStyle: &astisub.StyleAttributes{SSAAlignment: 7, SSAAlphaLevel: 0.1, SSABackColour: &astisub.Color{Alpha: 128, Red: 8}, SSABold: true, SSABorderStyle: 7, SSAFontName: "f1", SSAFontSize: 4, SSAOutline: 1, SSAOutlineColour: &astisub.Color{Green: 255, Red: 255}, SSAMarginLeft: 1, SSAMarginRight: 4, SSAMarginVertical: 7, SSAPrimaryColour: &astisub.Color{Green: 255, Red: 255}, SSASecondaryColour: &astisub.Color{Green: 255, Red: 255}, SSAShadow: 4}}, *s.Styles["1"])
	assert.Equal(t, astisub.Style{ID: "2", InlineStyle: &astisub.StyleAttributes{SSAAlignment: 8, SSAAlphaLevel: 0.2, SSABackColour: &astisub.Color{Blue: 15, Green: 15, Red: 15}, SSABold: true, SSABorderStyle: 8, SSAEncoding: 1, SSAFontName: "f2", SSAFontSize: 5, SSAOutline: 2, SSAOutlineColour: &astisub.Color{Green: 255, Red: 255}, SSAMarginLeft: 2, SSAMarginRight: 5, SSAMarginVertical: 8, SSAPrimaryColour: &astisub.Color{Blue: 239, Green: 239, Red: 239}, SSASecondaryColour: &astisub.Color{Green: 255, Red: 255}, SSAShadow: 5}}, *s.Styles["2"])
	assert.Equal(t, astisub.Style{ID: "3", InlineStyle: &astisub.StyleAttributes{SSAAlignment: 9, SSAAlphaLevel: 0.3, SSABackColour: &astisub.Color{Red: 8}, SSABorderStyle: 9, SSAEncoding: 2, SSAFontName: "f3", SSAFontSize: 6, SSAOutline: 3, SSAOutlineColour: &astisub.Color{Red: 8}, SSAMarginLeft: 3, SSAMarginRight: 6, SSAMarginVertical: 9, SSAPrimaryColour: &astisub.Color{Blue: 180, Green: 252, Red: 252}, SSASecondaryColour: &astisub.Color{Blue: 180, Green: 252, Red: 252}, SSAShadow: 6}}, *s.Styles["3"])
	return

	// No subtitles to write
	w := &bytes.Buffer{}
	err = astisub.Subtitles{}.WriteToSTL(w)
	assert.EqualError(t, err, astisub.ErrNoSubtitlesToWrite.Error())

	// Write
	c, err := ioutil.ReadFile("./testdata/example-out.stl")
	assert.NoError(t, err)
	err = s.WriteToSTL(w)
	assert.NoError(t, err)
	assert.Equal(t, string(c), w.String())
}
