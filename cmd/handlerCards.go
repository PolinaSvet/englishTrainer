package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"dictionary/pkg/cardlern"
)

// .EXT.

// 1. Получили запрос на начало изучения карточек
func HandlerCards(w http.ResponseWriter, r *http.Request) {
	session := GetOrCreateSession(w, r)

	// Создаем контекст с таймаутом в 2 секунды
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Канал для получения карточек и ошибки
	cardsChan := make(chan *cardlern.Cards, 1)
	errChan := make(chan error, 1)

	// Запускаем инициализацию карточек в отдельной горутине
	go func() {
		cards, err := initializeCards(map[string]interface{}{
			"limit": 5,
		})
		if err != nil {
			errChan <- err
			return
		}
		cardsChan <- cards
	}()

	// Ожидаем завершения инициализации или таймаута. Отправляем данные клиенту
	select {
	case cards := <-cardsChan:
		session.Cards = cards
	case err := <-errChan:
		session.Cards = cardlern.GetEmptyCard()
		session.ErrorMsg = err
	case <-ctx.Done():
		session.Cards = cardlern.GetEmptyCard()
		session.ErrorMsg = fmt.Errorf("Таймаут инициализации карточек")
	}

	// Отправляем данные клиенту
	handleCardRequest(w, session)
}

func handleCardRequest(w http.ResponseWriter, session *Session) {
	session.Mutex.RLock() // Блокируем чтение для безопасного доступа к данным
	defer session.Mutex.RUnlock()

	data := TemplateData{
		Cards:   session.Cards,
		Message: fmt.Sprintf("%v", session.ErrorMsg),
	}

	// Отправляем данные в шаблон
	tmpl, err := template.ParseFiles("ui/base.html", "ui/cards.html")
	if err != nil {
		log.Println("Ошибка при загрузке шаблона:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)

	// Сбрасываем таймер неактивности
	session.Timer.Reset(30 * time.Minute)
}

// 2. Получили запрос на фиксацию результата
func HandlerCardsFix(w http.ResponseWriter, r *http.Request) {
	fmt.Println("cardsFixHandler")
	//session := GetOrCreateSession(w, r)

	var p interface{}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Создаем контекст с таймаутом в 2 секунды
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Канал для получения карточек и ошибки
	cardsChan := make(chan *cardlern.Cards, 1)
	errChan := make(chan error, 1)

	go func() {

		cards, err := initializeCards(map[string]interface{}{
			"limit": 5,
		})
		if err != nil {
			errChan <- err
			return
		}
		cardsChan <- cards
	}()

	// Ожидаем завершения инициализации или таймаута. Отправляем данные клиенту
	select {
	case _ = <-cardsChan:
		w.WriteHeader(http.StatusOK)
	case _ = <-errChan:
		w.WriteHeader(http.StatusBadRequest)
	case <-ctx.Done():
		w.WriteHeader(http.StatusGatewayTimeout)
	}

}

// .INT.

// загрузить карточки из базы
func initializeCards(query map[string]interface{}) (*cardlern.Cards, error) {

	// Загружаем данные из базы
	glossary, err := SRV.dbPool.ViewRandomGlossary(query)
	if err != nil {
		return nil, fmt.Errorf("Ошибка загрузки данных:", err)
	}

	// Создаем карточки
	cards, err := cardlern.New(glossary)
	if err != nil {
		return nil, fmt.Errorf("Ошибка создания cards:", err)
	}

	//test
	//fmt.Println("Начало работы...")
	//time.Sleep(5 * time.Second)
	//fmt.Println("Работа завершена.")

	return cards, nil
}
