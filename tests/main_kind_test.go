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
)

var kindMemoryMaps map[string]*domains.ItemKind

func SetupKindRouter() *gin.Engine {
	r := gin.Default()
	kind := r.Group("/item/kind")
	{
		repo := memory.NewItemKindMemoryRepository()
		repo.Reset()
		kindMemoryMaps = repo.GetMemory()
		useCase := itemUseCase.NewItemKindUseCase(repo)
		handler := itemHandler.NewItemKindHandler(useCase)
		kind.GET("/:id", handler.Get)
		kind.GET("/", handler.GetAll)
		kind.POST("/", handler.Post)
		kind.PUT("/:id", handler.Put)
		kind.DELETE("/:id", handler.Delete)
	}
	return r
}

func TestItemKindHandler_GETALL(t *testing.T) {
	r := SetupKindRouter()
	req, _ := http.NewRequest("GET", "/item/kind/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if http.StatusOK != w.Code {
		t.Errorf("Status Code should be OK:%d", w.Code)
		return
	}
	// Convert the JSON response to a map
	var response []map[string]interface{}
	_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

	if len(response) != len(kindMemoryMaps) {
		t.Errorf("response length should be %d, actual:%d", len(kindMemoryMaps), len(response))
		return
	}
	// var want []map[string]interface{}
	wants := []map[string]interface{}{
		{"priority": 1, "name": "item1"},
		{"priority": 2, "name": "item2"},
	}
	for index, item := range response {
		fmt.Println("item", item)
		AssertMaps(t, item, wants[index])
	}
}

func TestItemKindHandler_GET(t *testing.T) {
	r := SetupKindRouter()

	wants := map[string]map[string]interface{}{
		"item1": {"priority": 1, "name": "item1"},
		"item2": {"priority": 2, "name": "item2"},
	}

	for id := range kindMemoryMaps {
		req, _ := http.NewRequest("GET", "/item/kind/"+id, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if http.StatusOK != w.Code {
			t.Errorf("Status Code should be OK:%d", w.Code)
			return
		}

		var response map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		nameKey := response["name"].(string)
		AssertMaps(t, response, wants[nameKey])
	}
}

func TestItemKindHandler_GET_NotFound(t *testing.T) {
	r := SetupKindRouter()

	req, _ := http.NewRequest("GET", "/item/kind/12345", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if http.StatusNotFound != w.Code {
		t.Errorf("Status Code should be NotFound:%d", w.Code)
		return
	}
	var response map[string]interface{}
	_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)
	if response != nil {
		t.Errorf("response should be null:%s", response)
		return
	}
}

func TestItemKindHandler_POST(t *testing.T) {
	r := SetupKindRouter()

	want := map[string]interface{}{"priority": 4, "name": "add"}
	body := map[string]interface{}{
		"name":     "add",
		"priority": 4,
	}
	jBytes, err := json.Marshal(body)
	if err != nil {
		t.Errorf("failed to marshal json for test.:%s", err)
		return
	}
	jsonStr := string(jBytes)
	fmt.Println(jsonStr)

	req, _ := http.NewRequest("POST", "/item/kind/", bytes.NewBuffer(jBytes))
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
	getReq, _ := http.NewRequest("GET", "/item/kind/"+id, nil)
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

func TestItemKindHandler_POST_BadRequest(t *testing.T) {
	type input struct {
		name string
		args map[string]interface{}
		want int
	}

	r := SetupKindRouter()

	inputs := []input{
		{name: "only name", args: map[string]interface{}{"name": "added"}, want: 2},
		{name: "only priority", args: map[string]interface{}{"priority": 3}, want: 2},
		{name: "empty", args: map[string]interface{}{}, want: 2},
	}

	for _, param := range inputs {
		jBytes, err := json.Marshal(param.args)
		if err != nil {
			t.Errorf("failed to marshal json for test.:%s", err)
			return
		}
		jsonStr := string(jBytes)
		fmt.Println(jsonStr)

		req, _ := http.NewRequest("POST", "/item/kind/", bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if http.StatusBadRequest != w.Code {
			t.Errorf("Status Code should be BadRequest:%d", w.Code)
			return
		}

		// confirm result
		getReq, _ := http.NewRequest("GET", "/item/kind/", nil)
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

func TestItemKindHandler_PUT(t *testing.T) {
	r := SetupKindRouter()

	ids := []string{}
	for id := range kindMemoryMaps {
		ids = append(ids, id)
	}

	want := map[string]interface{}{"priority": 3, "name": "changed"}
	body := map[string]interface{}{
		"name":     "changed",
		"priority": 3,
	}
	jBytes, err := json.Marshal(body)
	if err != nil {
		t.Errorf("failed to marshal json for test.:%s", err)
		return
	}
	jsonStr := string(jBytes)
	fmt.Println(jsonStr)

	req, _ := http.NewRequest("PUT", "/item/kind/"+ids[0], bytes.NewBuffer(jBytes))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if http.StatusOK != w.Code {
		t.Errorf("Status Code should be OK:%d", w.Code)
		return
	}

	// confirm result
	getReq, _ := http.NewRequest("GET", "/item/kind/"+ids[0], nil)
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

func TestItemKindHandler_PUT_BadRequest(t *testing.T) {
	type input struct {
		name string
		args map[string]interface{}
		want map[string]interface{}
	}

	r := SetupKindRouter()

	ids := []string{}
	for id := range kindMemoryMaps {
		ids = append(ids, id)
	}

	notchangedWant := map[string]interface{}{"priority": 1, "name": "item1"}

	inputs := []input{
		{name: "only name", args: map[string]interface{}{"name": "changed"}, want: notchangedWant},
		{name: "only priority", args: map[string]interface{}{"priority": 3}, want: notchangedWant},
		{name: "empty", args: map[string]interface{}{}, want: notchangedWant},
	}

	for _, param := range inputs {
		jBytes, err := json.Marshal(param.args)
		if err != nil {
			t.Errorf("failed to marshal json for test.:%s", err)
			return
		}
		jsonStr := string(jBytes)
		fmt.Println(jsonStr)

		req, _ := http.NewRequest("PUT", "/item/kind/"+ids[0], bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if http.StatusBadRequest != w.Code {
			t.Errorf("Status Code should be BadRequest:%d", w.Code)
			return
		}

		// confirm result
		getReq, _ := http.NewRequest("GET", "/item/kind/"+ids[0], nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, getReq)
		if http.StatusOK != w.Code {
			t.Errorf("Status Code should be OK:%d", w.Code)
		}
		fmt.Println("response", w.Body)

		var response map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		AssertMaps(t, response, param.want)
	}
}

func TestItemKindHandler_PUT_NotFound(t *testing.T) {

	body := map[string]interface{}{
		"name":     "changed",
		"priority": 3,
	}

	r := SetupKindRouter()
	jBytes, err := json.Marshal(body)
	if err != nil {
		t.Errorf("failed to marshal json for test.:%s", err)
		return
	}
	jsonStr := string(jBytes)
	fmt.Println(jsonStr)

	req, _ := http.NewRequest("PUT", "/item/kind/1234", bytes.NewBuffer(jBytes))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if http.StatusNotFound != w.Code {
		t.Errorf("Status Code should be NotFound:%d", w.Code)
		return
	}
}

func TestItemKindHandler_DELETE(t *testing.T) {
	r := SetupKindRouter()

	ids := []string{}
	for id := range kindMemoryMaps {
		ids = append(ids, id)
	}

	req, _ := http.NewRequest("DELETE", "/item/kind/"+ids[0], nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if http.StatusOK != w.Code {
		t.Errorf("Status Code should be OK:%d", w.Code)
		return
	}

	// confirm result
	getReq, _ := http.NewRequest("GET", "/item/kind/"+ids[0], nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, getReq)
	if http.StatusNotFound != w.Code {
		t.Errorf("Status Code should be NotFound after delete:%d", w.Code)
	}
}

func TestItemKindHandler_DELETE_NotFound(t *testing.T) {
	r := SetupKindRouter()

	req, _ := http.NewRequest("DELETE", "/item/kind/1234", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if http.StatusNotFound != w.Code {
		t.Errorf("Status Code should be NotFound:%d", w.Code)
		return
	}
}
