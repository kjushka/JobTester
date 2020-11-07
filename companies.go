package main

import (
	mod "./model"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type BranchesCompany struct {
	Branches []*mod.Branch
	Company  *mod.User
}

//Получение всех компаний
func (h *Handler) GetCompanies(w http.ResponseWriter, r *http.Request) {
	username, role, status := h.checkCookie(w, r)

	if !status {
		return
	}

	if role != 0 {
		http.Redirect(w, r, "/home/companies/"+username, http.StatusFound)
		return
	}

	rows, err := h.DB.Query(
		"SELECT * FROM amaker.user AS AU WHERE AU.isCompany=1",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	companiesArray := []*mod.User{}

	for rows.Next() {
		element := &mod.User{}
		err := rows.Scan(&element.Iduser,
			&element.Username,
			&element.Password,
			&element.Name,
			&element.IsCompany,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		companiesArray = append(companiesArray, element)
	}
	h.Tmpl.ExecuteTemplate(w, "companies.html", companiesArray)
}

//инф о конкретной компании
func (h *Handler) GetCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyName := vars["username"]

	branch, errbr := h.DB.Query(
		"SELECT AB.idbranch, AB.name FROM amaker.branch AS AB WHERE AB.idcompany=(SELECT AU.iduser FROM amaker.user AS AU WHERE AU.name=?)",
		companyName,
	)

	if errbr != nil {
		return
	}

	defer branch.Close()
	branchesArray := []*mod.Branch{}

	for branch.Next() {
		element := &mod.Branch{}
		err := branch.Scan(&element.Idbranch,
			&element.Name,
			&element.Idcompany,
		)

		if err != nil {
			http.Error(w, errbr.Error(), http.StatusInternalServerError)
			return
		}
		branchesArray = append(branchesArray, element)
	}

	companyRow, er := h.findUser(companyName)
	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
		return
	}

	branchComp := &BranchesCompany{
		Branches: branchesArray,
		Company:  companyRow,
	}

	h.Tmpl.ExecuteTemplate(w, "company.html", branchComp)
}

//Заявка на прохождение тестирования
func (h *Handler) SendRequest(w http.ResponseWriter, r *http.Request) {
	username, role, status := h.checkCookie(w, r)

	if !status {
		return
	}

	if role != 0 {
		http.Redirect(w, r, "/home/companies/"+username, http.StatusFound)
		return
	}

	userInfom, err := h.findUser(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mapfromPath := mux.Vars(r)
	idBranch, err := strconv.ParseInt(mapfromPath["idBranch"], 10, 32)

	if err == nil {
		fmt.Printf("Error formatter from string to int in copanies %s", mapfromPath["idBranch"])
	}

	answerFromCompany := h.DB.QueryRow(
		"SELECT COUNT(*) FROM amaker.request AS AR WHERE AR.idbranch=? AND AR.iduser=?",
		idBranch,
		userInfom.Iduser,
	)

	var count int
	err = answerFromCompany.Scan(&count)

	if count != 0 {
		_, err := h.DB.Exec(
			"INSERT INTO amaker.request ('idbranch','iduser') VALUES (?,?)",
			idBranch,
			userInfom.Iduser,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
