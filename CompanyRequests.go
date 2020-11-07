package main

//import (
//	mod "./model"
//	"fmt"
//	"github.com/gorilla/mux"
//	"net/http"
//	"strconv"
//)
//
//type CustomRequest struct{
//	IdUser int
//	NameUser int
//	Status int
//}
//
////все запросы для компании
//func (h *Handler) GetRequests(w http.ResponseWriter, r *http.Request){
//	usernameCompany, role, status := h.checkCookie(w, r)
//
//	if !status {
//		return
//	}
//
//	if role != 0 {
//		http.Redirect(w, r, "/requests/"+usernameCompany, http.StatusFound)
//		return
//	}
//
//	idCompany, err:= h.findUser(usernameCompany)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	rows, err:= h.DB.Query("SELECT AU.iduser, AU.name, AB.status FROM amaker.user AS AU," +
//		" amaker.branch AS AB WHERE AU.name=(SELECT AB.name FROM amaker.branch AS AB WHERE AB.idcompany=?)",
//		idCompany,
//	)
//
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	defer rows.Close()
//	arrayCustomRequest := []*CustomRequest{}
//
//	for rows.Next(){
//		element:= &CustomRequest{}
//		err:= rows.Scan(
//			&element.IdUser,
//			&element.NameUser,
//			&element.Status,)
//
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		arrayCustomRequest = append(arrayCustomRequest, element)
//	}
//	//h.Tmpl.ExecuteTemplate(w, "requests",allCompanyRequestsArray)
//}
//
////отдельно запрос по id
//func (h *Handler) GetRequest(w http.ResponseWriter, r *http.Request){
//	vars := mux.Vars(r)
//	idPretendent, err:= strconv.ParseInt(vars["id"], 10, 32)
//
//	if err == nil {
//		fmt.Printf("Error formatter from string to int in copanies %s", vars["id"])
//	}
//
//	rows:= h.DB.QueryRow("SELECT AU.name FROM amake.user AS AU WHERE AR.idUser=?",
//		idPretendent,
//	)
//
//	element:= &mod.User{}
//	err= rows.Scan(
//		&element.Name,
//		)
//
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	//вывод
//}
//
//
//func (h *Handler) SetBranchAccess(w http.ResponseWriter, r *http.Request) {
//	nameCompany, role, status := h.checkCookie(w, r)
//
//	if !status {
//		return
//	}
//
//	if role != 0 {
//		http.Redirect(w, r, "/requests/"+nameCompany, http.StatusFound)
//		return
//	}
//
//	strAccess := "1"
//	access, err := strconv.ParseInt(strAccess, 10, 32)
//
//	if err == nil {
//		fmt.Printf("Error formatter from string to int %s", strAccess)
//	}
//
//	idCompany, err:= h.findUser(nameCompany)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	vars := mux.Vars(r)
//	idUser, err := strconv.ParseInt(vars["id"], 10, 32)
//
//	if (err == nil) {
//		fmt.Printf("Error formatter from string to int %s", vars["id"])
//	}
//
//	if (access == 1) {
//		_, err := h.DB.Exec(
//			"INSERT INTO amaker.access ('iduser','idcom', 'idbranch') VALUES (?,?,?)",
//			idUser, idCompany, idBranch,
//		)
//
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//	} else{
//
//
//
//	}
//
//}
//
//func (h *Handler) GetCurrentWorker(w http.ResponseWriter, r *http.Request){
//
//}
