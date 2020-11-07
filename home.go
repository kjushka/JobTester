package main

import (
	"database/sql"
	"net/http"
	"strconv"
)

type LogErr struct {
	Error string
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	cName, err := r.Cookie("user-name")
	if err != http.ErrNoCookie {
		cRole, err := r.Cookie("user-role")
		if err != http.ErrNoCookie {
			role, atoiErr := strconv.Atoi(cRole.Value)
			if atoiErr != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if role == 0 {
				http.Redirect(w, r, "/home/worker/"+cName.Value, http.StatusFound)
				return
			} else {
				http.Redirect(w, r, "/home/company/"+cName.Value, http.StatusFound)
				return
			}
		}
	}

	err = h.Tmpl.ExecuteTemplate(w, "home.html", LogErr{
		Error: "",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HomeWorker(w http.ResponseWriter, r *http.Request) {
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
	_, err := h.updateUser(w, r, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/home/company/"+username, 302)
}

func (h *Handler) HomeCompany(w http.ResponseWriter, r *http.Request) {
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
	h.Tmpl.ExecuteTemplate(w, "worker.html", company)
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
	_, err := h.updateUser(w, r, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/home/company/"+username, 302)
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request, username string) (sql.Result, error) {
	name := r.FormValue("name")
	if name == "" {
		return nil, nil
	}
	worker, err := h.findUser(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	result, err := h.DB.Exec(
		"update amaker.user as u set name = ? where u.iduser = ?",
		name,
		worker.Iduser,
	)
	return result, err
}
