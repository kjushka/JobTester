package main

import (
	mod "./model"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type BranchesNew struct {
	Branch *mod.Branch
	Count  int
}

type BranchesCompany struct {
	Branches []*BranchesNew
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
	h.Tmpl.ExecuteTemplate(w, "company.html", companiesArray)
}

//инф о конкретной компании
func (h *Handler) GetCompany(w http.ResponseWriter, r *http.Request) {
	count := 0
	vars := mux.Vars(r)

	companyName, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		fmt.Printf("Error formatter from string to int in copanies %s", vars["id"])
	}

	branch, errbr := h.DB.Query(
		"SELECT AB.idbranch, AB.name FROM amaker.branch AS AB WHERE AB.idcompany=?",
		companyName,
	)

	if errbr != nil {
		fmt.Print("errooooooooooo")
		return
	}

	defer branch.Close()
	branchesArray := []*BranchesNew{}

	for branch.Next() {

		element := &mod.Branch{}
		err := branch.Scan(&element.Idbranch,
			&element.Name,
		)

		if err != nil {
			return
		}
		fmt.Println(element.String())
		branchesArray = append(branchesArray, &BranchesNew{
			Branch: element,
			Count:  count,
		})
		count++
	}

	companyRo := h.DB.QueryRow(
		"SELECT * FROM amaker.user AS AU WHERE AU.iduser=?",
		companyName,
	)

	company := &mod.User{}

	err = companyRo.Scan(
		&company.Iduser,
		&company.Username,
		&company.Password,
		&company.Name,
		&company.IsCompany,
	)
	fmt.Println(company.String())

	branchComp := &BranchesCompany{
		Branches: branchesArray,
		Company:  company,
	}

	h.Tmpl.ExecuteTemplate(w, "comp.html", branchComp)
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
