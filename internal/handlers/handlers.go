package handlers

import (
	"bootch/internal/worker"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	Isbn10Type = 10
	Isbn13Type = 13
)

func GetBookByIsbn10(c *gin.Context) {
	isbnStr, exists := c.GetQuery("isbn")
	if !exists {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "'isbn' param required"})
		return
	}

	books, err := worker.GetBookWithIsbn(isbnStr, Isbn10Type)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
}

func GetBookByIsbn13(c *gin.Context) {
	isbnStr, exists := c.GetQuery("isbn")
	if !exists {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "'isbn' param required"})
		return
	}

	books, err := worker.GetBookWithIsbn(isbnStr, Isbn13Type)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
}
