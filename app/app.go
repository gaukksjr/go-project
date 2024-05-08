package app

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"go/project/handlers"
	"go/project/repository"

	"github.com/gorilla/mux"
)

type Application struct {
}

func NewApp() *Application {
	return &Application{}
}

var (
	Repo *repository.Repository
	Hand *handlers.Handler
)

func init() {
	connectionString := os.Getenv("MYSQL_CONNECTION_STRING")
	if connectionString == "" {
		panic("empty env")
	}

	Repo = repository.NewRepository(connectionString)
	Hand = handlers.NewHandler(Repo)
}

func (a *Application) StartServer() {
	r := mux.NewRouter()

	tmpl := template.Must(template.ParseGlob("static/*.html"))

	r.HandleFunc("/home", Hand.HomeHandler)
	r.HandleFunc("/menu", Hand.MenuHandler)
	r.HandleFunc("/register", Hand.RegisterHandler)
	r.HandleFunc("/login", Hand.LoginHandler)
	r.HandleFunc("/add-menu", Hand.AddMenuItemsHandler)
	r.HandleFunc("/save-order", Hand.SaveOrderHandler)
	r.HandleFunc("/logout", Hand.LogoutHandler)
	r.HandleFunc("/update-order", Hand.UpdateOrderStatusHandler)
	r.HandleFunc("/get-order", Hand.GetOrdersHandler)
	r.HandleFunc("/get-username", Hand.GetUsernameByIDHandler)

	r.HandleFunc("/login-page", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "login.html", nil)
	})
	r.HandleFunc("/register-page", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "register.html", nil)
	})
	r.HandleFunc("/home-page", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "home.html", nil)
	})
	r.HandleFunc("/menu-page", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "menu.html", nil)
	})
	r.HandleFunc("/order-page", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "order.html", nil)
	})
	http.Handle("/", r)

	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      r,
	}

	log.Println("Listening on port 8080...")
	log.Fatal(server.ListenAndServe())
}
