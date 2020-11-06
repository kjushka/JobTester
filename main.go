package main

import (
	mod "./model"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type Response struct {
	Status  bool
	Message string
}

type Handler struct {
	DB   *sql.DB
	Tmpl *template.Template
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	row := h.DB.QueryRow(
		"SELECT u.iduser, u.username, u.password, u.name, u.isCompany FROM amaker.user as u where u.username = ? and  u.password = ?",
		r.FormValue("username"),
		r.FormValue("password"),
	)
	user := &mod.User{}
	err := row.Scan(
		&user.Iduser,
		&user.Username,
		&user.Password,
		&user.Name,
		&user.IsCompany,
	)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var response *Response
	if err != nil {
		response = &Response{
			Status:  false,
			Message: "Incorrect username or password!",
		}
	} else {
		response = &Response{
			Status:  true,
			Message: "Logged successfully!",
		}
	}
	jsonResponse, marshalError := json.Marshal(response)
	if marshalError != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	setCookie(w, user.Username, strconv.Itoa(user.IsCompany))
	w.Write(jsonResponse)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	fmt.Println(username)
	_, err := h.findUser(username)
	if err == sql.ErrNoRows {
		response := &Response{
			Status:  false,
			Message: "This username already registered",
		}
		jsonResponse, marshalError := json.Marshal(response)
		if marshalError != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
		return
	}
	role := r.FormValue("isCompany")
	intRole, err := strconv.Atoi(role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result, err := h.DB.Exec(
		"INSERT INTO amaker.user (`username`, `password`, `name`, `isCompany`) VALUES (?, ?, ?, ?)",
		username,
		r.FormValue("password"),
		r.FormValue("name"),
		intRole,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	affected, err := result.RowsAffected()
	lastID, err := result.LastInsertId()
	fmt.Println("Registered", affected,
		"with id:", lastID)
	setCookie(w, username, role)
	response := &Response{
		Status:  true,
		Message: "Registered successfully!",
	}
	jsonResponse, marshalError := json.Marshal(response)
	if marshalError != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (h *Handler) findUser(username string) (*mod.User, error) {
	fmt.Println(username)
	row := h.DB.QueryRow(
		"SELECT u.id, u.username, u.password, u.name, u.isCompany FROM amaker.user as u where u.username = ?",
		username,
	)
	user := &mod.User{}
	err := row.Scan(
		&user.Iduser,
		&user.Username,
		&user.Password,
		&user.Name,
		&user.IsCompany,
	)
	return user, err
}

func setCookie(w http.ResponseWriter, username string, isCompany string) {
	expiration := time.Now().Add(time.Minute * 30)
	cookie := &http.Cookie{
		Name:    "user-name",
		Value:   username,
		Expires: expiration,
	}
	http.SetCookie(w, cookie)
	cookie = &http.Cookie{
		Name:    "user-role",
		Value:   isCompany,
		Expires: expiration,
	}
	http.SetCookie(w, cookie)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("user-name")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/home", http.StatusFound)
		return
	}
	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	session, err = r.Cookie("user-role")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/home", http.StatusFound)
		return
	}
	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	http.Redirect(w, r, "/home", http.StatusFound)
}

func checkCookie(w http.ResponseWriter, r *http.Request) (string, int, bool) {
	cName, err := r.Cookie("user-number")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/home", http.StatusFound)
		return "", 0, false
	}
	cRole, err := r.Cookie("user-role")
	if err == nil {
		http.Redirect(w, r, "/home", http.StatusFound)
		return "", 0, false
	}
	role, atoiErr := strconv.Atoi(cRole.Value)
	if atoiErr != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", 0, false
	}
	return cName.Value, role, true
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "home.html", struct{}{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetWorker(w http.ResponseWriter, r *http.Request) {
	username, role, status := checkCookie(w, r)
	if !status {
		return
	}
	if role != 0 {
		http.Redirect(w, r, "/home/company/"+username, http.StatusFound)
	}
	user, err := h.findUser(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.Tmpl.ExecuteTemplate(w, "worker.html", user)
}

func (h *Handler) GetCompany(w http.ResponseWriter, r *http.Request) {
	username, role, status := checkCookie(w, r)
	if !status {
		return
	}
	if role != 1 {
		http.Redirect(w, r, "/home/worker/"+username, http.StatusFound)
	}
	company, err := h.findUser(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.Tmpl.ExecuteTemplate(w, "company.html", company)
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

func (h *Handler) GetCompanies(w http.ResponseWriter, r *http.Request) {

}

func main() {
	//dsn := "root:pass@tcp(localhost:3306)/amaker?" +
	//	"&charset=utf8&interpolateParams=true"
	//db, err := sql.Open("mysql", dsn)
	//db.Ping()
	//if err != nil {
	//	panic(err)
	//}
	//defer db.Close()
	//fmt.Println("Connected to db")

	handler := &Handler{
		//DB:   db,
		Tmpl: template.Must(template.ParseGlob("./webapp/view/*")),
	}
	router := mux.NewRouter()
	router.HandleFunc("/login", handler.Login).Methods("POST")
	router.HandleFunc("/register", handler.Register).Methods("POST")
	router.HandleFunc("/logout", handler.Logout).Methods("GET")

	router.HandleFunc("/home", handler.Home).Methods("GET")
	router.HandleFunc("/home/worker/{username}", handler.GetWorker).Methods("GET")
	//router.HandleFunc("/home/worker/{username}", handler.EditWorker).Methods("POST")
	router.HandleFunc("/home/company/{username}", handler.GetCompany).Methods("GET")
	//router.HandleFunc("/home/company/{username}", handler.EditCompany).Methods("POST")

	//router.HandleFunc("/branches/{username}", handler.GetBranches).Methods("GET")
	//router.HandleFunc("/branches/{username}/{id}", handler.GetBranch).Methods("GET")
	//router.HandleFunc("/branches/{username}/{id}", handler.EditBranch).Methods("POST")
	//router.HandleFunc("/branches/{username}/{id}", handler.DeleteBranch).Methods("DELETE")
	//router.HandleFunc("/branches/{username}/create", handler.CreateBranch).Methods("POST")
	//
	//router.HandleFunc("{username}/answers", handler.GetAnswers).Methods("GET")
	//router.HandleFunc("{username}/answers/{id}", handler.GetAnswer).Methods("GET")
	//router.HandleFunc("{username}/answers/{id}/status", handler.SetAnswerStatus).Methods("POST")
	//router.HandleFunc("{username}/answers/{id}/download", handler.DownloadAnswer).Methods("GET")
	//
	//router.HandleFunc("/companies", handler.GetCompanies).Methods("GET")
	//router.HandleFunc("/companies/{id}", handler.GetCompany).Methods("GET")
	//router.HandleFunc("/companies/{id}/request", handler.SendRequest).Methods("POST")

	//router.HandleFunc("/requests/{username}", handler.GetRequests).Methods("GET")
	//router.HandleFunc("/requests/{username}/{id}", handler.GetRequest).Methods("GET")
	//router.HandleFunc("/requests/{username}/{id}/access", handler.SetBranchAccess).Methods("POST")
	fmt.Println("Listen port 8080")
	http.ListenAndServe(":8080", router)
}
