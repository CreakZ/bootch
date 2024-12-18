package routing

import (
	"bootch/internal/handlers"

	"github.com/gin-gonic/gin"
)

func InitRouting(rg *gin.RouterGroup) {
	rg.GET("/isbn10", handlers.GetBookByIsbn10)
	rg.GET("/isbn13", handlers.GetBookByIsbn13)
}
