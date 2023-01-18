package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	itemHandler "chico/takeout/handlers/item"
	"chico/takeout/infrastructures/memory"
	itemUseCase "chico/takeout/usecase/item"

	"github.com/gin-gonic/gin"
)

func SetupOptionItemRouter() *gin.Engine {
	r := gin.Default()
	kind := r.Group("/item/option")
	{
		repo := memory.NewOptionItemMemoryRepository()
		repo.Reset()
		// kindMemoryMaps = repo.GetMemory()
		useCase := itemUseCase.NewOptionItemUseCase(repo)
		handler := itemHandler.NewOptionItemHandler(useCase)
		kind.GET("/:id", handler.Get)
		kind.GET("/", handler.GetAll)
		kind.POST("/", handler.Post)
		kind.PUT("/:id", handler.Put)
		kind.DELETE("/:id", handler.Delete)
	}
	return r
}

func TestOptionItemHandler_GET_ALL(t *testing.T) {
	r := SetupOptionItemRouter()
	req, _ := http.NewRequest("GET", "/item/option/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if http.StatusOK != w.Code {
		t.Errorf("Status Code should be OK:%d", w.Code)
		return
	}
	// Convert the JSON response to a map
	var response []map[string]interface{}
	_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

	if len(response) != 3 {
		t.Errorf("response length should be %d, actual:%d", len(kindMemoryMaps), len(response))
		return
	}
	// var want []map[string]interface{}
	wants := []map[string]interface{}{
		{"id": "1", "priority": 1, "name": "item1", "description": "memo1", "price": 100, "enabled": true},
		{"id": "2", "priority": 2, "name": "item2", "description": "memo2", "price": 200, "enabled": true},
		{"id": "3", "priority": 3, "name": "item3", "description": "memo3", "price": 300, "enabled": false},
	}
	fmt.Println("response", response)
	for index, item := range response {
		fmt.Println("item", item)
		fmt.Println("wants[index]", wants[index])
		AssertMaps(t, item, wants[index])
	}
}

func TestOptionItemHandler_GET(t *testing.T) {
	r := SetupOptionItemRouter()
	req, _ := http.NewRequest("GET", "/item/option/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if http.StatusOK != w.Code {
		t.Errorf("Status Code should be OK:%d", w.Code)
		return
	}
	// Convert the JSON response to a map
	var response map[string]interface{}
	_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

	wants := []map[string]interface{}{
		{"id": "1", "priority": 1, "name": "item1", "description": "memo1", "price": 100, "enabled": true},
	}

	AssertMaps(t, response, wants[0])
}

func TestOptionItemHandler_GET_NotFound(t *testing.T) {
	r := SetupOptionItemRouter()
	req, _ := http.NewRequest("GET", "/item/option/1234", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if http.StatusNotFound != w.Code {
		t.Errorf("Status Code should be NotFound:%d", w.Code)
		return
	}
}

func TestOptionItemHandler_POST(t *testing.T) {
	r := SetupOptionItemRouter()

	want := map[string]interface{}{"priority": 1, "name": "itemP", "description": "memoP", "price": 123, "enabled": true}
	body := map[string]interface{}{"priority": 1, "name": "itemP", "description": "memoP", "price": 123, "enabled": true}
	jBytes, err := json.Marshal(body)
	if err != nil {
		t.Errorf("failed to marshal json for test.:%s", err)
		return
	}
	jsonStr := string(jBytes)
	fmt.Println(jsonStr)

	req, _ := http.NewRequest("POST", "/item/option/", bytes.NewBuffer(jBytes))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if http.StatusOK != w.Code {
		t.Errorf("Status Code should be OK:%d", w.Code)
		return
	}

	var idResponse map[string]string
	_ = json.Unmarshal([]byte(w.Body.Bytes()), &idResponse)

	id := idResponse["id"]
	if id == "" {
		t.Errorf("response id should no be empty")
		return
	}

	// confirm result from result id
	getReq, _ := http.NewRequest("GET", "/item/option/"+id, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, getReq)
	if http.StatusOK != w.Code {
		t.Errorf("Status Code should be OK:%d", w.Code)
	}
	fmt.Println("response", w.Body)

	var response map[string]interface{}
	_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

	AssertMaps(t, response, want)
}

func TestOptionItemHandler_POST_BadRequest(t *testing.T) {
	type input struct {
		name string
		args map[string]interface{}
		want int
	}

	r := SetupOptionItemRouter()

	inputs := []input{
		{name: "lack name", args: map[string]interface{}{"priority": 1, "description": "memoP", "price": 123, "enabled": true}, want: 3},
		{name: "lack priority", args: map[string]interface{}{"name": "1", "description": "memoP", "price": 123, "enabled": true}, want: 3},
		{name: "lack description", args: map[string]interface{}{"name": "1", "priority": 1, "price": 123, "enabled": true}, want: 3},
		{name: "lack price", args: map[string]interface{}{"name": "1", "priority": 1, "description": "123", "enabled": true}, want: 3},
		{name: "lack enabled", args: map[string]interface{}{"name": "1", "priority": 1, "description": "123", "price": 123}, want: 3},
		{name: "long name(26)", args: map[string]interface{}{"name": "12345678901234567890123455", "priority": 1, "description": "memoP", "price": 123, "enabled": true}, want: 3},
		{name: "0 price", args: map[string]interface{}{"name": "123", "priority": 1, "description": "memoP", "price": 0, "enabled": true}, want: 3},
		{name: "minus price", args: map[string]interface{}{"name": "123", "priority": 1, "description": "memoP", "price": -1, "enabled": true}, want: 3},
		{name: "minus priority", args: map[string]interface{}{"name": "123", "priority": -1, "description": "memoP", "price": 9, "enabled": true}, want: 3},
	}

	for _, param := range inputs {
		jBytes, err := json.Marshal(param.args)
		if err != nil {
			t.Errorf("failed to marshal json for test.:%s", err)
			return
		}
		jsonStr := string(jBytes)
		fmt.Println(jsonStr)

		req, _ := http.NewRequest("POST", "/item/option/", bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if http.StatusBadRequest != w.Code {
			t.Errorf("Status Code should be BadRequest:%d", w.Code)
			return
		}

		// confirm result
		getReq, _ := http.NewRequest("GET", "/item/option/", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, getReq)
		if http.StatusOK != w.Code {
			t.Errorf("Status Code should be OK:%d", w.Code)
		}
		fmt.Println("response", w.Body)

		var response []map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		if len(response) != param.want {
			t.Errorf("item should not ba add. expected:%d, actual:%d", param.want, len(response))
			return
		}
	}
}

func TestOptionItemHandler_PUT(t *testing.T) {
	r := SetupOptionItemRouter()

	want := map[string]interface{}{"priority": 2, "name": "itemPP", "description": "memoPP", "price": 1234, "enabled": false}
	body := map[string]interface{}{"priority": 2, "name": "itemPP", "description": "memoPP", "price": 1234, "enabled": false}
	jBytes, err := json.Marshal(body)
	if err != nil {
		t.Errorf("failed to marshal json for test.:%s", err)
		return
	}
	jsonStr := string(jBytes)
	fmt.Println(jsonStr)

	req, _ := http.NewRequest("PUT", "/item/option/1", bytes.NewBuffer(jBytes))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if http.StatusOK != w.Code {
		fmt.Println("response", w.Body)
		t.Errorf("Status Code should be OK:%d", w.Code)
		return
	}

	// confirm result from result id
	getReq, _ := http.NewRequest("GET", "/item/option/1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, getReq)
	if http.StatusOK != w.Code {
		t.Errorf("Status Code should be OK:%d", w.Code)
	}
	fmt.Println("response", w.Body)

	var response map[string]interface{}
	_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

	AssertMaps(t, response, want)
}

func TestOptionItemHandler_PUT_BadRequest(t *testing.T) {
	type input struct {
		name string
		args map[string]interface{}
		want int
	}

	r := SetupOptionItemRouter()

	inputs := []input{
		{name: "lack name", args: map[string]interface{}{"priority": 1, "description": "memoP", "price": 123, "enabled": true}, want: 3},
		{name: "lack priority", args: map[string]interface{}{"name": "1", "description": "memoP", "price": 123, "enabled": true}, want: 3},
		{name: "lack description", args: map[string]interface{}{"name": "1", "priority": 1, "price": 123, "enabled": true}, want: 3},
		{name: "lack price", args: map[string]interface{}{"name": "1", "priority": 1, "description": "123", "enabled": true}, want: 3},
		{name: "lack enabled", args: map[string]interface{}{"name": "1", "priority": 1, "description": "123", "price": 123}, want: 3},
		{name: "long name(26)", args: map[string]interface{}{"name": "12345678901234567890123455", "priority": 1, "description": "memoP", "price": 123, "enabled": true}, want: 3},
		{name: "0 price", args: map[string]interface{}{"name": "123", "priority": 1, "description": "memoP", "price": 0, "enabled": true}, want: 3},
		{name: "minus price", args: map[string]interface{}{"name": "123", "priority": 1, "description": "memoP", "price": -1, "enabled": true}, want: 3},
		{name: "minus priority", args: map[string]interface{}{"name": "123", "priority": -1, "description": "memoP", "price": 9, "enabled": true}, want: 3},
	}

	for _, param := range inputs {
		jBytes, err := json.Marshal(param.args)
		if err != nil {
			t.Errorf("failed to marshal json for test.:%s", err)
			return
		}
		jsonStr := string(jBytes)
		fmt.Println(jsonStr)

		req, _ := http.NewRequest("PUT", "/item/option/1", bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if http.StatusBadRequest != w.Code {
			t.Errorf("Status Code should be BadRequest:%d", w.Code)
			return
		}

		// confirm result
		getReq, _ := http.NewRequest("GET", "/item/option/", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, getReq)
		if http.StatusOK != w.Code {
			t.Errorf("Status Code should be OK:%d", w.Code)
		}
		fmt.Println("response", w.Body)

		var response []map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		if len(response) != param.want {
			t.Errorf("item should not ba add. expected:%d, actual:%d", param.want, len(response))
			return
		}
	}
}

func TestOptionItemHandler_PUT_NotFound(t *testing.T) {
	r := SetupOptionItemRouter()

	body := map[string]interface{}{"priority": 2, "name": "itemPP", "description": "memoPP", "price": 1234, "enabled": false}
	jBytes, err := json.Marshal(body)
	if err != nil {
		t.Errorf("failed to marshal json for test.:%s", err)
		return
	}
	jsonStr := string(jBytes)
	fmt.Println(jsonStr)

	req, _ := http.NewRequest("PUT", "/item/option/1234", bytes.NewBuffer(jBytes))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if http.StatusNotFound != w.Code {
		fmt.Println("response", w.Body)
		t.Errorf("Status Code should be NotFound:%d", w.Code)
		return
	}
}

func TestOptionItemHandler_DELETE(t *testing.T) {
	r := SetupOptionItemRouter()

	req, _ := http.NewRequest("DELETE", "/item/option/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if http.StatusOK != w.Code {
		t.Errorf("Status Code should be OK:%d", w.Code)
		return
	}

	// confirm result
	getReq, _ := http.NewRequest("GET", "/item/option/1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, getReq)
	if http.StatusNotFound != w.Code {
		t.Errorf("Status Code should be NotFound after delete:%d", w.Code)
	}
}

func TestOptionItemHandler_DELETE_NotFound(t *testing.T) {
	r := SetupOptionItemRouter()

	req, _ := http.NewRequest("DELETE", "/item/option/1234", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if http.StatusNotFound != w.Code {
		t.Errorf("Status Code should be NotFound:%d", w.Code)
		return
	}
}
