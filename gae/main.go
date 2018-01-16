package app

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/campoy/svg-badge/badge"
	"github.com/gorilla/mux"
)

var maker *badge.Maker

func init() {
	data, err := ioutil.ReadFile("DejaVuSans.ttf")
	if err != nil {
		log.Fatalf("could not load font file: %v", err)
	}

	m, err := badge.NewMaker("Dejavu Sans", data)
	if err != nil {
		log.Fatalf("could not created badge maker: %v", err)
	}
	maker = m

	r := mux.NewRouter()
	r.HandleFunc("/badge/{kind}/{label}", handler)
	http.Handle("/", r)
}

var colors = map[string]string{
	"deprecated":   "C62914",
	"experimental": "DD5F0A",
	"frozen":       "4b4b4b",
	"locked":       "14C6C6",
	"stable":       "74C614",
}

func handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	kind := vars["kind"]
	label := vars["label"]
	color := colors[label]
	if v := r.FormValue("color"); v != "" {
		color = v
	}

	b, err := maker.New(kind, label, color, 11, 20)
	if err != nil {
		http.Error(w, "could not create badge", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "max-age=31556926")
	w.Write(b)
}
