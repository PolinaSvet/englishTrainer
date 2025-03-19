package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	// URL сайта
	url := "https://englex.ru/ways-to-say-yes-in-english/"
	// Наименование таблицы для фильтров

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

	// Ищем все <h2> внутри articleBody
	articleBody.Find("h2").Each(func(i int, h2 *goquery.Selection) {
		// Разделяем текст <h2> на английскую и русскую части
		h2Text := h2.Text()
		h2Parts := strings.Split(h2Text, " — ")
		if len(h2Parts) == 2 {
			englishH2 := re.ReplaceAllString(strings.TrimSpace(h2Parts[0]), "")
			russianH2 := strings.TrimSpace(h2Parts[1])
			fmt.Printf("H2: English: %s, Russian: %s\n", englishH2, russianH2)

			// Ищем связанные элементы с классом article-example
			h2.NextUntil("h2").Filter(".article-example").Each(func(i int, example *goquery.Selection) {
				exampleText := example.Text()
				exampleParts := strings.Split(exampleText, " — ")
				if len(exampleParts) == 2 {
					englishExample := strings.TrimSpace(exampleParts[0])
					russianExample := strings.TrimSpace(exampleParts[1])
					fmt.Printf("  Example: English: %s, Russian: %s\n", englishExample, russianExample)
				}
			})

			// Ищем связанные элементы с классом article-dialog
			h2.NextUntil("h2").Filter(".article-dialog").Each(func(i int, dialog *goquery.Selection) {
				// Извлекаем сообщения
				dialogEnCnt := 0

				dialog.Find(".message-text-wrapper").Each(func(i int, message *goquery.Selection) {
					messageText := message.Text()

					// Выводим сообщения в формате [A]: и [B]:
					if dialogEnCnt%2 == 0 {
						fmt.Printf("  Dialog Message: [1]: %s\n", messageText)
					} else {
						fmt.Printf("  Dialog Message: [2]: %s\n", messageText)
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
					} else {
						fmt.Printf("  Dialog Translation: [2]: %s\n", messageText)
					}
					dialogRuCnt++
				})
			})
		}
	})

}
