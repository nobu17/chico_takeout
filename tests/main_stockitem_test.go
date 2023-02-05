package tests

import (
	//"bytes"
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

var stockMemoryMaps map[string]*domains.StockItem

func SetupStockItemRouter() *gin.Engine {
	r := gin.Default()
	kindRepo := memory.NewItemKindMemoryRepository()
	kindMemoryMaps = kindRepo.GetMemory()
	// stock
	stock := r.Group("/item/stock")
	{
		stockRepo := memory.NewStockItemMemoryRepository()
		stockRepo.Reset()
		stockMemoryMaps = stockRepo.GetMemory()
		useCase := itemUseCase.NewStockItemUseCase(stockRepo, kindRepo)
		handler := itemHandler.NewStockItemHandler(useCase)
		stock.GET("/:id", handler.Get)
		stock.GET("/", handler.GetAll)
		stock.POST("/", handler.Post)
		stock.PUT("/:id", handler.Put)
		stock.PUT("/:id/remain", handler.PutRemain)
		stock.DELETE("/:id", handler.Delete)
	}
	return r
}

func TestStockItemHandler_GETALL(t *testing.T) {

	kinds := []map[string]interface{}{
		{"priority": 1, "name": "item1"},
		{"priority": 2, "name": "item2"},
	}
	wants := []map[string]interface{}{
		{"priority": 1, "name": "stock1", "description": "item1", "maxOrder": 4, "price": 100, "enabled": true, "remain": 10, "kind": kinds[0], "imageUrl": "https://item1.png"},
		{"priority": 2, "name": "stock2", "description": "item2", "maxOrder": 5, "price": 200, "enabled": true, "remain": 0, "kind": kinds[1], "imageUrl": ""},
	}

	r := SetupStockItemRouter()
	req, _ := http.NewRequest("GET", "/item/stock/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Convert the JSON response to a map
	var response []map[string]interface{}
	_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

	assert.Equal(t, len(stockMemoryMaps), len(response), "response length is incorrect")

	for index, item := range response {
		fmt.Println("item", item)
		AssertMaps(t, item, wants[index])
	}
}

func TestStockItemHandler_GET(t *testing.T) {
	r := SetupStockItemRouter()

	kinds := []map[string]interface{}{
		{"priority": 1, "name": "item1"},
		{"priority": 2, "name": "item2"},
	}
	wants := []map[string]interface{}{
		{"priority": 1, "name": "stock1", "description": "item1", "maxOrder": 4, "price": 100, "enabled": true, "remain": 10, "kind": kinds[0], "imageUrl": "https://item1.png"},
		{"priority": 2, "name": "stock2", "description": "item2", "maxOrder": 5, "price": 200, "enabled": true, "remain": 0, "kind": kinds[1], "imageUrl": ""},
	}

	index := 0
	for id := range stockMemoryMaps {
		req, _ := http.NewRequest("GET", "/item/stock/"+id, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		fmt.Println("response", response)
		AssertMaps(t, response, wants[index])
		index++
	}
}

func TestStockItemHandler_GET_NotFound(t *testing.T) {
	r := SetupStockItemRouter()

	req, _ := http.NewRequest("GET", "/item/stock/1234", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

	assert.Nil(t, response)
}

func TestStockItemHandler_POST_CREATE(t *testing.T) {
	r := SetupStockItemRouter()

	kindIds := []string{}
	for id := range kindMemoryMaps {
		kindIds = append(kindIds, id)
	}
	kinds := []map[string]interface{}{
		{"priority": 1, "name": "item1"},
		{"priority": 2, "name": "item2"},
	}

	wants := []map[string]interface{}{
		{"priority": 1, "name": "stock1", "description": "item1", "maxOrder": 4, "price": 100, "enabled": false, "remain": 0, "kind": kinds[0], "imageUrl": "https://hoge.png"},
		{"priority": 2, "name": "stock2", "description": "item2", "maxOrder": 5, "price": 200, "enabled": true, "remain": 0, "kind": kinds[1], "imageUrl": ""},
	}
	bodies := []map[string]interface{}{
		{"priority": 1, "name": "stock1", "description": "item1", "maxOrder": 4, "price": 100, "enabled": false, "kindId": kindIds[0], "imageUrl": "https://hoge.png"},
		{"priority": 2, "name": "stock2", "description": "item2", "maxOrder": 5, "price": 200, "enabled": true, "kindId": kindIds[1], "imageUrl": ""},
	}
	for index, body := range bodies {
		jBytes, err := json.Marshal(body)
		assert.NoError(t, err, "init json is failed")

		jsonStr := string(jBytes)
		fmt.Println("req", jsonStr)

		req, _ := http.NewRequest("POST", "/item/stock/", bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var idResponse map[string]string
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &idResponse)
		fmt.Println("response", w.Body)

		id := idResponse["id"]
		assert.NotEmpty(t, id, "response id should not be empty.")

		// confirm result from result id
		getReq, _ := http.NewRequest("GET", "/item/stock/"+id, nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, getReq)

		assert.Equal(t, http.StatusOK, w.Code)
		fmt.Println("response", w.Body)

		var response map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		AssertMaps(t, response, wants[index])
	}
}

type stockItemErrorData struct {
	name string
	args map[string]interface{}
	want int
}

func GetStockItemErrorData(kindIds []string) []stockItemErrorData {
	var stockItemErrorInputs = []stockItemErrorData{
		{name: "lack priority", args: map[string]interface{}{
			"name": "stock1", "description": "item1",
			"maxOrder": 4, "price": 100, "enabled": true, "kindId": kindIds[0], "imageUrl": "https://image.png",
		}, want: 2},
		{name: "error priority(0)", args: map[string]interface{}{
			"priority": 0, "name": "stock1", "description": "item1",
			"maxOrder": 4, "price": 100, "enabled": true, "kindId": kindIds[0], "imageUrl": "https://image.png",
		}, want: 2},
		{name: "error priority(-1)", args: map[string]interface{}{
			"priority": -1, "name": "stock1", "description": "item1",
			"maxOrder": 4, "price": 100, "enabled": true, "kindId": kindIds[0], "imageUrl": "https://image.png",
		}, want: 2},
		{name: "lack name", args: map[string]interface{}{
			"priority": 1, "description": "item1",
			"maxOrder": 4, "price": 100, "enabled": true, "kindId": kindIds[0], "imageUrl": "https://image.png",
		}, want: 2},
		{name: "error name(empty)", args: map[string]interface{}{
			"priority": 1, "name": "", "description": "item1",
			"maxOrder": 4, "price": 100, "enabled": true, "kindId": kindIds[0], "imageUrl": "https://image.png",
		}, want: 2},
		{name: "error name(over limit(26))", args: map[string]interface{}{
			"priority": 1, "name": MakeRandomStr(26), "description": "item1",
			"maxOrder": 4, "price": 100, "enabled": true, "kindId": kindIds[0], "imageUrl": "https://image.png",
		}, want: 2},
		{name: "lack description", args: map[string]interface{}{
			"priority": 1, "name": "stock1",
			"maxOrder": 4, "price": 100, "enabled": true, "kindId": kindIds[0], "imageUrl": "https://image.png",
		}, want: 2},
		{name: "error description(empty)", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "",
			"maxOrder": 4, "price": 100, "enabled": true, "kindId": kindIds[0], "imageUrl": "https://image.png",
		}, want: 2},
		{name: "error description(over limit(151))", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": MakeRandomStr(151),
			"maxOrder": 4, "price": 100, "enabled": true, "kindId": kindIds[0], "imageUrl": "https://image.png",
		}, want: 2},
		{name: "lack maxOrder", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "item1",
			"price": 100, "enabled": true, "kindId": kindIds[0], "imageUrl": "https://image.png",
		}, want: 2},
		{name: "error maxOrder(0)", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "desc",
			"maxOrder": 0, "price": 100, "enabled": true, "kindId": kindIds[0], "imageUrl": "https://image.png",
		}, want: 2},
		{name: "error maxOrder(31)", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "desc",
			"maxOrder": 31, "price": 100, "enabled": true, "kindId": kindIds[0], "imageUrl": "https://image.png",
		}, want: 2},
		{name: "lack price", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "item1",
			"maxOrder": 4, "enabled": true, "kindId": kindIds[0], "imageUrl": "https://image.png",
		}, want: 2},
		{name: "error price(0)", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "desc",
			"maxOrder": 1, "price": 0, "enabled": true, "kindId": kindIds[0], "imageUrl": "https://image.png",
		}, want: 2},
		{name: "error price(20001)", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "desc",
			"maxOrder": 1, "price": 20001, "enabled": true, "kindId": kindIds[0], "imageUrl": "https://image.png",
		}, want: 2},
		{name: "lack enabled", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "item1",
			"maxOrder": 4, "price": 100, "kindId": kindIds[0], "imageUrl": "https://image.png",
		}, want: 2},
		{name: "lack kindId", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "item1",
			"maxOrder": 4, "price": 100, "enabled": true, "imageUrl": "https://image.png",
		}, want: 2},
		{name: "error kindId(empty)", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "desc",
			"maxOrder": 1, "price": 100, "enabled": true, "kindId": "", "imageUrl": "https://image.png",
		}, want: 2},
		{name: "not exist kindId", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "item1",
			"maxOrder": 4, "price": 100, "enabled": true, "kindId": "abc", "imageUrl": "https://image.png",
		}, want: 2},
		{name: "lack imageUrl", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "item1",
			"maxOrder": 4, "price": 100, "kindId": kindIds[0], "enabled": true,
		}, want: 2},
		{name: "incorrect imageUrl format", args: map[string]interface{}{
			"priority": 1, "name": "stock1", "description": "item1",
			"maxOrder": 4, "price": 100, "kindId": kindIds[0], "enabled": true, "imageUrl": "image.png",
		}, want: 2},
	}
	return stockItemErrorInputs
}

func TestStockItemHandler_POST_CREATE_BadRequest(t *testing.T) {
	r := SetupStockItemRouter()

	kindIds := []string{}
	for id := range kindMemoryMaps {
		kindIds = append(kindIds, id)
	}

	inputs := GetStockItemErrorData(kindIds)
	assert.NotEqual(t, 0, len(inputs), "input data is empty")

	for _, tt := range inputs {
		fmt.Println("case:", tt.name)
		jBytes, err := json.Marshal(tt.args)
		assert.NoError(t, err, "init json is failed")

		req, _ := http.NewRequest("POST", "/item/stock/", bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		// confirm stock is not added
		req, _ = http.NewRequest("GET", "/item/stock/", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "failed to get result to confirm.")

		var response []map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		assert.Equal(t, tt.want, len(stockMemoryMaps), "response length is incorrect")
	}
}

func TestStockItemHandler_PUT(t *testing.T) {
	r := SetupStockItemRouter()

	kindIds := []string{}
	for id := range kindMemoryMaps {
		kindIds = append(kindIds, id)
	}
	stockIds := []string{}
	for id := range stockMemoryMaps {
		stockIds = append(stockIds, id)
	}
	kinds := []map[string]interface{}{
		{"priority": 1, "name": "item1"},
		{"priority": 2, "name": "item2"},
	}

	wants := []map[string]interface{}{
		{"priority": 3, "name": "change1", "description": "desc123", "maxOrder": 8, "price": 1000, "enabled": true, "remain": 10, "kind": kinds[1], "imageUrl": "https://image.png"},
		{"priority": 4, "name": "change2", "description": "desc321", "maxOrder": 9, "price": 2000, "enabled": false, "remain": 10, "kind": kinds[0], "imageUrl": ""},
	}
	bodies := []map[string]interface{}{
		{"priority": 3, "name": "change1", "description": "desc123", "maxOrder": 8, "price": 1000, "enabled": true, "kindId": kindIds[1], "imageUrl": "https://image.png"},
		{"priority": 4, "name": "change2", "description": "desc321", "maxOrder": 9, "price": 2000, "enabled": false, "kindId": kindIds[0], "imageUrl": ""},
	}
	for index, body := range bodies {
		jBytes, err := json.Marshal(body)
		assert.NoError(t, err, "init json is failed")

		jsonStr := string(jBytes)
		fmt.Println("req", jsonStr)

		req, _ := http.NewRequest("PUT", "/item/stock/"+stockIds[0], bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		fmt.Println("response", w.Body)

		var idResponse map[string]string
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &idResponse)

		// confirm result from get result
		getReq, _ := http.NewRequest("GET", "/item/stock/"+stockIds[0], nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, getReq)

		assert.Equal(t, http.StatusOK, w.Code)
		fmt.Println("response", w.Body)

		var response map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		AssertMaps(t, response, wants[index])
	}
}

func TestStockItemHandler_PUT_NotFound(t *testing.T) {
	r := SetupStockItemRouter()

	kindIds := []string{}
	for id := range kindMemoryMaps {
		kindIds = append(kindIds, id)
	}

	bodies := []map[string]interface{}{
		{"priority": 1, "name": "stock1", "description": "item1", "maxOrder": 4, "price": 100, "enabled": false, "kindId": kindIds[0], "imageUrl": "https://image.png"},
	}

	jBytes, err := json.Marshal(bodies[0])
	assert.NoError(t, err, "init json is failed")

	req, _ := http.NewRequest("PUT", "/item/stock/1234", bytes.NewBuffer(jBytes))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	fmt.Println("response", w.Body)
}

func TestStockItemHandler_PUT_BadRequest(t *testing.T) {
	r := SetupStockItemRouter()

	kindIds := []string{}
	for id := range kindMemoryMaps {
		kindIds = append(kindIds, id)
	}
	stockIds := []string{}
	for id := range stockMemoryMaps {
		stockIds = append(stockIds, id)
	}

	kinds := []map[string]interface{}{
		{"priority": 1, "name": "item1"},
	}

	wants := []map[string]interface{}{
		{"priority": 1, "name": "stock1", "description": "item1", "maxOrder": 4, "price": 100, "enabled": true, "remain": 10, "kind": kinds[0], "imageUrl": "https://item1.png"},
	}

	inputs := GetStockItemErrorData(kindIds)
	assert.NotEqual(t, 0, len(inputs), "input data is empty")

	for _, tt := range inputs {
		fmt.Println("case:", tt.name)
		jBytes, err := json.Marshal(tt.args)
		assert.NoError(t, err, "init json is failed")

		req, _ := http.NewRequest("PUT", "/item/stock/"+stockIds[0], bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		// confirm stock is not changed
		req, _ = http.NewRequest("GET", "/item/stock/"+stockIds[0], nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		fmt.Println("response", w.Body)

		assert.Equal(t, http.StatusOK, w.Code, "failed to get result to confirm.")

		var response map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		AssertMaps(t, response, wants[0])
	}
}

func TestStockItemHandler_DELETE(t *testing.T) {
	r := SetupStockItemRouter()

	stockIds := []string{}
	for id := range stockMemoryMaps {
		stockIds = append(stockIds, id)
	}

	req, _ := http.NewRequest("DELETE", "/item/stock/"+stockIds[0], nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// confirm result
	getReq, _ := http.NewRequest("GET", "/item/stock/"+stockIds[0], nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, getReq)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestStockItemHandler_DELETE_NotFound(t *testing.T) {
	r := SetupStockItemRouter()

	req, _ := http.NewRequest("DELETE", "/item/stock/1234", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	// confirm result
	getReq, _ := http.NewRequest("GET", "/item/stock/1234", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, getReq)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestStockItemHandler_PUTRemain(t *testing.T) {
	r := SetupStockItemRouter()

	stockIds := []string{}
	for id := range stockMemoryMaps {
		stockIds = append(stockIds, id)
	}
	kinds := []map[string]interface{}{
		{"priority": 1, "name": "item1"},
	}

	wants := []map[string]interface{}{
		{"priority": 1, "name": "stock1", "description": "item1", "maxOrder": 4, "price": 100, "enabled": true, "remain": 20, "kind": kinds[0]},
		{"priority": 1, "name": "stock1", "description": "item1", "maxOrder": 4, "price": 100, "enabled": true, "remain": 5, "kind": kinds[0]},
	}
	bodies := []map[string]interface{}{
		{"remain": 20},
		{"remain": 5},
	}
	for index, body := range bodies {
		jBytes, err := json.Marshal(body)
		assert.NoError(t, err, "init json is failed")

		jsonStr := string(jBytes)
		fmt.Println("req", jsonStr)

		req, _ := http.NewRequest("PUT", "/item/stock/"+stockIds[0]+"/remain", bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		fmt.Println("response", w.Body)

		var idResponse map[string]string
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &idResponse)

		// confirm result from get result
		getReq, _ := http.NewRequest("GET", "/item/stock/"+stockIds[0], nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, getReq)

		assert.Equal(t, http.StatusOK, w.Code)
		fmt.Println("response", w.Body)

		var response map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		AssertMaps(t, response, wants[index])
	}
}

func TestStockItemHandler_PUTRemain_BadRequest(t *testing.T) {
	r := SetupStockItemRouter()

	stockIds := []string{}
	for id := range stockMemoryMaps {
		stockIds = append(stockIds, id)
	}

	inputs := []stockItemErrorData{
		{name: "error value(0)", args: map[string]interface{}{
			"remain": 0,
		}, want: 10},
		{name: "error over limits(1000)", args: map[string]interface{}{
			"remain": 1000,
		}, want: 10},
	}

	for _, input := range inputs {
		jBytes, err := json.Marshal(input.args)
		assert.NoError(t, err, "init json is failed")

		jsonStr := string(jBytes)
		fmt.Println("req", jsonStr)

		req, _ := http.NewRequest("PUT", "/item/stock/"+stockIds[0]+"/remain", bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		fmt.Println("response", w.Body)
	}
}

func TestStockItemHandler_PUTRemain_NotFound(t *testing.T) {
	r := SetupStockItemRouter()

	bodies := []map[string]interface{}{
		{"remain": 20},
	}
	jBytes, err := json.Marshal(bodies[0])
	assert.NoError(t, err, "init json is failed")

	jsonStr := string(jBytes)
	fmt.Println("req", jsonStr)

	req, _ := http.NewRequest("PUT", "/item/stock/1234/remain", bytes.NewBuffer(jBytes))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	fmt.Println("response", w.Body)
}
