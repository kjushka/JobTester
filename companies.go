package main

import (
	mod "./model"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type BranchesCompany struct {
	Branches []*mod.Branch
	Company  *mod.User
}

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
		"SELECT * FROM amaker.user WHERE user.isCompany=1",
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

//инф о конкретоной компании
func (h *Handler) GetCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idcompany, err := strconv.ParseInt(vars["id"], 10, 32)

	branch, errbr := h.DB.Query(
		"SELECT AB.idbranch, AB.name FROM amaker.branch AS AB WHERE AB.idcompany=?",
		idcompany,
	)

	if errbr != nil {

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

	companyRow := h.DB.QueryRow(
		"SELECT  * FROM amaker.user AS AU WHERE AU.iduser=?",
		idcompany,
	)

	company := &mod.User{}
	err = companyRow.Scan(&company.Iduser,
		&company.Username,
		&company.Password,
		&company.Name,
		&company.IsCompany,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	branchComp := &BranchesCompany{
		Branches: branchesArray,
		Company:  company,
	}

	h.Tmpl.ExecuteTemplate(w, "company.html", branchComp)
}

func (h *Handler) SendRequest(w http.ResponseWriter, r *http.Request) {
	username, role, status := h.checkCookie(w, r)
	statusForRequest := 0

	if !status {
		return
	}

	if role != 0 {
		http.Redirect(w, r, "/home/companies/"+username, http.StatusFound)
		return
	}

	mapfromPath := mux.Vars(r)
	idcompany, err := strconv.ParseInt(mapfromPath["id"], 10, 32)

	if err == nil {
		fmt.Printf("Error formatter from string to int in copanies %s", mapfromPath["id"])
	}

	companyRow := h.DB.QueryRow(
		"SELECT  AU.iduser FROM amaker.user AS AU WHERE AU.username=?",
		username,
	)

	idUser := &mod.User{}
	err = companyRow.Scan(&idUser.Iduser)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	answerFromCompany := h.DB.QueryRow(
		"SELECT  AR.status FROM amaker.request AS AR WHERE AR.idcomp=? AND AR.iduser=? AND AR.status IN (0,1)",
		idcompany,
		idUser,
	)

	statusFromDB := &mod.Request{}
	err = answerFromCompany.Scan(&statusFromDB.Status)

	if err != nil {
		if err != sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			_, err := h.DB.Exec(
				"INSERT INTO amaker.request ('idcomp','iduser', 'status') VALUES (?,?,?)",
				idcompany,
				idUser,
				statusForRequest,
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	} else {
		if statusFromDB.Status < 2 {
			w.Write([]byte("Sorry no vi uge otpravily zapros"))
		}
	}

}
