package main

import (
	"net/http"

	itemHandler "chico/takeout/handlers/item"
	storeHandler "chico/takeout/handlers/store"
	"chico/takeout/infrastructures/memory"
	itemUseCase "chico/takeout/usecase/item"
	storeUseCase "chico/takeout/usecase/store"

	"github.com/gin-gonic/gin"
)

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8089")
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	// r.LoadHTMLGlob("templates/**/*")
	// r.Static("/static", "./static")

	// itemkind
	kind := r.Group("/item/kind")
	{
		repo := memory.NewItemKindMemoryRepository()
		useCase := itemUseCase.NewItemKindUseCase(repo)
		handler := itemHandler.NewItemKindHandler(*useCase)
		kind.GET("/:id", handler.Get)
		kind.GET("/", handler.GetAll)
		kind.POST("/", handler.Post)
		kind.PUT("/:id", handler.Put)
		kind.DELETE("/:id", handler.Delete)
	}
	// stock
	stock := r.Group("/item/stock")
	{
		stockRepo := memory.NewStockItemMemoryRepository()
		kindRepo := memory.NewItemKindMemoryRepository()
		useCase := itemUseCase.NewStockItemUseCase(stockRepo, kindRepo)
		handler := itemHandler.NewStockItemHandler(*useCase)
		stock.GET("/:id", handler.Get)
		stock.GET("/", handler.GetAll)
		stock.POST("/", handler.Post)
		stock.PUT("/:id", handler.Put)
		stock.PUT("/:id/remain", handler.PutRemain)
		stock.DELETE("/:id", handler.Delete)
	}

	// hour
	hour := r.Group("/store/hour")
	{
		repo := memory.NewBusinessHoursMemoryRepository()
		useCase := storeUseCase.NewBusinessHoursUseCase(repo)
		handler := storeHandler.NewbusinessHoursHandler(*useCase)
		hour.GET("/", handler.Get)
		hour.PUT("/", handler.Put)
	}

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return r
}
