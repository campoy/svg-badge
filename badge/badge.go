// Package badge provides an easy way to create SVG badges.
package badge

import (
	"bytes"
	"fmt"

	svg "github.com/ajstarks/svgo"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
)

// A Maker knows how to make badges.
type Maker struct {
	fontName string
	font     *truetype.Font
}

// NewMaker returns a new Maker given the font it should use.
func NewMaker(fontName string, fontData []byte) (*Maker, error) {
	font, err := freetype.ParseFont(fontData)
	if err != nil {
		return nil, err
	}
	return &Maker{font: font, fontName: fontName}, nil
}

func (m *Maker) width(s string, size float64) (int, error) {
	ctx := freetype.NewContext()
	ctx.SetFont(m.font)
	ctx.SetFontSize(size)

	pair, err := ctx.DrawString(s, fixed.Point26_6{})
	if err != nil {
		return 0, err
	}
	return pair.X.Ceil(), nil
}

// New creates a new badge given the wished text, color, font size, and height.
func (m *Maker) New(left, right, color string, fontSize float64, height int) ([]byte, error) {
	leftWidth, err := m.width(left, fontSize)
	if err != nil {
		return nil, err
	}
	rightWidth, err := m.width(right, fontSize)
	if err != nil {
		return nil, err
	}
	padding, err := m.width("A", fontSize)
	if err != nil {
		return nil, err
	}

	bounds := m.font.Bounds(fixed.I(int(fontSize)))
	fontHeight := (bounds.Max.Y - bounds.Min.Y).Ceil()
	verticalCenter := (height + fontHeight/3) / 2

	radius := 5
	width := 4*padding + rightWidth + leftWidth

	buf := new(bytes.Buffer)
	s := svg.New(buf)
	s.Start(width, height)

	var textStyle = fmt.Sprintf(
		`text { font-family: %s,Verdana,Geneva,sans-serif; font-size: %dpx; fill: white; }`,
		m.fontName, int(fontSize))

	s.Style("text/css", textStyle)
	s.LinearGradient("gradient", 0, 0, 0, 255, []svg.Offcolor{
		{Offset: 0, Color: "#000", Opacity: 0},
		{Offset: 200, Color: "#000", Opacity: 0.25},
		{Offset: 225, Color: "#000", Opacity: 1.0},
	})

	s.Roundrect(0, 0, width, height, radius, radius, "fill:#666")
	s.Roundrect(2*padding+leftWidth, 0, 2*padding+rightWidth, height, radius, radius, "fill:#"+color)
	s.Rect(2*padding+leftWidth, 0, radius, height, "fill:#"+color)
	s.Roundrect(0, 0, width, height, radius, radius, "fill:url(#gradient)")
	s.Text(padding, verticalCenter, left)
	s.Text(3*padding+leftWidth, verticalCenter, right)
	s.End()

	return buf.Bytes(), nil
}
