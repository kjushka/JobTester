package main

import (
	mod "./model"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "home.html", struct{}{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetWorker(w http.ResponseWriter, r *http.Request) {
	username, role, status := h.checkCookie(w, r)
	if !status {
		return
	}
	if role != 0 {
		http.Redirect(w, r, "/home/company/"+username, http.StatusFound)
		return
	}
	user, err := h.findUser(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.Tmpl.ExecuteTemplate(w, "worker.html", user)
}

func (h *Handler) EditWorker(w http.ResponseWriter, r *http.Request) {
	username, role, status := h.checkCookie(w, r)
	if !status {
		return
	}
	if role != 0 {
		http.Redirect(w, r, "/home/company/"+username, http.StatusFound)
		return
	}
	result, err := h.updateUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(result.LastInsertId())
	fmt.Println(result.RowsAffected())
	response := &Response{
		Status:  true,
		Message: "Update successfully!",
	}
	jsonResponse, marshalError := json.Marshal(response)
	if marshalError != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (h *Handler) GetCompany(w http.ResponseWriter, r *http.Request) {
	username, role, status := h.checkCookie(w, r)
	if !status {
		return
	}
	if role != 1 {
		http.Redirect(w, r, "/home/worker/"+username, http.StatusFound)
		return
	}
	company, err := h.findUser(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.Tmpl.ExecuteTemplate(w, "company.html", company)
}

func (h *Handler) EditCompany(w http.ResponseWriter, r *http.Request) {
	username, role, status := h.checkCookie(w, r)
	if !status {
		return
	}
	if role != 1 {
		http.Redirect(w, r, "/home/worker/"+username, http.StatusFound)
		return
	}
	result, err := h.updateUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(result.LastInsertId())
	fmt.Println(result.RowsAffected())
	response := &Response{
		Status:  true,
		Message: "Update successfully!",
	}
	jsonResponse, marshalError := json.Marshal(response)
	if marshalError != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) (sql.Result, error) {
	decoder := json.NewDecoder(r.Body)
	worker := &mod.User{}
	err := decoder.Decode(&worker)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	result, err := h.DB.Exec(
		"update amaker.user as u where username = ?, password = ?, name = ? where u.iduser = ?",
		worker.Username,
		worker.Password,
		worker.Name,
		worker.Iduser,
	)
	return result, err
}
