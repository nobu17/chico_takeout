package tests

import (
	//"bytes"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	idomains "chico/takeout/domains/item"
	domains "chico/takeout/domains/order"
	sdomains "chico/takeout/domains/store"
	orderHandler "chico/takeout/handlers/order"
	"chico/takeout/infrastructures/memory"
	orderUseCase "chico/takeout/usecase/order"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	orderUrl = "/order"
)

var orderMemoryMaps map[string]*domains.OrderInfo

func SetupOrderInfoRouter() *gin.Engine {
	r := gin.Default()

	kindRepo := memory.NewItemKindMemoryRepository()
	kindRepo.Reset()
	kindMemoryMaps = kindRepo.GetMemory()
	kindIds := []string{}
	for kindId := range kindMemoryMaps {
		kindIds = append(kindIds, kindId)
	}
	businessHoursRepo := memory.NewBusinessHoursMemoryRepository()
	schedules := businessHoursRepo.GetMemory().GetSchedules()
	spBusinessHourRepo := memory.NewSpecialBusinessHourMemoryRepository()
	// create special lunch schedule
	specialSchedule, _ := sdomains.NewSpecialBusinessHour("特別日程3", "2055/05/08", "11:00", "14:00", schedules[1].GetId())
	spBusinessHourRepo.Create(specialSchedule)

	holidayRepo := memory.NewSpecialHolidayMemoryRepository()
	// create special holiday
	spHoliday1, _ := sdomains.NewSpecialHoliday("長期休暇", "2056/07/10", "2056/10/03")
	holidayRepo.Create(spHoliday1)

	schedule, _ := businessHoursRepo.Fetch()

	stockRepo := memory.NewStockItemMemoryRepository()
	// stockRepo.Reset()
	stockMemoryMaps = stockRepo.GetMemory()
	// add new stock item
	newStock1, _ := idomains.NewStockItem("stock3", "item3", 6, 4, 300, kindIds[0], true, "https://stock1.png")
	newStock1.SetRemain(99)
	stockRepo.Create(newStock1)
	newStock2, _ := idomains.NewStockItem("stock4", "item4", 7, 6, 400, kindIds[0], true, "https://stock2.jpg")
	newStock2.SetRemain(3)
	stockRepo.Create(newStock2)

	foodRepo := memory.NewFoodItemMemoryRepository()
	// foodRepo.Reset()
	foodMemoryMaps = foodRepo.GetMemory()
	// add new food item
	scheduleIds1 := []string{schedule.GetSchedules()[0].GetId(), schedule.GetSchedules()[1].GetId()}
	food1, _ := idomains.NewFoodItem("food3", "item3", 4, 10, 11, 222, kindIds[0], scheduleIds1, true, "https://food1.jpg")
	foodRepo.Create(food1)

	orderRepos := memory.NewOrderInfoMemoryRepository()
	orderRepos.Reset()
	orderMemoryMaps = orderRepos.GetMemory()
	order := r.Group(orderUrl)
	{
		useCase := orderUseCase.NewOrderInfoUseCase(orderRepos, stockRepo, foodRepo, businessHoursRepo, spBusinessHourRepo, holidayRepo)
		handler := orderHandler.NewOrderInfoHandler(useCase)
		order.GET("/:id", handler.Get)
		order.POST("/", handler.PostCreate)
		order.PUT("/:id", handler.PutCancel)
		order.GET("/user/:userId", handler.GetByUser)
		order.GET("/user/active/:userId", handler.GetActiveByUser)
	}
	return r
}

// func TestOrderInfoHandler_GET(t *testing.T) {
// 	r := SetupOrderInfoRouter()

// 	stockIds := []string{}
// 	for id := range stockMemoryMaps {
// 		stockIds = append(stockIds, id)
// 	}
// 	foodIds := []string{}
// 	for id := range foodMemoryMaps {
// 		foodIds = append(foodIds, id)
// 	}

// 	wants := []map[string]interface{}{
// 		{"userId": "user1", "memo": "memo1", "pickupDateTime": "2050/12/10 12:00",
// 			"stockItems": []map[string]interface{}{},
// 			"foodItems": []map[string]interface{}{
// 				{"itemId": foodIds[0], "name": "food1", "price": 100.0, "quantity": 3.0},
// 				{"itemId": foodIds[1], "name": "food2", "price": 200.0, "quantity": 1.0},
// 			},
// 		},
// 		{"userId": "user2", "memo": "memo2", "pickupDateTime": "2050/12/14 12:00",
// 			"stockItems": []map[string]interface{}{
// 				{"itemId": stockIds[0], "name": "stock1", "price": 100.0, "quantity": 2.0},
// 			},
// 			"foodItems": []map[string]interface{}{
// 				{"itemId": foodIds[0], "name": "food1", "price": 100.0, "quantity": 1.0},
// 			},
// 		},
// 	}
// 	index := 0
// 	for id := range orderMemoryMaps {

// 		req, _ := http.NewRequest("GET", orderUrl+"/"+id, nil)
// 		w := httptest.NewRecorder()
// 		r.ServeHTTP(w, req)

// 		assert.Equal(t, http.StatusOK, w.Code)

// 		fmt.Println("body", w.Body)
// 		var response map[string]interface{}
// 		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

// 		AssertMaps(t, response, wants[index])
// 		index++
// 	}
// }

func TestOrderInfoHandler_GETByUser(t *testing.T) {
	r := SetupOrderInfoRouter()

	stockIds := map[string]string{}
	for id, value := range stockMemoryMaps {
		stockIds[value.GetName()] = id
	}
	foodIds := map[string]string{}
	for id, value := range foodMemoryMaps {
		foodIds[value.GetName()] = id
	}
	userIds := []string{"user1"}

	wants := []map[string]interface{}{
		{"userId": "user1", "userName": "ユーザー1", "memo": "memo1", "pickupDateTime": "2050/12/10 12:00",
			"stockItems": []map[string]interface{}{},
			"foodItems": []map[string]interface{}{
				{"itemId": foodIds["food1"], "name": "food1", "price": 100.0, "quantity": 3.0},
				{"itemId": foodIds["food2"], "name": "food2", "price": 200.0, "quantity": 1.0},
			},
		},
		{"userId": "user2", "userName": "ユーザー2", "memo": "memo2", "pickupDateTime": "2050/12/14 12:00",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds["stock1"], "name": "stock1", "price": 100.0, "quantity": 2.0},
			},
			"foodItems": []map[string]interface{}{
				{"itemId": foodIds["food1"], "name": "food1", "price": 100.0, "quantity": 1.0},
			},
		},
	}
	index := 0
	for _, userId := range userIds {

		req, _ := http.NewRequest("GET", orderUrl+"/user/"+userId, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		fmt.Println("body", w.Body)
		var response []map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		fmt.Println("response", response)
		AssertMaps(t, response[0], wants[index])
		index++
	}
}

func TestOrderInfoHandler_POST_CREATE(t *testing.T) {
	r := SetupOrderInfoRouter()

	stockIds := map[string]string{}
	for id, value := range stockMemoryMaps {
		stockIds[value.GetName()] = id
	}
	foodIds := map[string]string{}
	for id, value := range foodMemoryMaps {
		foodIds[value.GetName()] = id
	}

	bodies := []map[string]interface{}{
		{"userId": "123", "Memo": "めも", "pickupDateTime": "2052/12/10 09:00", // tuesday morning
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds["stock1"], "quantity": 1},
			},
			"foodItems": []map[string]interface{}{},
		},
		{"userId": "1234", "Memo": "めも2", "pickupDateTime": "2052/12/10 11:30", // tuesday lunch start
			"userName":  "ユーザー1234",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds["stock1"], "quantity": 1},
			},
			"foodItems": []map[string]interface{}{},
		},
		{"userId": "12345", "Memo": "", "pickupDateTime": "2052/12/11 21:00", // wed dinner end and allow empty memo
			"userName":  "ユーザー12345",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds["stock1"], "quantity": 1},
			},
			"foodItems": []map[string]interface{}{},
		},
		{"userId": "123456", "Memo": "特別日程", "pickupDateTime": "2055/05/08 11:00", // special schedule lunch
			"userName":  "ユーザー123456",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds["stock1"], "quantity": 1},
			},
			"foodItems": []map[string]interface{}{},
		},
		{"userId": "123", "Memo": "めも", "pickupDateTime": "2052/12/10 09:00",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{},
			"foodItems": []map[string]interface{}{
				{"itemId": foodIds["food1"], "quantity": 1},
			}, // only food item
		},
		{"userId": "123", "Memo": "めも", "pickupDateTime": "2052/12/10 09:00",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds["stock1"], "quantity": 1},
			},
			"foodItems": []map[string]interface{}{
				{"itemId": foodIds["food1"], "quantity": 1},
			}, // both stock and food
		},
		{"userId": "123", "Memo": "めも", "pickupDateTime": "2052/12/10 09:00",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds["stock1"], "quantity": 1},
				{"itemId": stockIds["stock3"], "quantity": 3},
			},
			"foodItems": []map[string]interface{}{
				{"itemId": foodIds["food1"], "quantity": 1},
				{"itemId": foodIds["food2"], "quantity": 2},
			}, // both stock and food
		},
	}
	wants := []map[string]interface{}{
		{"userId": "123", "memo": "めも", "pickupDateTime": "2052/12/10 09:00",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds["stock1"], "name": "stock1", "price": 100.0, "quantity": 1.0},
			},
			"foodItems": []map[string]interface{}{},
		},
		{"userId": "1234", "memo": "めも2", "pickupDateTime": "2052/12/10 11:30",
			"userName":  "ユーザー1234",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds["stock1"], "name": "stock1", "price": 100.0, "quantity": 1.0},
			},
			"foodItems": []map[string]interface{}{},
		},
		{"userId": "12345", "memo": "", "pickupDateTime": "2052/12/11 21:00",
			"userName":  "ユーザー12345",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds["stock1"], "name": "stock1", "price": 100.0, "quantity": 1.0},
			},
			"foodItems": []map[string]interface{}{},
		},
		{"userId": "123456", "memo": "特別日程", "pickupDateTime": "2055/05/08 11:00",
			"userName":  "ユーザー123456",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds["stock1"], "name": "stock1", "price": 100.0, "quantity": 1.0},
			},
			"foodItems": []map[string]interface{}{},
		},
		{"userId": "123", "memo": "めも", "pickupDateTime": "2052/12/10 09:00",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{},
			"foodItems": []map[string]interface{}{
				{"itemId": foodIds["food1"], "name": "food1", "price": 100.0, "quantity": 1.0},
			},
		},
		{"userId": "123", "memo": "めも", "pickupDateTime": "2052/12/10 09:00",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds["stock1"], "name": "stock1", "price": 100.0, "quantity": 1.0},
			},
			"foodItems": []map[string]interface{}{
				{"itemId": foodIds["food1"], "name": "food1", "price": 100.0, "quantity": 1.0},
			},
		},
		{"userId": "123", "memo": "めも", "pickupDateTime": "2052/12/10 09:00",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds["stock1"], "name": "stock1", "price": 100.0, "quantity": 1.0},
				{"itemId": stockIds["stock3"], "name": "stock3", "price": 300.0, "quantity": 3.0},
			},
			"foodItems": []map[string]interface{}{
				{"itemId": foodIds["food1"], "name": "food1", "price": 100.0, "quantity": 1.0},
				{"itemId": foodIds["food2"], "name": "food2", "price": 200.0, "quantity": 2.0},
			},
		},
	}
	for index, body := range bodies {
		jBytes, err := json.Marshal(body)
		assert.NoError(t, err, "init json is failed")

		jsonStr := string(jBytes)
		fmt.Println(jsonStr)

		req, _ := http.NewRequest("POST", orderUrl+"/", bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		fmt.Println("body", w.Body)
		var idResponse map[string]string
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &idResponse)

		id := idResponse["id"]
		assert.NotEmpty(t, id, "response id should not be empty.")

		// get result to confirm
		req, _ = http.NewRequest("GET", orderUrl+"/"+id, nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		fmt.Println("body", w.Body)
		var response map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		AssertMaps(t, response, wants[index])
	}
}

type orderInfoErrorData struct {
	name string
	args map[string]interface{}
}

func getOrderInfoErrorData(stockIds, foodIds []string) []orderInfoErrorData {
	var items = []orderInfoErrorData{
		{name: "userId empty", args: map[string]interface{}{
			"userId": "", "memo": "めも", "pickupDateTime": "2052/12/10 09:00",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds[0], "quantity": 1.0},
			},
			"foodItems": []map[string]interface{}{},
		}},
		{name: "lack of userId", args: map[string]interface{}{
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"memo": "めも", "pickupDateTime": "2052/12/10 09:00",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds[0], "quantity": 1.0},
			},
			"foodItems": []map[string]interface{}{},
		}},
		{name: "userName empty", args: map[string]interface{}{
			"userId": "123", "memo": "めも", "pickupDateTime": "2052/12/10 09:00",
			"userName":  "",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds[0], "quantity": 1.0},
			},
			"foodItems": []map[string]interface{}{},
		}},
		{name: "userName over limit(10)", args: map[string]interface{}{
			"userId": "123", "memo": "めも", "pickupDateTime": "2052/12/10 09:00",
			"userName":  "12345678901",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds[0], "quantity": 1.0},
			},
			"foodItems": []map[string]interface{}{},
		}},
		{name: "lack of userName", args: map[string]interface{}{
			"userId": "123", "memo": "めも", "pickupDateTime": "2052/12/10 09:00",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds[0], "quantity": 1.0},
			},
			"foodItems": []map[string]interface{}{},
		}},
		{name: "memo is over limit length(500)", args: map[string]interface{}{
			"userId": "123", "memo": MakeRandomStr(501), "pickupDateTime": "2052/12/10 09:00",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds[0], "quantity": 1.0},
			},
			"foodItems": []map[string]interface{}{},
		}},
		{name: "pickup date is incorrect format(not have date time)", args: map[string]interface{}{
			"userId": "1234", "memo": "めも", "pickupDateTime": "2052/12/10",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds[0], "quantity": 1.0},
			},
			"foodItems": []map[string]interface{}{},
		}},
		{name: "pickup date is incorrect format", args: map[string]interface{}{
			"userId": "1234", "memo": "めも", "pickupDateTime": "abcd",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds[0], "quantity": 1.0},
			},
			"foodItems": []map[string]interface{}{},
		}},
		{name: "lack of pickup date", args: map[string]interface{}{
			"userId": "1234", "memo": "めも",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds[0], "quantity": 1.0},
			},
			"foodItems": []map[string]interface{}{},
		}},
		{name: "both stock items and food items are empty", args: map[string]interface{}{
			"userId": "123", "memo": "めも", "pickupDateTime": "2052/12/10 09:00",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{},
			"foodItems":  []map[string]interface{}{},
		}},
		{name: "stock item id is not exists", args: map[string]interface{}{
			"userId": "123", "memo": "めも", "pickupDateTime": "2052/12/10 09:00",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds[0], "quantity": 1.0},
				{"itemId": "1111", "quantity": 1.0},
			},
			"foodItems": []map[string]interface{}{
				{"itemId": foodIds[0], "quantity": 1.0},
			},
		}},
		{name: "food item id is not exists", args: map[string]interface{}{
			"userId": "123", "memo": "めも", "pickupDateTime": "2052/12/10 09:00",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds[0], "quantity": 1.0},
			},
			"foodItems": []map[string]interface{}{
				{"itemId": "123", "quantity": 1.0},
			},
		}},
		{name: "lack of fooditems", args: map[string]interface{}{
			"userId": "123", "memo": "めも", "pickupDateTime": "2052/12/10 09:00",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds[0], "quantity": 1.0},
			},
		}},
		{name: "lack of stockitems", args: map[string]interface{}{
			"userId": "123", "memo": "めも", "pickupDateTime": "2052/12/10 09:00",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"foodItems": []map[string]interface{}{
				{"itemId": foodIds[0], "quantity": 1.0},
			},
		}},
		{name: "stock items over max order", args: map[string]interface{}{
			"userId": "123", "memo": "めも", "pickupDateTime": "2052/12/10 09:00",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds[0], "quantity": 5.0},
			},
			"foodItems": []map[string]interface{}{},
		}},
		{name: "stock items lack of remain", args: map[string]interface{}{
			"userId": "123", "memo": "めも", "pickupDateTime": "2052/12/10 09:00",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"stockItems": []map[string]interface{}{
				{"itemId": stockIds[3], "quantity": 5.0},
			},
			"foodItems": []map[string]interface{}{},
		}},
		{name: "food items over max order", args: map[string]interface{}{
			"userId": "123", "memo": "めも", "pickupDateTime": "2052/12/10 09:00",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"foodItems": []map[string]interface{}{
				{"itemId": foodIds[0], "quantity": 5.0},
			},
			"stockItems": []map[string]interface{}{},
		}},
		{name: "ouf of bussiness hour", args: map[string]interface{}{
			"userId": "123", "memo": "めも", "pickupDateTime": "2052/12/10 05:00",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"foodItems": []map[string]interface{}{
				{"itemId": foodIds[0], "quantity": 1.0},
			},
			"stockItems": []map[string]interface{}{},
		}},
		{name: "ouf of bussiness day", args: map[string]interface{}{
			"userId": "123", "memo": "めも", "pickupDateTime": "2052/12/09 12:00",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"foodItems": []map[string]interface{}{
				{"itemId": foodIds[0], "quantity": 1.0},
			},
			"stockItems": []map[string]interface{}{},
		}},
		{name: "pickupdate is past", args: map[string]interface{}{
			"userId": "123", "memo": "めも", "pickupDateTime": "2020/12/09 12:00",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"foodItems": []map[string]interface{}{
				{"itemId": foodIds[0], "quantity": 1.0},
			},
			"stockItems": []map[string]interface{}{},
		}},
		{name: "pickupdate is out of sp hour", args: map[string]interface{}{
			"userId": "123", "memo": "めも", "pickupDateTime": "2055/05/08 14:10",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"foodItems": []map[string]interface{}{
				{"itemId": foodIds[0], "quantity": 1.0},
			},
			"stockItems": []map[string]interface{}{},
		}},
		{name: "pickupdate is in special holiday", args: map[string]interface{}{
			"userId": "123", "memo": "めも", "pickupDateTime": "2056/08/01 14:10",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"foodItems": []map[string]interface{}{
				{"itemId": foodIds[0], "quantity": 1.0},
			},
			"stockItems": []map[string]interface{}{},
		}},
	}
	return items
}

func TestOrderInfoHandler_POST_BadRequest(t *testing.T) {
	r := SetupOrderInfoRouter()

	stockIds := []string{}
	for id := range stockMemoryMaps {
		stockIds = append(stockIds, id)
	}
	foodIds := []string{}
	for id := range foodMemoryMaps {
		foodIds = append(foodIds, id)
	}

	inputs := getOrderInfoErrorData(stockIds, foodIds)
	assert.NotEqual(t, 0, len(inputs), "input data is empty")

	for _, tt := range inputs {
		fmt.Println("case:", tt.name)
		jBytes, err := json.Marshal(tt.args)
		assert.NoError(t, err, "init json is failed")

		req, _ := http.NewRequest("POST", orderUrl+"/", bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		fmt.Println("body", w.Body)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	}
}

func TestOrderInfoHandler_POST_BadRequest_FoodOrderLimits(t *testing.T) {
	r := SetupOrderInfoRouter()

	foodIds := []string{}
	for id := range foodMemoryMaps {
		foodIds = append(foodIds, id)
	}

	parepare := orderInfoErrorData{
		name: "consuming all food stock",
		args: map[string]interface{}{
			"userId": "123", "memo": "めも", "pickupDateTime": "2052/12/10 09:00",
			"userName":  "ユーザー123",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"foodItems": []map[string]interface{}{
				{"itemId": foodIds[2], "quantity": 10.0},
			},
			"stockItems": []map[string]interface{}{},
		},
	}

	input := orderInfoErrorData{
		name: "over limit foods order",
		args: map[string]interface{}{
			"userId": "1234", "memo": "めも", "pickupDateTime": "2052/12/10 09:00",
			"userName":  "ユーザー1234",
			"userEmail": "userx@hoge.com", "userTelNo": "123456789",
			"foodItems": []map[string]interface{}{
				{"itemId": foodIds[2], "quantity": 2.0},
			},
			"stockItems": []map[string]interface{}{},
		},
	}

	// at first consuming all stocks
	fmt.Println("case:", parepare.name)
	jBytes, err := json.Marshal(parepare.args)
	assert.NoError(t, err, "init json is failed")

	req, _ := http.NewRequest("POST", orderUrl+"/", bytes.NewBuffer(jBytes))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	fmt.Println("body", w.Body)
	assert.Equal(t, http.StatusOK, w.Code)

	// try to order limit
	fmt.Println("case:", input.name)
	jBytes, err = json.Marshal(input.args)
	assert.NoError(t, err, "init json is failed")

	req, _ = http.NewRequest("POST", orderUrl+"/", bytes.NewBuffer(jBytes))
	req.Header.Add("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	fmt.Println("body", w.Body)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
