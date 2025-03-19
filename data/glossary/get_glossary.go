package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode"

	"dictionary/pkg/storage/postgres"

	"github.com/PuerkitoBio/goquery"
)

//url := "https://verb.ru/"

func getFullInformation(href, glossaryWord, glossaryTranslation, glossaryMark string) (postgres.Glossary, error) {

	var glossary postgres.Glossary
	var err error

	fmt.Println(href)

	// Выполнение GET-запроса
	resp, err := http.Get(href)
	if err != nil {
		return glossary, fmt.Errorf("Ошибка при выполнении GET-запроса: %v", err)
	}
	defer resp.Body.Close()

	// Проверка на успешный ответ
	if resp.StatusCode != http.StatusOK {
		return glossary, fmt.Errorf("Ошибка: статус-код %d", resp.StatusCode)
	}

	// Парсинг HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return glossary, fmt.Errorf("Ошибка парсинга HTML: %v", err)
	}

	postContent := doc.Find("div.post_content")

	// Word
	glossary.Word, err = checkStringOnValidData(glossaryWord, "en")
	if err != nil {
		return glossary, fmt.Errorf("Ошибка при проверки %v: %v", glossary.Word, err)
	}

	// LetterName
	glossary.LetterName = string(glossary.Word[0])

	// Transcription
	postContent.Find("span[style='background:#fcff99;color:#000000']").Each(func(i int, s *goquery.Selection) {
		glossary.Transcription = strings.TrimSpace(s.Text())
	})
	glossary.Transcription, err = checkStringOnValidData(glossary.Transcription, "")
	if err != nil {
		return glossary, fmt.Errorf("Ошибка при проверки %v: %v", glossary.Transcription, err)
	}

	// Translation
	var translations []string
	postContent.Find("strong a").Each(func(i int, s *goquery.Selection) {
		translations = append(translations, s.Text())
	})
	glossary.Translation = glossaryTranslation
	if len(translations) > 1 {
		glossary.Translation = strings.Join(translations, ", ")
	}

	glossary.Translation, err = checkStringOnValidData(glossary.Translation, "ru")
	if err != nil {
		return glossary, fmt.Errorf("Ошибка при проверки %v: %v", glossary.Translation, err)
	}

	// Example
	ex1List := ""
	ex2List := ""
	postContent.Find(".su-row").Each(func(i int, s *goquery.Selection) {
		ex1 := s.Find(".su-column-size-1-2.first_bold").Text()
		ex1List += ex1

		ex2 := s.Find(".su-column-size-1-2:not(.first_bold)").Text()
		ex2List += ex2
		glossary.Example = append(glossary.Example, postgres.ExampleColumn{
			Ex1: strings.TrimSpace(ex1),
			Ex2: strings.TrimSpace(ex2),
		})
	})

	// Example
	if len(glossary.Example) == 0 {

		postContent.Find("div.su-divider.su-divider-style-default").Each(func(i int, s *goquery.Selection) {

			prev := s.PrevFiltered("p")
			if prev.Length() > 0 {

				strongText := prev.Find("strong").Text()
				ex1 := strongText
				ex1List += ex1

				russianText := prev.Text()
				russianText = strings.TrimPrefix(russianText, strongText) // Удаляем английский текст
				russianText = strings.TrimSpace(russianText)              // Убираем лишние пробелы
				ex2 := russianText
				ex2List += ex2

				glossary.Example = append(glossary.Example, postgres.ExampleColumn{
					Ex1: strings.TrimSpace(ex1),
					Ex2: strings.TrimSpace(ex2),
				})
			}

		})
	}
	ex1List, err = checkStringOnValidData(ex1List, "")
	if err != nil {
		return glossary, fmt.Errorf("Ошибка при проверки example en %v: %v", glossary.Example, err)
	}
	ex2List, err = checkStringOnValidData(ex2List, "")
	if err != nil {
		return glossary, fmt.Errorf("Ошибка при проверки example ru %v: %v", glossary.Example, err)
	}

	glossary.MarkName = glossaryMark
	glossary.Enable = true

	return glossary, nil
}

func checkStringOnValidData(s string, typeRegExp string) (string, error) {

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

func main() {
	// URL сайта
	url := "https://verb.ru/"

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

	var glossaryList []postgres.Glossary
	var glossaryListError []postgres.Glossary
	var wg sync.WaitGroup

	// Ищем все элементы <li> внутри вложенных <ul>
	stop := false // Флаг для досрочного выхода

	doc.Find("ul ul li").Each(func(i int, s *goquery.Selection) {

		if stop {
			return // Выходим из текущей итерации
		}

		link := s.Find("a")

		if link.Length() > 0 {

			glossaryWord := ""
			glossaryTranslation := ""
			glossaryMark := "Phrasal verbs"

			//1.
			word := link.Text()
			glossaryWord = word
			//2.
			liText := s.Text()
			translation := strings.Replace(liText, word, "", 1)
			translation = strings.Replace(translation, " – ", "", 1)
			translation = strings.Replace(translation, " – ", "", 1)

			glossaryTranslation = translation

			//X.
			href, _ := link.Attr("href")

			wg.Add(1)
			go func(href, glossaryWord, glossaryTranslation, glossaryMark string) {
				defer wg.Done()

				glossary, err := getFullInformation(href, glossaryWord, glossaryTranslation, glossaryMark)
				if err != nil {
					glossaryListError = append(glossaryListError, glossary)
					log.Printf("Ошибка обработки ссылки: %v %v\n", href, err)
				} else {
					glossaryList = append(glossaryList, glossary)
					//log.Printf("Ссылка добавлена: %v %v\n", href, err)
				}

			}(href, glossaryWord, glossaryTranslation, glossaryMark)

			time.Sleep(100 * time.Millisecond)

		}

		// Условие для досрочного выхода
		/*if i == 20 {
			stop = true // Устанавливаем флаг
			return
		}*/

	})

	wg.Wait()

	//file.json
	// Convert to JSON
	jsonData, err := json.MarshalIndent(glossaryList, "", "  ")
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
	jsonData, err = json.MarshalIndent(glossaryListError, "", "  ")
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
