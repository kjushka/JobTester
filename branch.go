package main

import (
	mod "./model"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
	"strconv"
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
	fmt.Println("GBS")
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
	fmt.Println("GB")
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
				Idtask:   int(ctat.Idtask.Int32),
				Name:     ctat.TaskName.String,
				Text:     ctat.TaskText.String,
				Idtheme:  ctat.Idtheme,
				Answer:   answer,
				Username: username,
				Idbranch: branch.Idbranch,
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
				Username: username,
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

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DT")
	vars := mux.Vars(r)
	_, err := h.DB.Exec(
		"delete from amaker.task as t where t.idtask = ?",
		vars["idTask"],
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(
		w,
		r,
		"/branches/"+vars["username"]+"/"+vars["idBranch"],
		http.StatusFound,
	)
}

func (h *Handler) DeleteTheme(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DTH")
	vars := mux.Vars(r)
	_, err := h.DB.Exec(
		"delete from amaker.theme as t where t.idtheme = ?",
		vars["idTheme"],
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(
		w,
		r,
		"/branches/"+vars["username"]+"/"+vars["idBranch"],
		http.StatusFound,
	)
}

func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AT")
	name := r.FormValue("name")
	text := r.FormValue("text")
	idTheme := r.FormValue("idtheme")
	fmt.Println(idTheme)
	_, err := h.DB.Exec(
		"insert into amaker.task (`name`, `text`, `idtheme`) values (?, ?, ?)",
		name,
		text,
		idTheme,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	http.Redirect(
		w,
		r,
		"/branches/"+vars["username"]+"/"+vars["idBranch"],
		http.StatusFound,
	)
}

func (h *Handler) AddTheme(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ATH")
	name := r.FormValue("name")
	vars := mux.Vars(r)

	row := h.DB.QueryRow(
		"select count(*) from amaker.theme as th where th.idbranch = ?",
		vars["idBranch"],
	)
	var count int
	err := row.Scan(&count)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.DB.Exec(
		"insert into amaker.theme (`name`, `idbranch`, `index`) values (?, ?, ?)",
		name,
		vars["idBranch"],
		count+1,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(
		w,
		r,
		"/branches/"+vars["username"]+"/"+vars["idBranch"],
		http.StatusFound,
	)
}

func (h *Handler) EditTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ET")
	newName := r.FormValue("newName")
	newText := r.FormValue("newText")
	idtask := r.FormValue("idtask")
	vars := mux.Vars(r)
	_, err := h.DB.Exec(
		"update amaker.task set name = ?, text = ? where idtask = ?",
		newName,
		newText,
		idtask,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(
		w,
		r,
		"/branches/"+vars["username"]+"/"+vars["idBranch"],
		http.StatusFound,
	)
}

func (h *Handler) EditTheme(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ETH")
	newName := r.FormValue("newName")
	idtheme := r.FormValue("idtheme")
	fmt.Println(newName, idtheme)
	vars := mux.Vars(r)
	_, err := h.DB.Exec(
		"update amaker.theme set name = ? where idtheme = ?",
		newName,
		idtheme,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(
		w,
		r,
		"/branches/"+vars["username"]+"/"+vars["idBranch"],
		http.StatusFound,
	)
}

func (h *Handler) DeleteBranch(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DB")
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
	fmt.Println("CB")
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

	idtask, err := strconv.Atoi(r.FormValue("idtask"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := h.findUser(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	row := h.DB.QueryRow(
		"select count(*) from amaker.answer as a where a.idsender = ? and a.idtask =  ?",
		user.Iduser,
		idtask,
	)
	var count int
	err = row.Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count == 0 {
		answer := &mod.Answer{}

		answer.File = r.FormValue("file")
		answer.Idsender = user.Iduser
		answer.Idtask = idtask
		answer.Status = 0
		answer.Date = time.Now()
		_, err = h.DB.Exec(
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
	} else {
		_, err = h.DB.Exec(
			"update amaker.answer set status = 0, file = ?, date = ? where idtask = ? and idsender = ?",
			r.FormValue("file"),
			time.Now(),
			idtask,
			user.Iduser,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	vars := mux.Vars(r)
	http.Redirect(w, r, "/branches/worker/"+vars["username"]+"/"+vars["idBranch"], http.StatusFound)
}

func (h *Handler) GetWorkerBranches(w http.ResponseWriter, r *http.Request) {
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

	rows, err := h.DB.Query(
		"select br.idbranch, br.name, br.idcompany from amaker.branch as br inner join amaker.access as r on br.idbranch = r.idbranch where r.iduser = ?",
		user.Iduser,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	brMap := make(map[int]*mod.Branch)

	for rows.Next() {
		br := &mod.Branch{}
		err := rows.Scan(&br.Idbranch, &br.Name, &br.Idcompany)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		br.Username = username
		brMap[br.Idbranch] = br
	}

	branches := []*mod.Branch{}
	for i, _ := range brMap {
		branches = append(branches, brMap[i])
	}

	data := &BranchesData{
		Branches:  branches,
		UName:     username,
		IsCompany: role,
	}

	h.Tmpl.ExecuteTemplate(w, "workerBranches.html", data)
}

func (h *Handler) GetWorkerCurrentBranch(w http.ResponseWriter, r *http.Request) {
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

	vars := mux.Vars(r)
	row := h.DB.QueryRow(
		"select * from amaker.branch as b where b.idbranch = ?",
		vars["idBranch"],
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
		vars["idBranch"],
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	themeMap := make(map[int]*mod.Theme)
	taskMap := make(map[int]*mod.Task)
	flag := false
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
				Idtask:   int(ctat.Idtask.Int32),
				Name:     ctat.TaskName.String,
				Text:     ctat.TaskText.String,
				Idtheme:  ctat.Idtheme,
				Answer:   answer,
				Username: username,
				Idbranch: branch.Idbranch,
			}
			if _, ok := taskMap[task.Idtask]; !ok {
				taskMap[task.Idtask] = task
				flag = true
			} else {
				flag = false
			}
		}
		if val, ok := themeMap[ctat.ThemeIndex]; ok {
			if flag {
				val.Tasks = append(val.Tasks, task)
			}
		} else {
			theme := &mod.Theme{
				Idtheme:  ctat.Idtheme,
				Name:     ctat.ThemeName,
				Idbranch: branch.Idbranch,
				Index:    ctat.ThemeIndex,
				Tasks:    []*mod.Task{},
				Username: username,
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

	h.Tmpl.ExecuteTemplate(w, "workerTasks.html", branchData)
}
