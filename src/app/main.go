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
		shadow    = 2
		padding   = 10
		charWidth = 8
		radius    = 2
	)

	kindWidth := charWidth * len(kind)
	labelWidth := charWidth * len(label)
	width := 4*padding + kindWidth + labelWidth

	s := svg.New(w)
	s.Start(width+shadow, height+shadow)
	s.Style("text/css", `
		text {
			alignment-baseline: middle;
			font-family: Courier, monospace;
			font-size: 12px;
			fill: white;
			text-shadow: 1px 1px #666
		}`)

	s.Roundrect(shadow, shadow, width, height, radius, radius, "fill:#aaa")
	s.Roundrect(0, 0, width, height, radius, radius, "fill:#333")
	s.Roundrect(2*padding+kindWidth, 0, 2*padding+labelWidth, height, radius, radius, "fill:#"+color)
	s.Text(padding, height/2, kind)
	s.Text(3*padding+kindWidth, height/2, label)
	s.End()
}
