package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
)

type ExampleColumn struct {
	Ex1 string `json:"ex1"`
	Ex2 string `json:"ex2"`
}

type Glossary struct {
	Letter        string          `json:"letter"`
	Word          string          `json:"word"`
	Transcription string          `json:"transcription"`
	Translation   string          `json:"translation"`
	Example       []ExampleColumn `json:"example"`
}

type GlossaryLoad struct {
	TableName    string
	TableComment string
	Data         []Glossary
}

func main() {

	// URL сайта
	//url := "https://englex.ru/ways-to-say-yes-in-english/"
	//tableName := "greeting"
	//tableComment := "Как сказать «да» в английском языке"

	url := "https://englex.ru/ways-to-say-i-love-you/"
	tableName := "iloveyou"
	tableComment := "15 способов сказать I love you"

	// Наименование таблицы для фильтров

	var glossaryList []Glossary
	var glossaryListError []Glossary
	errorCheck := false

	// Выполнение GET-запроса
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Ошибка при выполнении GET-запроса: %v", err)
	}
	defer resp.Body.Close() // Закрытие тела ответа после завершения функции

	// Проверка на успешный ответ
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Ошибка: статус-код %d", resp.StatusCode)
	}

	// Парсинг HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Ошибка парсинга HTML: %v", err)
	}

	// Находим элемент с атрибутом itemprop="articleBody"
	articleBody := doc.Find("[itemprop=articleBody]")

	// Регулярное выражение для удаления цифры с точкой в начале строки
	re := regexp.MustCompile(`^\d+\.\s*`)
	reEx := regexp.MustCompile(`"`)

	// Ищем все <h2> внутри articleBody
	articleBody.Find("h2").Each(func(i int, h2 *goquery.Selection) {
		var glossary Glossary

		// Разделяем текст <h2> на английскую и русскую части
		h2Text := h2.Text()
		h2Parts := strings.Split(h2Text, " — ")
		if len(h2Parts) == 2 {
			englishH2 := re.ReplaceAllString(strings.TrimSpace(h2Parts[0]), "")
			russianH2 := strings.TrimSpace(h2Parts[1])
			fmt.Printf("H2: English: %s, Russian: %s\n", englishH2, russianH2)

			// Word
			if glossary.Word, err = checkStringOnValidDataEnglix(englishH2, "en"); err != nil {
				errorCheck = true
			}

			// LetterName
			if !errorCheck {
				glossary.Letter = string(glossary.Word[0])
			}

			// Translation
			if glossary.Translation, err = checkStringOnValidDataEnglix(russianH2, "en"); err != nil {
				errorCheck = true
			}

			// Translation
			glossary.Transcription = ""

			// Ищем связанные элементы с классом article-example
			/*h2.NextUntil("h2").Filter(".article-example").Each(func(i int, example *goquery.Selection) {
				exampleText := example.Text()
				exampleParts := strings.Split(exampleText, " — ")
				if len(exampleParts) == 2 {
					englishExample := strings.TrimSpace(reEx.ReplaceAllString(exampleParts[0], "'"))
					russianExample := strings.TrimSpace(reEx.ReplaceAllString(exampleParts[1], "'"))
					fmt.Printf("  Example: English: %s, Russian: %s\n", englishExample, russianExample)

					// Example
					glossary.Example = append(glossary.Example, ExampleColumn{
						Ex1: strings.TrimSpace(englishExample),
						Ex2: strings.TrimSpace(russianExample),
					})
				}
			})*/

			// Ищем связанные элементы с классом article-example
			h2.NextUntil("h2").Find("span[style='font-size:18px']").Each(func(i int, example *goquery.Selection) {
				exampleText := example.Text()
				exampleParts := strings.Split(exampleText, " — ")
				if len(exampleParts) == 2 {
					englishExample := strings.TrimSpace(reEx.ReplaceAllString(exampleParts[0], "'"))
					russianExample := strings.TrimSpace(reEx.ReplaceAllString(exampleParts[1], "'"))
					fmt.Printf("  Example: English: %s, Russian: %s\n", englishExample, russianExample)

					// Example
					glossary.Example = append(glossary.Example, ExampleColumn{
						Ex1: strings.TrimSpace(englishExample),
						Ex2: strings.TrimSpace(russianExample),
					})
				}
			})

			// Ищем связанные элементы с классом article-dialog
			h2.NextUntil("h2").Filter(".article-dialog").Each(func(i int, dialog *goquery.Selection) {

				englishExample := ""
				russianExample := ""

				// Извлекаем сообщения
				dialogEnCnt := 0
				dialog.Find(".message-text-wrapper").Each(func(i int, message *goquery.Selection) {
					//messageText := message.Text()
					messageText := strings.TrimSpace(reEx.ReplaceAllString(message.Text(), "'"))

					// Выводим сообщения в формате [A]: и [B]:
					if dialogEnCnt%2 == 0 {
						fmt.Printf("  Dialog Message: [1]: %s\n", messageText)
						englishExample += fmt.Sprintf("[1]: %s\n", messageText)
					} else {
						fmt.Printf("  Dialog Message: [2]: %s\n", messageText)
						englishExample += fmt.Sprintf("[2]: %s\n", messageText)
					}
					dialogEnCnt++
				})

				// Извлекаем перевод
				dialogRuCnt := 0
				translation := dialog.Find(".translation-content")
				translation.Find("em").Each(func(i int, em *goquery.Selection) {
					// Убираем символы ( - ) и лишние пробелы
					messageText := strings.TrimSpace(em.Text())
					messageText = strings.TrimPrefix(messageText, "—")
					messageText = strings.TrimSpace(messageText)

					// Выводим сообщения в формате [A]: и [B]:
					if dialogRuCnt%2 == 0 {
						fmt.Printf("  Dialog Translation: [1]: %s\n", messageText)
						russianExample += fmt.Sprintf("[1]: %s\n", messageText)
					} else {
						fmt.Printf("  Dialog Translation: [2]: %s\n", messageText)
						russianExample += fmt.Sprintf("[2]: %s\n", messageText)
					}
					dialogRuCnt++
				})

				// Example
				glossary.Example = append(glossary.Example, ExampleColumn{
					Ex1: strings.TrimSpace(englishExample),
					Ex2: strings.TrimSpace(russianExample),
				})

			})
		}

		if errorCheck {
			glossaryListError = append(glossaryListError, glossary)
			log.Printf("Ошибка обработки ссылки: %v\n", glossary)
		} else {
			glossaryList = append(glossaryList, glossary)
		}
	})

	//file.json
	// Convert to JSON
	jsonData, err := json.MarshalIndent(GlossaryLoad{
		TableName:    tableName,
		TableComment: tableComment,
		Data:         glossaryList,
	}, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}

	// Write to file
	fileName := "glossary.json"
	err = os.WriteFile(fileName, jsonData, 0644) // 0644 is the permission
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	//fileError.json
	// Convert to JSON
	jsonData, err = json.MarshalIndent(GlossaryLoad{
		TableName:    tableName,
		TableComment: tableComment,
		Data:         glossaryListError,
	}, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}

	// Write to file
	fileName = "glossaryError.json"
	err = os.WriteFile(fileName, jsonData, 0644) // 0644 is the permission
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Printf("tWordsList count: %v, tWordsListError count: %v\n", len(glossaryList), len(glossaryListError))

}

func checkStringOnValidDataEnglix(s string, typeRegExp string) (string, error) {

	// Убираем лишние пробелы
	s = strings.TrimSpace(s)
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")

	if len(s) > 0 {
		r := []rune(s)               // конвертируем строку в slice рубик
		r[0] = unicode.ToUpper(r[0]) // переводим первый символ в верхний регистр
		s = string(r)                // конвертируем обратно в строку
	}

	if len(s) == 0 {
		return "", fmt.Errorf("Ошибка пустая строка")
	}

	//hasEnglishChar := regexp.MustCompile(`[a-zA-Z]`).MatchString(s)
	//hasEnglishWord := regexp.MustCompile(`\b[a-zA-Z]+\b`).MatchString(s)
	//hasRussianChar := regexp.MustCompile(`[А-Яа-яЁё]`).MatchString(s)
	//hasRussianWord := regexp.MustCompile(`\b[А-Яа-яЁё]+\b`).MatchString(s)

	switch typeRegExp {
	case "ru":
		re := regexp.MustCompile(`\b[a-zA-Z]+\b`)
		if re.MatchString(s) {
			return "", fmt.Errorf("Ошибка строка содержит EN слова")
		}
		return s, nil
	case "en":
		re := regexp.MustCompile(`\b[А-Яа-яЁё]+\b`)
		if re.MatchString(s) {
			return "", fmt.Errorf("Ошибка строка содержит RU слова")
		}
		return s, nil
	default:
		return s, nil
	}

}
