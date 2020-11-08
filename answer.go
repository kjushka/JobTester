package main

import (
	mod "./model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type CustomAnswerUser struct {
	Answer   *mod.Answer
	Username string
	Number   int
}

type AnswersData struct {
	CAU   []*CustomAnswerUser
	UName string
}

func (h *Handler) GetAnswers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GA")
	username, role, status := h.checkCookie(w, r)
	if !status {
		return
	}
	if role != 1 {
		http.Redirect(w, r, "/home/worker/"+username, http.StatusFound)
		return
	}

	user, err := h.findUser(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := h.DB.Query(
		"select a.idanswer, a.file, a.idtask, a.idsender, a.status, a.date, user.username "+
			"from amaker.answer as a "+
			"inner join ("+
			"select t.idtask from amaker.task as t "+
			"right join ( "+
			"select theme.idtheme "+
			"from amaker.theme as theme "+
			"right join ( "+
			"select branch.idbranch "+
			"from amaker.branch as branch "+
			"right join amaker.user as u "+
			"on u.iduser = branch.idcompany "+
			"where u.iduser = ? "+
			") as b "+
			"on b.idbranch = theme.idbranch "+
			") as th "+
			"on t.idtheme = th.idtheme "+
			") as ta "+
			"on ta.idtask = a.idtask "+
			"left join amaker.user as user "+
			"on user.iduser = a.idsender "+
			"where a.status = 0",
		user.Iduser,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	answers := []*CustomAnswerUser{}
	count := 1
	for rows.Next() {
		answer := &mod.Answer{}
		var username string
		err := rows.Scan(
			&answer.Idanswer,
			&answer.File,
			&answer.Idtask,
			&answer.Idsender,
			&answer.Status,
			&answer.Date,
			&username,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cau := &CustomAnswerUser{
			Answer:   answer,
			Username: username,
			Number:   count,
		}
		answers = append(answers, cau)
		count++
	}

	data := &AnswersData{
		CAU:   answers,
		UName: username,
	}

	h.Tmpl.ExecuteTemplate(w, "answer.html", data)
}

func (h *Handler) SetAnswerRight(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SAR")
	username, role, status := h.checkCookie(w, r)
	if !status {
		return
	}
	if role != 1 {
		http.Redirect(w, r, "/home/worker/"+username, http.StatusFound)
		return
	}
	_, err := h.DB.Exec(
		"update amaker.answer as a set status = 1 where a.idanswer = ?",
		mux.Vars(r)["id"],
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/answers/"+mux.Vars(r)["username"], 302)
}

func (h *Handler) SetAnswerWrong(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SAW")
	username, role, status := h.checkCookie(w, r)
	if !status {
		return
	}
	if role != 1 {
		http.Redirect(w, r, "/home/worker/"+username, http.StatusFound)
		return
	}
	_, err := h.DB.Exec(
		"update amaker.answer as a set status = 2 where a.idanswer = ?",
		mux.Vars(r)["id"],
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/answers/"+mux.Vars(r)["username"], 302)
}

func (h *Handler) SetAnswerStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SAS")
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
