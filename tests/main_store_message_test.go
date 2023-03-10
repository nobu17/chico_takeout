package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	handler "chico/takeout/handlers/message"
	"chico/takeout/infrastructures/memory"
	useCase "chico/takeout/usecase/message"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetupStoreMessageRouter() *gin.Engine {
	r := gin.Default()
	message := r.Group("/message/store")
	{
		messageRepo := memory.NewStoreMessageRepository()
		messageRepo.Reset()
		useCase := useCase.NewStoreMessageUseCase(messageRepo)
		err := useCase.CreateInitialMessage()
		if err != nil {
			panic("unexpected error")
		}
		handler := handler.NewStoreMessageHandler(useCase)
		message.GET("/:id", handler.Get)
		message.POST("/", handler.Post)
		message.PUT("/:id", handler.Put)
	}
	return r
}

func TestStoreMessageHandler_GET(t *testing.T) {
	r := SetupStoreMessageRouter()

	inputs := []string{"1", "2"}
	wants := []map[string]interface{}{
		{"id": "1", "content": "トップメッセージです。"},
		{"id": "2", "content": "マイページメッセージです。"},
	}

	for idx, id := range inputs {
		req, _ := http.NewRequest("GET", "/message/store/"+id, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if http.StatusOK != w.Code {
			t.Errorf("Status Code should be OK:%d", w.Code)
			return
		}

		var response map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		fmt.Println(response)

		AssertMaps(t, response, wants[idx])
	}
}

func TestStoreMessageHandler_GET_NotFound(t *testing.T) {
	r := SetupStoreMessageRouter()

	req, _ := http.NewRequest("GET", "/message/store/3", nil)
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

func TestStoreMessageHandler_PUT_UPDATE(t *testing.T) {
	r := SetupStoreMessageRouter()

	stockIds := map[string]string{}
	for id, value := range stockMemoryMaps {
		stockIds[value.GetName()] = id
	}
	foodIds := map[string]string{}
	for id, value := range foodMemoryMaps {
		foodIds[value.GetName()] = id
	}

	ids := []string{"1", "2"}
	bodies := []map[string]interface{}{
		{"id": "1", "content": "アップデート"},
		{"id": "2", "content": "アップデートXXX"},
	}

	wants := []map[string]interface{}{
		{"id": "1", "content": "アップデート"},
		{"id": "2", "content": "アップデートXXX"},
	}
	for index, body := range bodies {
		jBytes, err := json.Marshal(body)
		assert.NoError(t, err, "init json is failed")

		jsonStr := string(jBytes)
		fmt.Println(jsonStr)

		req, _ := http.NewRequest("PUT", "/message/store/"+ids[index], bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		fmt.Println("body", w.Body)
		var idResponse map[string]string
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &idResponse)

		// get result to confirm
		req, _ = http.NewRequest("GET", "/message/store/"+ids[index], nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		fmt.Println("body", w.Body)
		var response map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		AssertMaps(t, response, wants[index])
	}
}
