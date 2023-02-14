package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	domains "chico/takeout/domains/item"
	itemHandler "chico/takeout/handlers/item"
	"chico/takeout/infrastructures/memory"
	itemUseCase "chico/takeout/usecase/item"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var foodMemoryMaps map[string]*domains.FoodItem

func SetupFoodItemRouter() *gin.Engine {
	r := gin.Default()
	kindRepo := memory.NewItemKindMemoryRepository()
	kindMemoryMaps = kindRepo.GetMemory()
	businessHourRepo := memory.NewBusinessHoursMemoryRepository()
	businessHoursMemory = businessHourRepo.GetMemory()
	food := r.Group("/item/food")
	{
		foodRepo := memory.NewFoodItemMemoryRepository()
		foodRepo.Reset()
		foodMemoryMaps = foodRepo.GetMemory()

		useCase := itemUseCase.NewFoodItemUseCase(foodRepo, kindRepo, businessHourRepo)
		handler := itemHandler.NewFoodItemHandler(useCase)
		food.GET("/:id", handler.Get)
		food.GET("/", handler.GetAll)
		food.POST("/", handler.Post)
		food.PUT("/:id", handler.Put)
		food.DELETE("/:id", handler.Delete)
	}
	return r
}

func TestFoodItemHandler_GET_ALL(t *testing.T) {

	kinds := []map[string]interface{}{
		{"priority": 1, "name": "item1"},
		{"priority": 2, "name": "item2"},
	}
	wants := []map[string]interface{}{
		{"priority": 1, "name": "food1", "description": "item1", "maxOrder": 4, "price": 100, "enabled": true, "kind": kinds[0], "maxOrderPerDay": 10, "imageUrl": "https://food1.jpg", "allowDates": []string{}},
		{"priority": 2, "name": "food2", "description": "item2", "maxOrder": 5, "price": 200, "enabled": true, "kind": kinds[1], "maxOrderPerDay": 18, "imageUrl": "", "allowDates": []string{}},
		{"priority": 3, "name": "food3", "description": "item3", "maxOrder": 6, "price": 110, "enabled": true, "kind": kinds[1], "maxOrderPerDay": 20, "imageUrl": "", "allowDates": []string{"2023/12/10", "2023/12/13"}},
	}
	wantsScheduleIdLength := []int{
		2, 2, 2,
	}

	r := SetupFoodItemRouter()
	req, _ := http.NewRequest("GET", "/item/food/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	fmt.Println("body", w.Body)

	// Convert the JSON response to a map
	var response []map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.Bytes()), &response)
	if err != nil {
		assert.Fail(t, "failed to Unmarshal json", err)
		return
	}

	assert.Equal(t, 3, len(response), "response length is incorrect")

	for index, item := range response {
		scheduleIds, ok := item["scheduleIds"].([]interface{})
		if !ok {
			assert.Fail(t, "failed to get id")
			continue
		}
		assert.Equal(t, wantsScheduleIdLength[index], len(scheduleIds))
		AssertMaps(t, item, wants[index])
	}
}

func TestFoodItemHandler_GET(t *testing.T) {
	r := SetupFoodItemRouter()

	kinds := []map[string]interface{}{
		{"priority": 1, "name": "item1"},
		{"priority": 2, "name": "item2"},
	}
	wants := []map[string]interface{}{
		{"priority": 1, "name": "food1", "description": "item1", "maxOrder": 4, "price": 100, "enabled": true, "kind": kinds[0], "maxOrderPerDay": 10, "imageUrl": "https://food1.jpg", "allowDates": []string{}},
		{"priority": 2, "name": "food2", "description": "item2", "maxOrder": 5, "price": 200, "enabled": true, "kind": kinds[1], "maxOrderPerDay": 18, "imageUrl": "", "allowDates": []string{}},
		{"priority": 3, "name": "food3", "description": "item3", "maxOrder": 6, "price": 110, "enabled": true, "kind": kinds[1], "maxOrderPerDay": 20, "imageUrl": "", "allowDates": []string{"2023/12/10", "2023/12/13"}},
	}
	wantsScheduleIdLength := []int{
		2, 2, 2,
	}

	index := 0
	for id := range foodMemoryMaps {
		req, _ := http.NewRequest("GET", "/item/food/"+id, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		fmt.Println("response", response)
		scheduleIds, ok := response["scheduleIds"].([]interface{})
		if !ok {
			assert.Fail(t, "failed to get id")
			continue
		}
		assert.Equal(t, wantsScheduleIdLength[index], len(scheduleIds))
		AssertMaps(t, response, wants[index])
		index++
	}
}

func TestFoodItemHandler_GET_NotFound(t *testing.T) {
	r := SetupFoodItemRouter()

	req, _ := http.NewRequest("GET", "/item/food/1234", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

	assert.Nil(t, response)
}

func TestFoodItemHandler_POST_CREATE(t *testing.T) {
	r := SetupFoodItemRouter()

	scheduleIds := []string{}
	for _, sch := range businessHoursMemory.GetSchedules() {
		scheduleIds = append(scheduleIds, sch.GetId())
	}
	kindIds := []string{}
	for id := range kindMemoryMaps {
		kindIds = append(kindIds, id)
	}

	kinds := []map[string]interface{}{
		{"priority": 1, "name": "item1"},
		{"priority": 2, "name": "item2"},
	}

	wants := []map[string]interface{}{
		{"priority": 3, "name": "create_food1", "description": "desc1", "maxOrder": 1, "price": 1000, "enabled": false, "kind": kinds[1], "maxOrderPerDay": 11, "scheduleIds": []string{scheduleIds[2]}, "imageUrl": "https://hoge.png", "allowDates": []string{}},
		{"priority": 5, "name": "create_food2", "description": "desc2", "maxOrder": 2, "price": 2000, "enabled": true, "kind": kinds[0], "maxOrderPerDay": 14, "scheduleIds": []string{scheduleIds[0], scheduleIds[1]}, "imageUrl": "", "allowDates": []string{}},
		{"priority": 6, "name": "create_food3", "description": "desc3", "maxOrder": 2, "price": 2000, "enabled": true, "kind": kinds[0], "maxOrderPerDay": 14, "scheduleIds": []string{scheduleIds[0], scheduleIds[1]}, "imageUrl": "", "allowDates": []string{"2023/12/10", "2023/12/20"}},
		{"priority": 7, "name": "free_food", "description": "free", "maxOrder": 2, "price": 0, "enabled": true, "kind": kinds[0], "maxOrderPerDay": 12, "scheduleIds": []string{scheduleIds[0], scheduleIds[1]}, "imageUrl": "", "allowDates": []string{}},
	}
	bodies := []map[string]interface{}{
		{"priority": 3, "name": "create_food1", "description": "desc1", "maxOrder": 1, "price": 1000, "enabled": false, "kindId": kindIds[1], "maxOrderPerDay": 11, "scheduleIds": []string{scheduleIds[2]}, "imageUrl": "https://hoge.png", "allowDates": []string{}},
		{"priority": 5, "name": "create_food2", "description": "desc2", "maxOrder": 2, "price": 2000, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 14, "scheduleIds": []string{scheduleIds[0], scheduleIds[1]}, "imageUrl": "", "allowDates": []string{}},
		{"priority": 6, "name": "create_food3", "description": "desc3", "maxOrder": 2, "price": 2000, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 14, "scheduleIds": []string{scheduleIds[0], scheduleIds[1]}, "imageUrl": "", "allowDates": []string{"2023/12/10", "2023/12/20"}},
		{"priority": 7, "name": "free_food", "description": "free", "maxOrder": 2, "price": 0, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 12, "scheduleIds": []string{scheduleIds[0], scheduleIds[1]}, "imageUrl": "", "allowDates": []string{}},
	}
	for index, body := range bodies {
		jBytes, err := json.Marshal(body)
		assert.NoError(t, err, "init json is failed")

		req, _ := http.NewRequest("POST", "/item/food/", bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		fmt.Println("body", w.Body)
		assert.Equal(t, http.StatusOK, w.Code)

		var idResponse map[string]string
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &idResponse)

		id := idResponse["id"]
		assert.NotEmpty(t, id, "response id should not be empty.")

		// confirm result from result id
		getReq, _ := http.NewRequest("GET", "/item/food/"+id, nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, getReq)

		assert.Equal(t, http.StatusOK, w.Code)
		fmt.Println("response", w.Body)

		var response map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		AssertMaps(t, response, wants[index])
	}
}

type foodItemErrorData struct {
	name string
	args map[string]interface{}
	want int
}

func GetFoodItemErrorData(kindIds, scheduleIds []string) []foodItemErrorData {
	var foodItemErrorInputs = []foodItemErrorData{
		{name: "lack priority", args: map[string]interface{}{
			"name": "stock1", "description": "item1",
			"maxOrder": 1, "price": 100, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 5, "scheduleIds": []string{scheduleIds[0]}, "imageUrl": "https://hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "error priority(0)", args: map[string]interface{}{
			"priority": 0, "name": "stock1", "description": "item1",
			"maxOrder": 1, "price": 100, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 5, "scheduleIds": []string{scheduleIds[0]}, "imageUrl": "https://hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "error priority(-1)", args: map[string]interface{}{
			"priority": -1, "name": "stock1", "description": "item1",
			"maxOrder": 1, "price": 100, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 5, "scheduleIds": []string{scheduleIds[0]}, "imageUrl": "https://hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "lack name", args: map[string]interface{}{
			"priority": 1, "description": "item1",
			"maxOrder": 1, "price": 100, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 5, "scheduleIds": []string{scheduleIds[0]}, "imageUrl": "https://hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "error name(empty)", args: map[string]interface{}{
			"priority": 1, "name": "", "description": "item1",
			"maxOrder": 1, "price": 100, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 5, "scheduleIds": []string{scheduleIds[0]}, "imageUrl": "https://hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "error name(over limit(26))", args: map[string]interface{}{
			"priority": 1, "name": MakeRandomStr(26), "description": "item1",
			"maxOrder": 1, "price": 100, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 1, "scheduleIds": []string{scheduleIds[0]}, "imageUrl": "https://hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "lack description", args: map[string]interface{}{
			"priority": 1, "name": "stock1",
			"maxOrder": 1, "price": 100, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 1, "scheduleIds": []string{scheduleIds[0]}, "imageUrl": "https://hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "error description(empty)", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "",
			"maxOrder": 1, "price": 100, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 1, "scheduleIds": []string{scheduleIds[0]}, "imageUrl": "https://hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "error description(over limit(151))", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": MakeRandomStr(151),
			"maxOrder": 1, "price": 100, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 1, "scheduleIds": []string{scheduleIds[0]}, "imageUrl": "https://hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "lack maxOrder", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "item1",
			"price": 100, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 1, "scheduleIds": []string{scheduleIds[0]}, "imageUrl": "https://hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "error maxOrder(0)", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "desc",
			"maxOrder": 0, "price": 100, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 1, "scheduleIds": []string{scheduleIds[0]}, "imageUrl": "https://hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "error maxOrder(31)", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "desc",
			"maxOrder": 31, "price": 100, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 1, "scheduleIds": []string{scheduleIds[0]}, "imageUrl": "https://hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "lack price", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "item1",
			"maxOrder": 1, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 1, "scheduleIds": []string{scheduleIds[0]}, "imageUrl": "https://hoge.png",
		}, want: 3},
		{name: "error price(-1)", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "desc",
			"maxOrder": 1, "price": -1, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 1, "scheduleIds": []string{scheduleIds[0]}, "imageUrl": "https://hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "error price(20001)", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "desc",
			"maxOrder": 1, "price": 20001, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 1, "scheduleIds": []string{scheduleIds[0]}, "imageUrl": "https://hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "lack enabled", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "item1",
			"maxOrder": 1, "price": 100, "kindId": kindIds[0], "maxOrderPerDay": 1, "scheduleIds": []string{scheduleIds[0]}, "imageUrl": "https://hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "lack kindId", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "item1",
			"maxOrder": 1, "price": 100, "enabled": true, "maxOrderPerDay": 1, "scheduleIds": []string{scheduleIds[0]}, "imageUrl": "https://hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "error kindId(empty)", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "desc",
			"maxOrder": 1, "price": 100, "enabled": true, "kindId": "", "maxOrderPerDay": 1, "scheduleIds": []string{scheduleIds[0]}, "imageUrl": "https://hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "not exist kindId", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "item1",
			"maxOrder": 1, "price": 100, "enabled": true, "kindId": "abc", "maxOrderPerDay": 1, "scheduleIds": []string{scheduleIds[0]}, "imageUrl": "https://hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "lack maxOrderPerDay", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "item1", "maxOrder": 1,
			"price": 100, "enabled": true, "kindId": kindIds[0], "scheduleIds": []string{scheduleIds[0]}, "imageUrl": "https://hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "error maxOrderPerDay(101)", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "item1", "maxOrder": 1,
			"price": 100, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 101, "scheduleIds": []string{scheduleIds[0]}, "imageUrl": "https://hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "error maxOrderPerDay(0)", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "item1", "maxOrder": 1,
			"price": 100, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": -1, "scheduleIds": []string{scheduleIds[0]}, "imageUrl": "https://hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "lack scheduleIds", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "item1", "maxOrder": 1,
			"price": 100, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 101, "imageUrl": "https://hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "error scheduleIds(empty)", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "item1", "maxOrder": 1,
			"price": 100, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 99, "scheduleIds": []string{}, "imageUrl": "https://hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "error not exists scheduleIds", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "item1", "maxOrder": 2,
			"price": 100, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 10, "scheduleIds": []string{"1234"}, "imageUrl": "https://hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "error not exists scheduleIds and exist id", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "item1", "maxOrder": 2,
			"price": 100, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 10, "scheduleIds": []string{scheduleIds[1], "1234"}, "imageUrl": "https://hoge.png",
		}, want: 3},
		{name: "error maxOrder > maxOrderPerDay", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "item1", "maxOrder": 2,
			"price": 100, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 1, "scheduleIds": []string{scheduleIds[1]}, "imageUrl": "https://hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "error lack of imageUrl", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "item1", "maxOrder": 1,
			"price": 100, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 2, "scheduleIds": []string{scheduleIds[0]},
			"allowDates": []string{},
		}, want: 3},
		{name: "error incorrect imageUrl format", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "item1", "maxOrder": 1,
			"price": 100, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 2, "scheduleIds": []string{scheduleIds[0]}, "imageUrl": "hoge.png",
			"allowDates": []string{},
		}, want: 3},
		{name: "error incorrect date format", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "item1", "maxOrder": 1,
			"price": 100, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 2, "scheduleIds": []string{scheduleIds[0]}, "imageUrl": "",
			"allowDates": []string{"2023"},
		}, want: 3},
		{name: "error incorrect date format mixed", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "item1", "maxOrder": 1,
			"price": 100, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 2, "scheduleIds": []string{scheduleIds[0]}, "imageUrl": "",
			"allowDates": []string{"2023/12/10", "abcd"},
		}, want: 3},
	}
	return foodItemErrorInputs
}

func TestFoodItemHandler_POST_CREATE_BadRequest(t *testing.T) {
	r := SetupFoodItemRouter()

	scheduleIds := []string{}
	for _, sch := range businessHoursMemory.GetSchedules() {
		scheduleIds = append(scheduleIds, sch.GetId())
	}

	kindIds := []string{}
	for id := range kindMemoryMaps {
		kindIds = append(kindIds, id)
	}

	inputs := GetFoodItemErrorData(kindIds, scheduleIds)
	assert.NotEqual(t, 0, len(inputs), "input data is empty")

	for _, tt := range inputs {
		fmt.Println("case:", tt.name)
		jBytes, err := json.Marshal(tt.args)
		assert.NoError(t, err, "init json is failed")

		req, _ := http.NewRequest("POST", "/item/food/", bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		fmt.Println(w.Body)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// confirm stock is not added
		req, _ = http.NewRequest("GET", "/item/food/", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "failed to get result to confirm.")

		var response []map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		assert.Equal(t, tt.want, len(foodMemoryMaps), "response length is incorrect")
	}
}

func TestFoodItemHandler_PUT(t *testing.T) {
	r := SetupFoodItemRouter()

	scheduleIds := []string{}
	for _, sch := range businessHoursMemory.GetSchedules() {
		scheduleIds = append(scheduleIds, sch.GetId())
	}

	kindIds := []string{}
	for id := range kindMemoryMaps {
		kindIds = append(kindIds, id)
	}
	foodIds := []string{}
	for id := range foodMemoryMaps {
		foodIds = append(foodIds, id)
	}

	kinds := []map[string]interface{}{
		{"priority": 1, "name": "item1"},
		{"priority": 2, "name": "item2"},
	}

	wants := []map[string]interface{}{
		{"priority": 3, "name": "create_food1", "description": "desc1", "maxOrder": 1, "price": 1000, "enabled": false, "kind": kinds[1], "maxOrderPerDay": 11, "scheduleIds": []string{scheduleIds[2]}, "imageUrl": "https://hoge.png", "allowDates": []string{}},
		{"priority": 5, "name": "create_food2", "description": "desc2", "maxOrder": 2, "price": 2000, "enabled": true, "kind": kinds[0], "maxOrderPerDay": 14, "scheduleIds": []string{scheduleIds[0], scheduleIds[1]}, "imageUrl": "", "allowDates": []string{}},
		{"priority": 5, "name": "create_food3", "description": "desc3", "maxOrder": 3, "price": 2200, "enabled": true, "kind": kinds[0], "maxOrderPerDay": 14, "scheduleIds": []string{scheduleIds[0], scheduleIds[1]}, "imageUrl": "", "allowDates": []string{"2020/12/10", "2023/12/10"}},
		{"priority": 5, "name": "free_food", "description": "free", "maxOrder": 2, "price": 0, "enabled": true, "kind": kinds[0], "maxOrderPerDay": 14, "scheduleIds": []string{scheduleIds[0], scheduleIds[1]}, "imageUrl": "", "allowDates": []string{}},
	}
	bodies := []map[string]interface{}{
		{"priority": 3, "name": "create_food1", "description": "desc1", "maxOrder": 1, "price": 1000, "enabled": false, "kindId": kindIds[1], "maxOrderPerDay": 11, "scheduleIds": []string{scheduleIds[2]}, "imageUrl": "https://hoge.png", "allowDates": []string{}},
		{"priority": 5, "name": "create_food2", "description": "desc2", "maxOrder": 2, "price": 2000, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 14, "scheduleIds": []string{scheduleIds[0], scheduleIds[1]}, "imageUrl": "", "allowDates": []string{}},
		{"priority": 5, "name": "create_food3", "description": "desc3", "maxOrder": 3, "price": 2200, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 14, "scheduleIds": []string{scheduleIds[0], scheduleIds[1]}, "imageUrl": "", "allowDates": []string{"2020/12/10", "2023/12/10"}},
		{"priority": 5, "name": "free_food", "description": "free", "maxOrder": 2, "price": 0, "enabled": true, "kindId": kindIds[0], "maxOrderPerDay": 14, "scheduleIds": []string{scheduleIds[0], scheduleIds[1]}, "imageUrl": "", "allowDates": []string{}},
	}
	for index, body := range bodies {
		jBytes, err := json.Marshal(body)
		assert.NoError(t, err, "init json is failed")

		//jsonStr := string(jBytes)
		//fmt.Println("req", jsonStr)

		req, _ := http.NewRequest("PUT", "/item/food/"+foodIds[0], bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		fmt.Println("response", w.Body)

		var idResponse map[string]string
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &idResponse)

		// confirm result from get result
		getReq, _ := http.NewRequest("GET", "/item/food/"+foodIds[0], nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, getReq)

		assert.Equal(t, http.StatusOK, w.Code)
		// fmt.Println("response", w.Body)

		var response map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		AssertMaps(t, response, wants[index])
	}
}

func TestFoodItemHandler_PUT_NotFound(t *testing.T) {
	r := SetupFoodItemRouter()

	scheduleIds := []string{}
	for _, sch := range businessHoursMemory.GetSchedules() {
		scheduleIds = append(scheduleIds, sch.GetId())
	}

	kindIds := []string{}
	for id := range kindMemoryMaps {
		kindIds = append(kindIds, id)
	}

	bodies := []map[string]interface{}{
		{"priority": 3, "name": "create_food1", "description": "desc1", "maxOrder": 1, "price": 1000, "enabled": false, "kindId": kindIds[1], "maxOrderPerDay": 11, "scheduleIds": []string{scheduleIds[2]}, "imageUrl": "https://food1.jpg", "allowDates": []string{}},
		{"priority": 3, "name": "create_food1", "description": "desc1", "maxOrder": 1, "price": 1000, "enabled": false, "kindId": kindIds[1], "maxOrderPerDay": 11, "scheduleIds": []string{scheduleIds[2]}, "imageUrl": "", "allowDates": []string{}},
	}

	jBytes, err := json.Marshal(bodies[0])
	assert.NoError(t, err, "init json is failed")

	req, _ := http.NewRequest("PUT", "/item/food/1234", bytes.NewBuffer(jBytes))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	fmt.Println("response", w.Body)
}

func TestFoodItemHandler_PUT_CREATE_BadRequest(t *testing.T) {
	r := SetupFoodItemRouter()

	foodIds := []string{}
	for id := range foodMemoryMaps {
		foodIds = append(foodIds, id)
	}

	scheduleIds := []string{}
	for _, sch := range businessHoursMemory.GetSchedules() {
		scheduleIds = append(scheduleIds, sch.GetId())
	}

	kinds := []map[string]interface{}{
		{"priority": 1, "name": "item1"},
	}

	kindIds := []string{}
	for id := range kindMemoryMaps {
		kindIds = append(kindIds, id)
	}

	wants := []map[string]interface{}{
		{"priority": 1, "name": "food1", "description": "item1", "maxOrder": 4, "price": 100, "enabled": true, "kind": kinds[0], "maxOrderPerDay": 10, "scheduleIds": []string{scheduleIds[0], scheduleIds[1]}, "imageUrl": "https://food1.jpg", "allowDates": []string{}},
	}

	inputs := GetFoodItemErrorData(kindIds, scheduleIds)
	assert.NotEqual(t, 0, len(inputs), "input data is empty")

	for _, tt := range inputs {
		fmt.Println("case:", tt.name)
		jBytes, err := json.Marshal(tt.args)
		assert.NoError(t, err, "init json is failed")

		req, _ := http.NewRequest("PUT", "/item/food/"+foodIds[0], bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		// confirm item is not update
		req, _ = http.NewRequest("GET", "/item/food/"+foodIds[0], nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		//fmt.Println("re", w.Body)

		assert.Equal(t, http.StatusOK, w.Code, "failed to get result to confirm.")

		var response map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		AssertMaps(t, response, wants[0])
	}
}

func TestFoodItemHandler_DELETE(t *testing.T) {
	r := SetupFoodItemRouter()

	foodIds := []string{}
	for id := range foodMemoryMaps {
		foodIds = append(foodIds, id)
	}
	req, _ := http.NewRequest("DELETE", "/item/food/"+foodIds[0], nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// confirm result
	getReq, _ := http.NewRequest("GET", "/item/food/"+foodIds[0], nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, getReq)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestFoodItemHandler_DELETE_Notfound(t *testing.T) {
	r := SetupFoodItemRouter()

	req, _ := http.NewRequest("DELETE", "/item/food/1234", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
