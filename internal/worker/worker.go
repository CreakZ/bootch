package worker

import (
	"bootch/internal/models"
	"bootch/internal/usecases"
	"bootch/pkg/utils"
	"bootch/pkg/validation"
	"fmt"
	"sync"
)

func GetBookWithIsbn(isbnStr string, isbnType int) ([]models.Book, error) {
	isbn := utils.CleanIsbn(isbnStr)
	switch isbnType {
	case 10:
		if valid := validation.IsIsbn10Valid(isbn); !valid {
			return []models.Book{}, fmt.Errorf("provided ISBN-10 is invalid")
		}
	case 13:
		if valid := validation.IsIsbn13Valid(isbn); !valid {
			return []models.Book{}, fmt.Errorf("provided ISBN-13 is invalid")
		}
	default:
		return []models.Book{}, fmt.Errorf("isbn type is neither 10 nor 13")
	}

	booksChan := make(chan models.Book, 3)
	errsChan := make(chan error, 3)

	wg := sync.WaitGroup{}

	// this is horrible but i am lazy :)

	wg.Add(1)

	// go func(wg *sync.WaitGroup) {
	// 	defer wg.Done()
	// 	book, err := usecases.GetBookViaChitaiGorod(isbn)
	// 	if err != nil {
	// 		errsChan <- err
	// 		return
	// 	}
	// 	booksChan <- book
	// }(&wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		book, err := usecases.GetBookViaGoogleBookApi(isbn)
		if err != nil {
			errsChan <- err
			return
		}
		booksChan <- book
	}(&wg)

	// go func(wg *sync.WaitGroup) {
	// 	defer wg.Done()
	// 	book, err := usecases.GetBookViaLivelib(isbn)
	// 	if err != nil {
	// 		errsChan <- err
	// 		return
	// 	}
	// 	booksChan <- book
	// }(&wg)

	wg.Wait()

	if len(booksChan) == 0 {
		return []models.Book{}, fmt.Errorf("no books with '%s' ISBN-%d", isbn, isbnType)
	}

	books := make([]models.Book, 3)
	for book := range booksChan {
		books = append(books, book)
	}
	return books, nil
}
