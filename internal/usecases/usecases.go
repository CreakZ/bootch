package usecases

import (
	"bootch/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	googleAPI   = "GoogleAPI"
	chitaiGorod = "Читай Город"
	livelib     = "Livelib"

	GoogleAPIUrl   = "https://www.googleapis.com/books/v1/volumes?q=isbn:"
	ChitaiGorodUrl = "https://www.chitai-gorod.ru/search?phrase="
	LivelibUrl     = "https://www.livelib.ru/find/books/"
)

func GetBookViaGoogleBookApi(isbn string) (models.Book, error) {
	URL := GoogleAPIUrl + isbn

	resp, err := http.Get(URL)
	if err != nil {
		return models.Book{}, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	info := struct {
		Items      int           `json:"totalItems"`
		VolumeInfo []interface{} `json:"items"`
	}{
		Items:      0,
		VolumeInfo: nil,
	}

	if err = json.Unmarshal(body, &info); err != nil {
		return models.Book{}, err
	}

	if info.Items == 0 {
		return models.Book{}, newErrorWithCaller(googleAPI, fmt.Sprintf("no books with %s isbn found", isbn))
	}

	volumeInfo, ok := info.VolumeInfo[0].(map[string]interface{})["volumeInfo"].(map[string]interface{})
	if !ok {
		return models.Book{}, newErrorWithCaller(googleAPI, "response json has unexpected structure")
	}

	var book models.Book

	book.Title, ok = volumeInfo["title"].(string)
	if !ok {
		book.Title = ""
	}

	temp, exists := volumeInfo["authors"].([]interface{})
	if !exists {
		book.Authors = []string{""}
	}

	for _, author := range temp {
		book.Authors = append(book.Authors, author.(string))
	}

	var desc string
	desc, ok = volumeInfo["description"].(string)
	if !ok {
		fmt.Println("no description provided")
	} else {
		book.Description = desc
	}

	return book, nil
}

func GetBookViaChitaiGorod(isbn string) (models.Book, error) {
	url := ChitaiGorodUrl + isbn

	resp, err := http.Get(url)
	if err != nil {
		return models.Book{}, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return models.Book{}, err
	}

	val, _ := doc.Find("h4").Attr("class")
	if val == "catalog-empty-result__header" {
		return models.Book{}, newErrorWithCaller(chitaiGorod, fmt.Sprintf("no books with '%s' isbn found", isbn))
	}

	var (
		link string
		book models.Book
	)

	doc.Find("article").
		EachWithBreak(func(i int, s *goquery.Selection) bool {
			value, exists := s.Find("a").Attr("href")
			if exists {
				link = value
				return false
			}
			return true
		})

	bookUrl := "https://www.chitai-gorod.ru" + link

	resp, err = http.Get(bookUrl)
	if err != nil {
		return models.Book{}, err
	}
	defer resp.Body.Close()

	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return models.Book{}, err
	}

	doc.Find("h1").
		EachWithBreak(func(i int, s *goquery.Selection) bool {
			class, _ := s.Attr("class")
			if class == "detail-product__header-title" {
				book.Title = strings.TrimSpace(s.Text())
				return false
			}
			return true
		})

	doc.Find("a").
		EachWithBreak(func(i int, s *goquery.Selection) bool {
			itemprop, _ := s.Attr("class")
			if itemprop == "product-info-authors__author" {
				author, _ := strings.CutSuffix(strings.TrimSpace(s.Text()), ",")
				book.Authors = append(book.Authors, author)

				return false
			}
			return true
		})

	if desc := doc.Find("article"); desc != nil {
		text := strings.TrimSpace(desc.Text())
		book.Description = text
	}

	return book, nil
}

func GetBookViaLivelib(isbn string) (models.Book, error) {
	url := LivelibUrl + isbn

	resp, err := http.Get(url)
	if err != nil {
		return models.Book{}, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return models.Book{}, err
	}

	var book models.Book

	doc.Find("a").
		Each(func(i int, s *goquery.Selection) {
			class, _ := s.Attr("class")

			switch class {
			case "title":
				book.Title = s.Text()

			case "description":
				book.Authors = append(book.Authors, s.Text())
			}
		})

	doc.Find("span").
		EachWithBreak(func(i int, s *goquery.Selection) bool {
			class, _ := s.Attr("class")
			if class == "description" {
				text := s.Text()
				book.Description = text
				return false
			}
			return true
		})

	return book, nil
}

func newErrorWithCaller(caller, msg string) error {
	return fmt.Errorf("%s: %s", caller, msg)
}
