package main

import (
	mod "./model"
	"github.com/gorilla/mux"
	"strconv"
	"time"

	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"net/http"
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
	fmt.Println("Log")

	row := h.DB.QueryRow(
		"SELECT * FROM amaker.user as u where u.username = ? and  u.password = ?",
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err != nil {
		err := h.Tmpl.ExecuteTemplate(w, "home.html", LogErr{
			"Incorrect username or password!",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	setCookie(w, user.Username, strconv.Itoa(user.IsCompany))

	if user.IsCompany == 0 {
		http.Redirect(w, r, "/home/worker/"+user.Username, 302)
	} else {
		http.Redirect(w, r, "/home/company/"+user.Username, 302)
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Reg")
	username := r.FormValue("username")
	_, err := h.findUser(username)
	if err != nil {
		if err != sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			err = nil
		}
	} else {
		err := h.Tmpl.ExecuteTemplate(w, "home.html", LogErr{
			"This username already registered",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	role := r.FormValue("role")
	intRole := 0
	if role == "on" {
		intRole = 1
	}
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
	setCookie(w, username, strconv.Itoa(intRole))
	if intRole == 0 {
		http.Redirect(w, r, "/home/worker/"+username, 302)
	} else {
		http.Redirect(w, r, "/home/company/"+username, 302)
	}
}

func (h *Handler) findUser(username string) (*mod.User, error) {
	row := h.DB.QueryRow(
		"SELECT * FROM amaker.user as u where u.username = ?",
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

	fmt.Println("Logout")

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

func (h *Handler) checkCookie(w http.ResponseWriter, r *http.Request) (string, int, bool) {
	cName, err := r.Cookie("user-name")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/home", http.StatusFound)
		return "", -1, false
	}
	cRole, err := r.Cookie("user-role")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/home", http.StatusFound)
		return "", -1, false
	}
	cName.Expires = time.Now().Add(time.Minute * 30)
	cRole.Expires = time.Now().Add(time.Minute * 30)
	role, atoiErr := strconv.Atoi(cRole.Value)
	if atoiErr != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "error", -1, false
	}
	return cName.Value, role, true
}

func main() {
	dsn := "root:pass@tcp(localhost:3306)/amaker?" +
		"&charset=utf8&interpolateParams=true&parseTime=true"
	db, err := sql.Open("mysql", dsn)
	db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("Connected to db")

	handler := &Handler{
		DB:   db,
		Tmpl: template.Must(template.ParseGlob("./webapp/view/*")),
	}
	router := mux.NewRouter()
	router.HandleFunc("/login", handler.Login).Methods("POST")
	router.HandleFunc("/register", handler.Register).Methods("POST")
	router.HandleFunc("/logout", handler.Logout).Methods("GET")

	router.HandleFunc("/home", handler.Home).Methods("GET")
	router.HandleFunc("/home/worker/{username}", handler.HomeWorker).Methods("GET")
	router.HandleFunc("/home/worker/{username}/edit", handler.EditWorker).Methods("POST")
	router.HandleFunc("/home/company/{username}", handler.HomeCompany).Methods("GET")
	router.HandleFunc("/home/company/{username}/edit", handler.EditCompany).Methods("POST")

	router.HandleFunc("/branches/{username}", handler.GetBranches).Methods("GET")
	router.HandleFunc("/branches/{username}/{id:[0-9]+}", handler.GetBranch).Methods("GET")
	router.HandleFunc("/branches/{username}/{id:[0-9]+}", handler.EditBranch).Methods("POST")
	router.HandleFunc("/branches/{username}/{id:[0-9]+}/delete", handler.DeleteBranch).Methods("GET")
	router.HandleFunc("/branches/{username}/create", handler.CreateBranch).Methods("POST")
	router.HandleFunc(
		"/branches/{username}/{id}/send-answer/{idtask}",
		handler.SendAnswer,
	).Methods("POST")

	router.HandleFunc("/answers/{username}", handler.GetAnswers).Methods("GET")
	router.HandleFunc("/answers/{username}/{id}", handler.GetAnswer).Methods("GET")
	router.HandleFunc("/answers/{username}/{id}/status", handler.SetAnswerStatus).Methods("POST")
	router.HandleFunc("/answers/{username}/{id}/download", handler.DownloadAnswer).Methods("GET")

	/*router.HandleFunc("/companies", handler.GetCompanies).Methods("GET")
	router.HandleFunc("/companies/{id}", handler.GetCompany).Methods("GET")
	router.HandleFunc("/companies/{id}/request", handler.SendRequest).Methods("POST")

	router.HandleFunc("/requests/{username}", handler.GetRequests).Methods("GET")
	router.HandleFunc("/requests/{username}/{id}", handler.GetRequest).Methods("GET")
	router.HandleFunc("/requests/{username}/{id}/access", handler.SetBranchAccess).Methods("POST")*/

	router.HandleFunc("/resources/home.js", handler.SendHomeJs)
	router.HandleFunc("/resources/css.css", handler.SendCssCss)
	router.HandleFunc("/resources/branches.js", handler.SendBranchesJs)
	router.HandleFunc("/resources/branches.css", handler.SendBranchesCss)
	router.HandleFunc("/home/resources/branches.css", handler.SendBranchesCss)

	fmt.Println("Listen port 8080")
	http.ListenAndServe(":8080", router)
}
