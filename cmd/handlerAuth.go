package cmd

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"dictionary/pkg/pass"
)

// .EXT.

// 1. Обработчик для загрузки формы авторизации MethodGet
func HandlerAuth(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode") // mode (register, login, logout)
	dataAuth := DataAuth{
		Title:      "Register",
		ButtonText: "Register",
		Mode:       mode,
	}

	switch mode {
	case "login":
		dataAuth.Title = "Login"
		dataAuth.ButtonText = "Login"
	case "logout":
		dataAuth.Title = "Logout"
		dataAuth.ButtonText = "Logout"
	}

	data := TemplateData{
		Auth: dataAuth,
	}

	tmpl, err := template.ParseFiles("ui/base.html", "ui/auth.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

// 2. Обработчик для регистрации нового пользователя MethodPost
func HandlerAuthRegister(w http.ResponseWriter, r *http.Request) {

	log.Println("authRegisterHandler")
	var user map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request json data", http.StatusBadRequest)
		return
	}

	// Проверка сложности пароля
	if err := pass.ValidatePassword(user["password"].(string)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Хеширование пароля
	hashedPassword, err := pass.HashPassword(user["password"].(string))
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user["password"] = hashedPassword

	// Загружаем данные из базы
	id, err := SRV.dbPool.InsertUsers(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if id == 0 {
		http.Error(w, "Registration error", http.StatusBadRequest)
		return
	}

	//data := map[string]string{}
	data := DataUser{}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// 3. Обработчик для входа в систему MethodPut
func HandlerAuthLogin(w http.ResponseWriter, r *http.Request) {
	session := GetOrCreateSession(w, r)
	session.Mutex.Lock()
	defer session.Mutex.Unlock()

	var user map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request json data", http.StatusBadRequest)
		return
	}

	mail, ok := user["mail"].(string)
	if !ok {
		http.Error(w, "Invalid request load: mail", http.StatusBadRequest)
		return
	}

	password, ok := user["password"].(string)
	if !ok {
		http.Error(w, "Invalid request load: password", http.StatusBadRequest)
		return
	}

	// Загружаем данные из базы
	usersDb, err := SRV.dbPool.ViewUsers(map[string]interface{}{
		"mail": mail,
	})
	if err != nil || len(usersDb) == 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Хеширование пароля
	ok = pass.CheckPasswordHash(password, usersDb[0].Password)
	if !ok {
		http.Error(w, "The password is incorrect", http.StatusInternalServerError)
		return
	}

	// Если авторизация успешна
	session.UserLogin = true
	session.User.Id = usersDb[0].Id
	session.User.Name = usersDb[0].Name
	session.User.Mail = usersDb[0].Mail

	data := TemplateData{
		User: session.User,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)

}

// 4. Обработчик для выхода из системы MethodDelete
func HandlerAuthLogout(w http.ResponseWriter, r *http.Request) {
	session := GetOrCreateSession(w, r)
	session.Mutex.Lock()
	defer session.Mutex.Unlock()

	var user interface{}
	json.NewDecoder(r.Body).Decode(&user)

	//data := map[string]string{}

	session.UserLogin = false
	session.User = DataUser{}

	data := TemplateData{
		User: session.User,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
