package main

import (
	"io/ioutil"
	"net/http"
)

func (h *Handler) SendBranchesJs(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("./webapp/resources/branches.js")
	if err != nil {
		http.Error(w, "Couldn't read file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Write(data)
}

func (h *Handler) SendBranchesCss(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("./webapp/resources/branches.css")
	if err != nil {
		http.Error(w, "Couldn't read file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Write(data)
}

func (h *Handler) SendHomeJs(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("./webapp/resources/home.js")
	if err != nil {
		http.Error(w, "Couldn't read file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Write(data)
}

func (h *Handler) SendCssCss(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("./webapp/resources/css.css")
	if err != nil {
		http.Error(w, "Couldn't read file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Write(data)
}
