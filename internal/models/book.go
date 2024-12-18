package models

type Book struct {
	Title       string   `json:"title"`
	Authors     []string `json:"authors"`
	Description string   `json:"description"`
	PageCount   int      `json:"page_count"`

	// Cover URLs is commented in indefinite time
	// CoverUrls   []string `json:"cover_urls"`
}

func NewBook(title string, authors []string, description string, pageCount int) Book {
	return Book{
		Title:       title,
		Authors:     authors,
		Description: description,
		PageCount:   pageCount,
	}
}

func NewDefaultBook() Book {
	return Book{
		Title:       "Преступление и наказание",
		Authors:     []string{"Фёдор Достоевский"},
		Description: "«Преступление и наказание», по мнению многих критиков, является лучшим романом Федора Достоевского, который оказал значительное влияние на русскую и мировую литературу.",
		PageCount:   1337,
	}
}
