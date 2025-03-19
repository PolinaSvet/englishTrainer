package cmd

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"dictionary/pkg/cardlern"
	"dictionary/pkg/storage/postgres"

	"github.com/gorilla/mux"
)

type Session struct {
	Cards     *cardlern.Cards
	Mutex     sync.RWMutex
	Timer     *time.Timer
	ErrorMsg  error
	User      DataUser
	UserLogin bool
}

type TemplateData struct {
	Cards   *cardlern.Cards `json:"cards"`
	Message string          `json:"message"`
	User    DataUser        `json:"user"`
	Auth    DataAuth        `json:"auth"`
}

type DataAuth struct {
	Title      string `json:"title"`
	ButtonText string `json:"buttonText"`
	Mode       string `json:"mode"`
}

type DataUser struct {
	Name string `json:"name"`
	Mail string `json:"mail"`
	Id   int    `json:"id"`
}

type server struct {
	sessions map[string]*Session
	mutex    sync.RWMutex
	dbPool   *postgres.Storage
}

var SRV server

func InitData() {
	connString := os.Getenv("PG_URL_DBLANGUAGE")

	dbPool, err := postgres.New(connString)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
	SRV.dbPool = dbPool
	SRV.sessions = make(map[string]*Session)
	log.Println("InitData")
}

func Handler() {
	router := mux.NewRouter()
	router.HandleFunc("/", HandlerHome).Methods(http.MethodGet)
	router.HandleFunc("/cards", HandlerCards).Methods(http.MethodGet)
	router.HandleFunc("/cards", HandlerCardsFix).Methods(http.MethodPost)

	// Обновление роутов
	router.HandleFunc("/auth", HandlerAuth).Methods(http.MethodGet)
	router.HandleFunc("/auth/register", HandlerAuthRegister).Methods(http.MethodPost)
	router.HandleFunc("/auth/login", HandlerAuthLogin).Methods(http.MethodPut)
	router.HandleFunc("/auth/logout", HandlerAuthLogout).Methods(http.MethodDelete)

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func HandlerHome(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("ui/base.html", "ui/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func GetOrCreateSession(w http.ResponseWriter, r *http.Request) *Session {
	sessionID, err := r.Cookie("session_id")
	if err != nil {
		sessionID = &http.Cookie{
			Name:  "session_id",
			Value: GenerateSessionID(),
		}
		http.SetCookie(w, sessionID)
	}

	SRV.mutex.Lock()
	defer SRV.mutex.Unlock()

	session, exists := SRV.sessions[sessionID.Value]
	if !exists {
		session = &Session{
			Cards: nil,
			Timer: time.AfterFunc(30*time.Minute, func() {
				SRV.mutex.Lock()
				delete(SRV.sessions, sessionID.Value)
				SRV.mutex.Unlock()
				log.Println("Сессия удалена:", sessionID.Value)
			}),
			ErrorMsg:  nil,
			UserLogin: false,
		}
		SRV.sessions[sessionID.Value] = session
		log.Println("Новая сессия создана:", sessionID.Value)
	}

	log.Println(sessionID, session.User, session.UserLogin)
	return session
}

func GenerateSessionID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
