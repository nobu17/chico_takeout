package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	domains "chico/takeout/domains/store"
	storeHandler "chico/takeout/handlers/store"
	"chico/takeout/infrastructures/memory"
	storeUseCase "chico/takeout/usecase/store"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var spHolidayMemory map[string]*domains.SpecialHoliday

const holidayUrl = "/store/holiday"

func SetupSpecialHolidayRouter() *gin.Engine {
	r := gin.Default()
	businessHourRepo := memory.NewBusinessHoursMemoryRepository()
	businessHoursMemory = businessHourRepo.GetMemory()

	holidayRepo := memory.NewSpecialHolidayMemoryRepository()
	holidayRepo.Reset()
	spHolidayMemory = holidayRepo.GetMemory()
	holiday := r.Group(holidayUrl)
	{
		useCase := storeUseCase.NewSpecialHolidayUseCase(holidayRepo)
		handler := storeHandler.NewSpecialHolidayHandler(useCase)
		holiday.GET("/:id", handler.Get)
		holiday.GET("/", handler.GetAll)
		holiday.POST("/", handler.Post)
		holiday.PUT("/:id", handler.Put)
		holiday.DELETE("/:id", handler.Delete)
	}
	return r
}

func TestSpecialHolidayHandler_GETALL(t *testing.T) {

	wants := []map[string]interface{}{
		{"name": "おやすみ１", "start": "2022/05/06", "end": "2022/06/03"},
		{"name": "おやすみ2", "start": "2022/07/06", "end": "2022/08/01"},
	}

	r := SetupSpecialHolidayRouter()
	req, _ := http.NewRequest("GET", holidayUrl+"/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Convert the JSON response to a map
	var response []map[string]interface{}
	_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)
	fmt.Println("res", w.Body)

	for index, item := range response {
		AssertMaps(t, item, wants[index])
	}
}

func TestSpecialHolidayHandler_GET(t *testing.T) {

	wants := []map[string]interface{}{
		{"name": "おやすみ１", "start": "2022/05/06", "end": "2022/06/03"},
		{"name": "おやすみ2", "start": "2022/07/06", "end": "2022/08/01"},
	}

	r := SetupSpecialHolidayRouter()

	index := 0
	for id := range spHolidayMemory {
		req, _ := http.NewRequest("GET", holidayUrl+"/"+id, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Convert the JSON response to a map
		var response map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)
		fmt.Println("res", w.Body)

		AssertMaps(t, response, wants[index])
		index++
	}
}

func TestSpecialHolidayHandler_GET_NotFound(t *testing.T) {

	r := SetupSpecialHolidayRouter()
	req, _ := http.NewRequest("GET", holidayUrl+"/1234", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSpecialHolidayHandler_POST_CREATE(t *testing.T) {
	r := SetupSpecialHolidayRouter()

	bodies := []map[string]interface{}{
		{"name": "おやすみ３", "start": "2022/09/01", "end": "2022/09/01"},
		{"name": "おやすみ４", "start": "2022/09/11", "end": "2022/09/12"},
	}
	wants := []map[string]interface{}{
		{"name": "おやすみ３", "start": "2022/09/01", "end": "2022/09/01"},
		{"name": "おやすみ４", "start": "2022/09/11", "end": "2022/09/12"},
	}
	for index, body := range bodies {
		jBytes, err := json.Marshal(body)
		assert.NoError(t, err, "init json is failed")

		req, _ := http.NewRequest("POST", holidayUrl+"/", bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var idResponse map[string]string
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &idResponse)

		id := idResponse["id"]
		assert.NotEmpty(t, id, "response id should not be empty.")

		// confirm result from result id
		getReq, _ := http.NewRequest("GET", holidayUrl+"/"+id, nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, getReq)

		assert.Equal(t, http.StatusOK, w.Code)
		fmt.Println("response", w.Body)

		var response map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		AssertMaps(t, response, wants[index])
	}
}

type specilaHolidayItemErrorData struct {
	name string
	args map[string]interface{}
	want int
}

func getspecilaHolidayItemCommonErrorData() []specilaHolidayItemErrorData {
	var items = []specilaHolidayItemErrorData{
		{name: "lack name", args: map[string]interface{}{
			"start": "2022/11/01", "end": "2022/11/04",
		}, want: 2},
		{name: "empty name", args: map[string]interface{}{
			"name": "", "start": "2022/11/01", "end": "2022/11/04",
		}, want: 2},
		{name: "over limit name(21)", args: map[string]interface{}{
			"name": "123456789012345678901", "start": "2022/11/01", "end": "2022/11/04",
		}, want: 2},
		{name: "lack start", args: map[string]interface{}{
			"name": "1234", "end": "2022/11/04",
		}, want: 2},
		{name: "incorrect format start", args: map[string]interface{}{
			"name": "1234", "start": "20221101", "end": "2022/11/04",
		}, want: 2},
		{name: "lack end", args: map[string]interface{}{
			"name": "1234", "start": "2022/11/04",
		}, want: 2},
		{name: "incorrect end start", args: map[string]interface{}{
			"name": "1234", "start": "2022/11/01", "end": "20221104",
		}, want: 2},
		{name: "start > end", args: map[string]interface{}{
			"name": "1234", "start": "2022/11/05", "end": "2022/11/04",
		}, want: 2},
	}
	return items
}

func getSpecilaHolidayItemPostErrorData() []specilaHolidayItemErrorData {
	var commonError = getspecilaHolidayItemCommonErrorData()
	var items = []specilaHolidayItemErrorData{
		{name: "overlap date(1)", args: map[string]interface{}{
			"name": "1234", "start": "2022/05/01", "end": "2022/05/07",
		}, want: 2},
		{name: "overlap date(2)", args: map[string]interface{}{
			"name": "1234", "start": "2022/05/01", "end": "2022/06/07",
		}, want: 2},
		{name: "overlap date(3)", args: map[string]interface{}{
			"name": "1234", "start": "2022/05/06", "end": "2022/05/17",
		}, want: 2},
		{name: "overlap date(4)", args: map[string]interface{}{
			"name": "1234", "start": "2022/05/06", "end": "2022/06/17",
		}, want: 2},
	}
	commonError = append(commonError, items...)
	return commonError
}

func getSpecilaHolidayItemPutErrorData() []specilaHolidayItemErrorData {
	var commonError = getspecilaHolidayItemCommonErrorData()
	var items = []specilaHolidayItemErrorData{
		{name: "overlap date(1)", args: map[string]interface{}{
			"name": "1234", "start": "2022/07/10", "end": "2022/10/07",
		}, want: 2},
	}
	commonError = append(commonError, items...)
	return commonError
}

func TestSpecialHolidayHandler_POST_BadRequest(t *testing.T) {
	r := SetupSpecialHolidayRouter()

	inputs := getSpecilaHolidayItemPostErrorData()
	assert.NotEqual(t, 0, len(inputs), "input data is empty")

	for _, tt := range inputs {
		fmt.Println("case:", tt.name)
		jBytes, err := json.Marshal(tt.args)
		assert.NoError(t, err, "init json is failed")

		req, _ := http.NewRequest("POST", holidayUrl+"/", bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		fmt.Println("body", w.Body)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// confirm stock is not added
		req, _ = http.NewRequest("GET", holidayUrl+"/", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "failed to get result to confirm.")

		var response []map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		assert.Equal(t, tt.want, len(spHolidayMemory), "response length is incorrect")
	}
}

func TestSpecialHolidayHandler_PUT(t *testing.T) {
	r := SetupSpecialHolidayRouter()

	spIds := []string{}
	for id := range spHolidayMemory {
		spIds = append(spIds, id)
	}

	bodies := []map[string]interface{}{
		{"name": "おやすみ３", "start": "2022/09/01", "end": "2022/09/01"},
		{"name": "おやすみ４", "start": "2022/09/11", "end": "2022/09/12"},
	}
	wants := []map[string]interface{}{
		{"name": "おやすみ３", "start": "2022/09/01", "end": "2022/09/01"},
		{"name": "おやすみ４", "start": "2022/09/11", "end": "2022/09/12"},
	}
	for index, body := range bodies {
		jBytes, err := json.Marshal(body)
		assert.NoError(t, err, "init json is failed")

		req, _ := http.NewRequest("PUT", holidayUrl+"/"+spIds[index], bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		fmt.Println("response", w.Body)

		assert.Equal(t, http.StatusOK, w.Code)

		var idResponse map[string]string
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &idResponse)

		// confirm result from result id
		getReq, _ := http.NewRequest("GET", holidayUrl+"/"+spIds[index], nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, getReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		AssertMaps(t, response, wants[index])
	}
}

func TestSpecialHolidayHandler_PUT_NotFound(t *testing.T) {
	r := SetupSpecialHolidayRouter()

	bodies := []map[string]interface{}{
		{"name": "おやすみ３", "start": "2022/09/01", "end": "2022/09/01"},
	}
	jBytes, err := json.Marshal(bodies[0])
	assert.NoError(t, err, "init json is failed")

	req, _ := http.NewRequest("PUT", holidayUrl+"/1234", bytes.NewBuffer(jBytes))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	fmt.Println("response", w.Body)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSpecialHolidayHandler_PUT_BadRequest(t *testing.T) {
	r := SetupSpecialHolidayRouter()

	spIds := []string{}
	for id := range spHolidayMemory {
		spIds = append(spIds, id)
	}

	inputs := getSpecilaHolidayItemPutErrorData()
	assert.NotEqual(t, 0, len(inputs), "input data is empty")

	for _, tt := range inputs {
		fmt.Println("case:", tt.name)
		jBytes, err := json.Marshal(tt.args)
		assert.NoError(t, err, "init json is failed")

		req, _ := http.NewRequest("PUT", holidayUrl+"/"+spIds[0], bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		fmt.Println("body", w.Body)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// confirm stock is not added
		req, _ = http.NewRequest("GET", holidayUrl+"/", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "failed to get result to confirm.")

		var response []map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		assert.Equal(t, tt.want, len(spHolidayMemory), "response length is incorrect")
	}
}

func TestSpecialHolidayHandler_DELETE(t *testing.T) {
	r := SetupSpecialHolidayRouter()

	spIds := []string{}
	for id := range spHolidayMemory {
		spIds = append(spIds, id)
	}

	req, _ := http.NewRequest("DELETE", holidayUrl+"/"+spIds[0], nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// confirm result
	getReq, _ := http.NewRequest("GET", holidayUrl+"/"+spIds[0], nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, getReq)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSpecialHolidayHandler_DELETE_NotFound(t *testing.T) {
	r := SetupSpecialHolidayRouter()

	req, _ := http.NewRequest("DELETE", holidayUrl+"/1234", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
