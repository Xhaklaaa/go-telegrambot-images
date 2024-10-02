package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Quote struct {
	Author string   `json:"author"`
	Text   string   `json:"quote"`
	Tags   []string `json:"tags"`
}

func ShouldRunParser(flagFile string) bool {
	_, err := os.Stat(flagFile)
	return os.IsNotExist(err)
}

func CreateFlagFile(flagFile string) {
	file, err := os.Create(flagFile)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
}

func ParseQuotes() error {
	indexURL := "https://www.azquotes.com/quotes/topics/index.html"

	doc, err := goquery.NewDocument(indexURL)
	if err != nil {
		return fmt.Errorf("failed to fetch index page: %w", err)
	}

	var allQuotes []Quote

	doc.Find("section.authors-page a").Each(func(i int, s *goquery.Selection) {
		topicURL, exists := s.Attr("href")
		if !exists {
			log.Println("Topic URL not found")
			return
		}

		if strings.Contains(topicURL, "/tags/") {
			log.Printf("Ignoring tag page: %s", topicURL)
			return
		}

		fullTopicURL := "https://www.azquotes.com" + topicURL
		log.Printf("Fetching topic page: %s", fullTopicURL)

		topicDoc, err := goquery.NewDocument(fullTopicURL)
		if err != nil {
			log.Printf("Failed to fetch topic page: %s", err)
			return
		}

		topicDoc.Find(".wrap-block").Each(func(j int, q *goquery.Selection) {
			quoteText := strings.TrimSpace(q.Find("p").Text())
			if quoteText == "" {
				log.Println("Quote text not found")
				return
			}

			// Получаем имя автора
			author := strings.TrimSpace(q.Find(".author a").Text())
			if author == "" {
				log.Println("Author not found")
				return
			}

			// Получаем теги
			var tags []string
			q.Find(".mytags a").Each(func(k int, t *goquery.Selection) {
				tags = append(tags, strings.TrimSpace(t.Text()))
			})

			allQuotes = append(allQuotes, Quote{
				Author: author,
				Text:   quoteText,
				Tags:   tags,
			})
		})
	})

	if len(allQuotes) == 0 {
		return fmt.Errorf("no quotes found")
	}

	file, err := os.Create("internal/image_creator/quotes.json")
	if err != nil {
		return fmt.Errorf("failed to create quotes.json: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(allQuotes); err != nil {
		return fmt.Errorf("failed to encode quotes to JSON: %w", err)
	}

	return nil
}

func ReadQuotes(filename string) ([]Quote, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var quotes []Quote
	err = json.Unmarshal(file, &quotes)
	if err != nil {
		return nil, err
	}

	return quotes, nil
}

func WriteQuotes(filename string, quotes []Quote) error {
	file, err := json.MarshalIndent(quotes, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, file, 0644)
}

func GetUniqueQuote(quotes []Quote, lastTags map[string]bool) (Quote, error) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(quotes), func(i, j int) { quotes[i], quotes[j] = quotes[j], quotes[i] })
	for _, quote := range quotes {
		unique := true
		for _, tag := range quote.Tags {
			if lastTags[tag] {
				unique = false
				break
			}
		}
		if unique {
			return quote, nil
		}
	}

	return Quote{}, fmt.Errorf("no unique quote found")
}

func RemoveQuote(quotes []Quote, quote Quote) []Quote {
	for i, q := range quotes {
		if q.Author == quote.Author && q.Text == quote.Text {
			return append(quotes[:i], quotes[i+1:]...)
		}
	}
	return quotes
}

func FormatTags(tags []string) string {
	var formattedTags []string
	for _, tag := range tags {
		formattedTags = append(formattedTags, "#"+tag)
	}
	return strings.Join(formattedTags, " ")
}
