package cardlern

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

type ExampleColumn struct {
	Ex1 string `json:"ex1"`
	Ex2 string `json:"ex2"`
}

type Card struct {
	Id            int             `json:"id"`
	MarkId        int             `json:"mark_id"`
	MarkName      string          `json:"mark_name"`
	LetterId      int             `json:"letter_id"`
	LetterName    string          `json:"letter_name"`
	Word          string          `json:"word"`
	Transcription string          `json:"transcription"`
	Translation   string          `json:"translation"`
	Example       []ExampleColumn `json:"example"`
	DtAdd         int             `json:"dt_add"`
	DtAddTxt      string          `json:"dt_add_txt"`
	Enable        bool            `json:"enable"`
	Attempt       int             `json:"attempt"`
	Guess         bool            `json:"guess"`
	Answers       []string        `json:"answers"`
}

// Data store
type Cards struct {
	Data       []Card `json:"data"`
	ScoreAll   int    `json:"scoreAll"`
	ScoreGuess int    `json:"scoreGuess"`
	Finish     bool   `json:"finish"`
}

// The constructor accepts a connection string to the database
func New(jsonRequest []map[string]interface{}) (*Cards, error) {
	if jsonRequest == nil {
		return nil, fmt.Errorf("empty glossary")
	}

	data := convertMapToCard(jsonRequest)
	if data == nil {
		return nil, fmt.Errorf("unknown format glossary")
	}

	cards := Cards{
		Data:       data,
		ScoreAll:   len(data),
		ScoreGuess: 0,
		Finish:     false,
	}
	return &cards, nil
}

func GetEmptyCard() *Cards {
	cards := Cards{
		Data:       nil,
		ScoreAll:   0,
		ScoreGuess: 0,
		Finish:     false,
	}
	return &cards
}

func (cards *Cards) GetUnGuessCard() (Card, error) {
	var falseGuesses []Card
	var card Card

	for _, card := range cards.Data {
		if !card.Guess {
			falseGuesses = append(falseGuesses, card)
		}
	}

	if falseGuesses == nil {
		return card, fmt.Errorf("no unguess word in card")
	}

	randomIndex := rand.Intn(len(falseGuesses))

	return falseGuesses[randomIndex], nil
}

// .INNER

func fillAnswers(cards []Card) {
	rand.Seed(time.Now().UnixNano())

	for i := range cards {
		// Создаем временный массив для хранения всех Translation, кроме текущего
		otherTranslations := make([]string, 0, len(cards)-1)
		for j, card := range cards {
			if j != i {
				otherTranslations = append(otherTranslations, card.Translation)
			}
		}

		// Выбираем 3 случайных Translation из других элементов
		selectedTranslations := make([]string, 0, 3)
		for len(selectedTranslations) < 3 {
			index := rand.Intn(len(otherTranslations))
			selectedTranslations = append(selectedTranslations, otherTranslations[index])
			// Удаляем выбранный элемент, чтобы избежать повторений
			otherTranslations = append(otherTranslations[:index], otherTranslations[index+1:]...)
		}

		// Добавляем текущий Translation и перемешиваем
		answers := append(selectedTranslations, cards[i].Translation)
		rand.Shuffle(len(answers), func(k, l int) {
			answers[k], answers[l] = answers[l], answers[k]
		})

		// Записываем результат в поле Answers
		cards[i].Answers = answers
	}
}

func convertMapToCard(data []map[string]interface{}) []Card {
	var cards []Card

	for _, item := range data {
		card := Card{
			Id:            int(item["id"].(int)),
			MarkId:        int(item["mark_id"].(int)),
			MarkName:      item["mark_name"].(string),
			LetterId:      int(item["letter_id"].(int)),
			LetterName:    item["letter_name"].(string),
			Word:          item["word"].(string),
			Transcription: item["transcription"].(string),
			Translation:   item["translation"].(string),
			Example:       convertInterfaceToExampleColumn(item["example"]),
			DtAdd:         int(item["dt_add"].(int)),
			DtAddTxt:      item["dt_add_txt"].(string),
			Enable:        item["enable"].(bool),
			Attempt:       0,     // default
			Guess:         false, // default
		}
		cards = append(cards, card)
	}

	fillAnswers(cards)

	return cards
}

func convertInterfaceToExampleColumn(data interface{}) []ExampleColumn {
	var examples []ExampleColumn
	slice := reflect.ValueOf(data)

	for i := 0; i < slice.Len(); i++ {
		item := slice.Index(i).Interface().(map[string]interface{})
		example := ExampleColumn{
			Ex1: item["ex1"].(string),
			Ex2: item["ex2"].(string),
		}
		examples = append(examples, example)
	}

	return examples
}
