package main

import (
	mod "./model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type CustomAnswerUser struct {
	Answer *mod.Answer
	User   *mod.User
}

func (h *Handler) GetAnswers(w http.ResponseWriter, r *http.Request) {
	username, role, status := h.checkCookie(w, r)
	if !status {
		return
	}
	if role != 1 {
		http.Redirect(w, r, "/home/worker/"+username, http.StatusFound)
		return
	}
}

func (h *Handler) GetAnswer(w http.ResponseWriter, r *http.Request) {
	username, role, status := h.checkCookie(w, r)
	if !status {
		return
	}
	if role != 1 {
		http.Redirect(w, r, "/home/worker/"+username, http.StatusFound)
		return
	}
	vars := mux.Vars(r)
	row := h.DB.QueryRow(
		"select * from amaker.answer as a join amaker.user as u on a.idsender = u.iduser where a.idanswer = ?",
		vars["id"],
	)
	answer := &mod.Answer{}
	sender := &mod.User{}
	err := row.Scan(
		&answer.Idanswer,
		&answer.File,
		&answer.Idsender,
		&answer.Idtask,
		&answer.Status,
		&answer.Date,
		&sender.Iduser,
		&sender.Username,
		&sender.Password,
		&sender.Name,
		&sender.IsCompany,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cau := &CustomAnswerUser{
		Answer: answer,
		User:   sender,
	}
	h.Tmpl.ExecuteTemplate(w, "answer.html", cau)
}

func (h *Handler) SetAnswerStatus(w http.ResponseWriter, r *http.Request) {
	username, role, status := h.checkCookie(w, r)
	if !status {
		return
	}
	if role != 1 {
		http.Redirect(w, r, "/home/worker/"+username, http.StatusFound)
		return
	}
	vars := mux.Vars(r)
	result, err := h.DB.Exec(
		"update amaker.answer as a set status = ? where a.idanswer = ?",
		vars["id"],
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(result.LastInsertId())
	fmt.Println(result.RowsAffected())

	response := &Response{
		Status:  true,
		Message: "Status update successfully!",
	}
	jsonResponse, marshalError := json.Marshal(response)
	if marshalError != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (h *Handler) DownloadAnswer(w http.ResponseWriter, r *http.Request) {
}
