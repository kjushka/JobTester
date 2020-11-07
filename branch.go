package main

import (
	mod "./model"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
	"time"
)

type BranchesData struct {
	Branches  []*mod.Branch
	UName     string
	IsCompany int
}

type BranchData struct {
	Branch    *mod.Branch
	UName     string
	IsCompany int
}

type customThemeAndTask struct {
	Idtheme      int
	ThemeName    string
	ThemeIndex   int
	Idtask       sql.NullInt32
	TaskName     sql.NullString
	TaskText     sql.NullString
	Idanswer     sql.NullInt32
	File         sql.NullString
	Idsender     sql.NullInt32
	AnswerStatus sql.NullInt32
	AnswerDate   sql.NullTime
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
		"select * from amaker.branch as br where br.idcompany = ?",
		company.Iduser,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	branches := []*mod.Branch{}

	for rows.Next() {
		br := &mod.Branch{}
		err := rows.Scan(&br.Idbranch, &br.Name, &br.Idcompany)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		br.Username = username
		branches = append(branches, br)
	}

	data := &BranchesData{
		Branches:  branches,
		UName:     username,
		IsCompany: role,
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

	user, err := h.findUser(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	row := h.DB.QueryRow(
		"select * from amaker.branch as b where b.idbranch = ?",
		vars["id"],
	)
	branch := &mod.Branch{}
	err = row.Scan(
		&branch.Idbranch,
		&branch.Name,
		&branch.Idcompany,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := h.DB.Query(
		"select th.idtheme, th.name, th.index, t.idtask, t.name, t.text, t.idanswer, t.file, t.idsender, t.status, t.date "+
			"from amaker.theme as th left join ( "+
			"select task.idtask, task.name, task.text, task.idtheme, a.idanswer, a.file, a.idsender, a.status, a.date from amaker.task as task "+
			"left join ( "+
			"select * from amaker.answer as ans "+
			"where ans.idsender = ? "+
			") as a "+
			"on task.idtask = a.idtask "+
			") as t "+
			"on th.idtheme = t.idtheme "+
			"where th.idbranch = ? ",
		user.Iduser,
		vars["id"],
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	themeMap := make(map[int]*mod.Theme)

	for rows.Next() {
		ctat := &customThemeAndTask{}
		err := rows.Scan(
			&ctat.Idtheme,
			&ctat.ThemeName,
			&ctat.ThemeIndex,
			&ctat.Idtask,
			&ctat.TaskName,
			&ctat.TaskText,
			&ctat.Idanswer,
			&ctat.File,
			&ctat.Idsender,
			&ctat.AnswerStatus,
			&ctat.AnswerDate,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var task *mod.Task = nil
		if ctat.Idtask.Valid {
			var answer *mod.Answer = nil
			if ctat.Idanswer.Valid {
				answer = &mod.Answer{
					Idanswer: int(ctat.Idanswer.Int32),
					File:     ctat.File.String,
					Idsender: int(ctat.Idsender.Int32),
					Idtask:   int(ctat.Idtask.Int32),
					Status:   int(ctat.AnswerStatus.Int32),
					Date:     ctat.AnswerDate.Time,
				}
			}
			task = &mod.Task{
				Idtask:  int(ctat.Idtask.Int32),
				Name:    ctat.TaskName.String,
				Text:    ctat.TaskText.String,
				Idtheme: ctat.Idtheme,
				Answer:  answer,
			}
		}
		if val, ok := themeMap[ctat.ThemeIndex]; ok {
			if task != nil {
				val.Tasks = append(val.Tasks, task)
			}
		} else {
			theme := &mod.Theme{
				Idtheme:  ctat.Idtheme,
				Name:     ctat.ThemeName,
				Idbranch: branch.Idbranch,
				Index:    ctat.ThemeIndex,
				Tasks:    []*mod.Task{},
			}
			if task != nil {
				theme.Tasks = append(theme.Tasks, task)
			}
			themeMap[ctat.ThemeIndex] = theme
		}
	}

	keys := make([]int, 0, len(themeMap))
	for k, _ := range themeMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	themeSlice := make([]*mod.Theme, 0, len(themeMap))
	for _, k := range keys {
		themeSlice = append(themeSlice, themeMap[k])
	}

	branch.Themes = themeSlice

	branchData := &BranchData{
		Branch:    branch,
		UName:     username,
		IsCompany: role,
	}

	h.Tmpl.ExecuteTemplate(w, "tasks.html", branchData)
}

func (h *Handler) EditBranch(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) DeleteBranch(w http.ResponseWriter, r *http.Request) {
	username, role, status := h.checkCookie(w, r)
	if !status {
		return
	}
	if role != 1 {
		http.Redirect(w, r, "/home/worker/"+username, http.StatusFound)
		return
	}
	vars := mux.Vars(r)
	result, err := h.DB.Exec("delete from amaker.branch where idbranch = ?", vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(result.LastInsertId())
	fmt.Println(result.RowsAffected())
	http.Redirect(w, r, "/branches/"+username, 302)
}

func (h *Handler) CreateBranch(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.DB.Exec(
		"INSERT INTO amaker.branch (`name`, `idcompany`) VALUES (?, ?)",
		r.FormValue("name"),
		user.Iduser,
	)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	affected, err := result.RowsAffected()
	lastID, err := result.LastInsertId()
	fmt.Println(affected, "rows affected and last ID is", lastID)
	http.Redirect(w, r, "/branches/"+username, 302)
}

func (h *Handler) SendAnswer(w http.ResponseWriter, r *http.Request) {

	fmt.Println("sans")
	username, role, status := h.checkCookie(w, r)
	if !status {
		return
	}
	if role != 0 {
		http.Redirect(w, r, "/home/company/"+username, http.StatusFound)
		return
	}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	answer := &mod.Answer{}
	err := decoder.Decode(answer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := h.findUser(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	answer.Idsender = user.Iduser
	answer.Date = time.Now()
	result, err := h.DB.Exec(
		"insert into amaker.answer (`file`, `idsender`, `idtask`, `status`, `date`) values (?, ?, ?, ?, ?)",
		answer.File,
		answer.Idsender,
		answer.Idtask,
		answer.Status,
		answer.Date,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	affected, err := result.RowsAffected()
	lastID, err := result.LastInsertId()
	fmt.Println(affected, "rows affected and last ID is", lastID)

	response := &Response{
		Status:  true,
		Message: "Answer sent successfully!",
	}
	jsonResponse, marshalError := json.Marshal(response)
	if marshalError != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
