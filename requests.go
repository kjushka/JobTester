package main

import (
	mod "./model"
	"github.com/gorilla/mux"
	"net/http"
)

type CustomRequest struct {
	Request *mod.Request
	User    *mod.User
}

type CustomRequestData struct {
	CRs      []*CustomRequest
	Username string
}

//все запросы для компании
func (h *Handler) GetRequests(w http.ResponseWriter, r *http.Request) {
	usernameCompany, role, status := h.checkCookie(w, r)

	if !status {
		return
	}

	if role != 1 {
		http.Redirect(w, r, "/home/worker/"+usernameCompany, http.StatusFound)
		return
	}

	company, err := h.findUser(usernameCompany)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	arrayCustomRequest := []*CustomRequest{}

	rows, err := h.DB.Query(
		"select r.idrequest, r.iduser, r.idbranch, u.iduser, u.username, u.password, u.name, u.isCompany "+
			"from amaker.request as r "+
			"inner join (select b.idbranch from amaker.branch as b where b.idcompany = ?) as t "+
			"on t.idbranch = r.idbranch "+
			"left join amaker.user as u "+
			"on u.iduser = r.iduser",
		company.Iduser,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	for rows.Next() {
		element := &mod.Request{}
		user := &mod.User{}
		err := rows.Scan(
			&element.Idrequest,
			&element.Iduser,
			&element.Idbranch,
			&user.Iduser,
			&user.Username,
			&user.Password,
			&user.Name,
			&user.IsCompany,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cr := &CustomRequest{
			Request: element,
			User:    user,
		}
		arrayCustomRequest = append(arrayCustomRequest, cr)
	}

	crd := &CustomRequestData{
		CRs:      arrayCustomRequest,
		Username: usernameCompany,
	}

	h.Tmpl.ExecuteTemplate(w, "hiring.html", crd)
}

func (h *Handler) AcceptRequest(w http.ResponseWriter, r *http.Request) {
	username, role, status := h.checkCookie(w, r)

	if !status {
		return
	}

	if role != 1 {
		http.Redirect(w, r, "/home/worker/"+username, http.StatusFound)
		return
	}

	vars := mux.Vars(r)
	_, err := h.DB.Exec(
		"delete from amaker.request as r where r.iduser = ? and r.idbranch = ?",
		vars["iduser"],
		vars["idbranch"],
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = h.DB.Exec(
		"insert into amaker.access (`iduser`, `idbranch`) values (?, ?)",
		vars["iduser"],
		vars["idbranch"],
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/requests/"+username, http.StatusFound)
}

func (h *Handler) DeclineRequest(w http.ResponseWriter, r *http.Request) {
	username, role, status := h.checkCookie(w, r)
	if !status {
		return
	}

	if role != 1 {
		http.Redirect(w, r, "/home/worker/"+username, http.StatusFound)
		return
	}

	vars := mux.Vars(r)
	_, err := h.DB.Exec(
		"delete from amaker.request as r where r.iduser = ? and r.idbranch = ?",
		vars["iduser"],
		vars["idbranch"],
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/requests/"+username, http.StatusFound)
}

//отдельно запрос по id
/*func (h *Handler) GetRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idPretendent, err := strconv.ParseInt(vars["id"], 10, 32)

	if err == nil {
		fmt.Printf("Error formatter from string to int in copanies %s", vars["id"])
	}

	rows := h.DB.QueryRow("SELECT AU.name FROM amake.user AS AU WHERE AR.idUser=?",
		idPretendent,
	)

	element := &mod.User{}
	err = rows.Scan(
		&element.Name,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//вывод
}

func (h *Handler) SetBranchAccess(w http.ResponseWriter, r *http.Request) {
	nameCompany, role, status := h.checkCookie(w, r)

	if !status {
		return
	}

	if role != 0 {
		http.Redirect(w, r, "/requests/"+nameCompany, http.StatusFound)
		return
	}

	strAccess := "1"
	access, err := strconv.ParseInt(strAccess, 10, 32)

	if err == nil {
		fmt.Printf("Error formatter from string to int %s", strAccess)
	}

	idCompany, err := h.findUser(nameCompany)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	idUser, err := strconv.ParseInt(vars["id"], 10, 32)

	if err == nil {
		fmt.Printf("Error formatter from string to int %s", vars["id"])
	}

	if access == 1 {
		_, err := h.DB.Exec(
			"INSERT INTO amaker.access ('iduser','idcom', 'idbranch') VALUES (?,?,?)",
			idUser, idCompany,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {

	}

}

func (h *Handler) GetCurrentWorker(w http.ResponseWriter, r *http.Request) {

}*/
