package app

import (
	"net/http"

	svg "github.com/ajstarks/svgo"
	"github.com/gorilla/mux"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/badge/{kind}/{label}", badgeHandler)
	http.Handle("/", r)
}

var colors = map[string]string{
	"deprecated":   "C62914",
	"experimental": "DD5F0A",
	"frozen":       "4b4b4b",
	"locked":       "14C6C6",
	"stable":       "74C614",
}

func badgeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "max-age=31556926")

	vars := mux.Vars(r)
	kind := vars["kind"]
	label := vars["label"]
	color := colors[label]
	if v := r.FormValue("color"); v != "" {
		color = v
	}

	const (
		height    = 20
		padding   = 7
		charWidth = 7
		radius    = 3
	)

	kindWidth := charWidth * len(kind)
	labelWidth := charWidth * len(label)
	width := 4*padding + kindWidth + labelWidth

	s := svg.New(w)
	s.Start(width, height)
	s.Style("text/css", `
		text {
			font-family: DejaVu Sans,Verdana,Geneva,sans-serif;
			font-size: 11px;
			fill: white;
		}`)
	s.LinearGradient("gradient", 0, 0, 0, 255, []svg.Offcolor{
		{Offset: 0, Color: "#000", Opacity: 0},
		{Offset: 200, Color: "#000", Opacity: 0.25},
		{Offset: 225, Color: "#000", Opacity: 1.0},
	})
	s.Roundrect(0, 0, width, height, radius, radius, "fill:#666")
	s.Roundrect(padding+kindWidth, 0, 3*padding+labelWidth, height, radius, radius, "fill:#"+color)
	s.Rect(padding+kindWidth, 0, 10, height, "fill:#666")
	s.Roundrect(0, 0, width, height, radius, radius, "fill:url(#gradient)")
	s.Text(padding, 14, kind)
	s.Text(3*padding+kindWidth, 14, label)
	s.End()
}
