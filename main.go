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
		handler := itemHandler.NewItemKindHandler(useCase)
		kind.GET("/:id", handler.Get)
		kind.GET("/", handler.GetAll)
		kind.POST("/", handler.Post)
		kind.PUT("/:id", handler.Put)
		kind.DELETE("/:id", handler.Delete)
	}
	kindRepo := memory.NewItemKindMemoryRepository()
	// stock
	stock := r.Group("/item/stock")
	{
		stockRepo := memory.NewStockItemMemoryRepository()
		useCase := itemUseCase.NewStockItemUseCase(stockRepo, kindRepo)
		handler := itemHandler.NewStockItemHandler(useCase)
		stock.GET("/:id", handler.Get)
		stock.GET("/", handler.GetAll)
		stock.POST("/", handler.Post)
		stock.PUT("/:id", handler.Put)
		stock.PUT("/:id/remain", handler.PutRemain)
		stock.DELETE("/:id", handler.Delete)
	}
	// food
	// todo idのGET紐付け
	food := r.Group("/item/food")
	{
		foodRepo := memory.NewFoodItemMemoryRepository()
		useCase := itemUseCase.NewFoodItemUseCase(foodRepo, kindRepo)
		handler := itemHandler.NewFoodItemHandler(useCase)
		food.GET("/:id", handler.Get)
		food.GET("/", handler.GetAll)
		food.POST("/", handler.Post)
		food.PUT("/:id", handler.Put)
		food.DELETE("/:id", handler.Delete)
	}

	// hour
	businessHourRepo := memory.NewBusinessHoursMemoryRepository()
	spBusinessHourRepo := memory.NewSpecialBusinessHourMemoryRepository()
	hour := r.Group("/store/hour")
	{
		useCase := storeUseCase.NewBusinessHoursUseCase(businessHourRepo, spBusinessHourRepo)
		handler := storeHandler.NewbusinessHoursHandler(useCase)
		hour.GET("/", handler.Get)
		hour.PUT("/", handler.Put)
	}
	specialHour := r.Group("/store/special_hour")
	{
		useCase := storeUseCase.NewSpecialBusinessHoursUseCase(businessHourRepo, spBusinessHourRepo)
		handler := storeHandler.NewSpecialBusinessHourHandler(useCase)
		specialHour.GET("/:id", handler.Get)
		specialHour.GET("/", handler.GetAll)
		specialHour.POST("/", handler.Post)
		specialHour.PUT("/:id", handler.Put)
		specialHour.DELETE("/:id", handler.Delete)
	}

	//holiday
	holidayRepo := memory.NewSpecialHolidayMemoryRepository()
	holiday := r.Group("/store/holiday")
	{
		useCase := storeUseCase.NewSpecialHolidayUseCase(holidayRepo)
		handler := storeHandler.NewSpecialHolidayHandler(useCase)
		holiday.GET("/:id", handler.Get)
		holiday.GET("/", handler.GetAll)
		holiday.POST("/", handler.Post)
		holiday.PUT("/:id", handler.Put)
		holiday.DELETE("/:id", handler.Delete)
	}

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return r
}
