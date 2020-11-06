package main

import (
	mod "./model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type BranchesData struct {
	Branches []mod.Branch
}

func (h *Handler) GetBranches(w http.ResponseWriter, r *http.Request) {
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

	rows, err := h.DB.Query(
		"select * from amaker.branches as br where br.idcompany = ?",
		company.Iduser,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	branches := []mod.Branch{}

	for rows.Next() {
		br := mod.Branch{}
		err := rows.Scan(&br.Idbranch, &br.Name, &br.Idcompany)
		if err != nil {
			continue
		}
		branches = append(branches, br)
	}

	data := &BranchesData{
		Branches: branches,
	}

	h.Tmpl.ExecuteTemplate(w, "branches.html", data)
}

func (h *Handler) GetBranch(w http.ResponseWriter, r *http.Request) {
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
		"select * from amaker.branch as b where b.idbranch = ?",
		vars["id"],
	)
	branch := &mod.Branch{}
	err := row.Scan(
		&branch.Idbranch,
		&branch.Name,
		&branch.Idcompany,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := h.DB.Query(
		"select * from amaker.theme as theme where theme.idbranch = ?", vars["id"])

}

func (h *Handler) EditBranch(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) DeleteBranch(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) CreateBranch(w http.ResponseWriter, r *http.Request) {
	cName, err := r.Cookie("user-number")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/home", http.StatusFound)
		return
	}
	cRole, err := r.Cookie("user-role")
	if role, atoiErr := strconv.Atoi(cRole.Value); err == http.ErrNoCookie || role != 1 || atoiErr != nil {
		if atoiErr != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			http.Redirect(w, r, "/home/worker/"+cName.Value, http.StatusFound)
			return
		}
	}

	user, err := h.findUser(cName.Value)

	result, err := h.DB.Exec(
		"INSERT INTO amaker.branch (`name`, `idcompany`) VALUES (?, ?)",
		r.FormValue("name"),
		user.Iduser,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	affected, err := result.RowsAffected()
	lastID, err := result.LastInsertId()
	fmt.Println(affected, "rows affected and last ID is", lastID)
	response := &Response{
		Status:  true,
		Message: "Branch created with success!",
	}
	jsonResponse, marshalError := json.Marshal(response)
	if marshalError != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
