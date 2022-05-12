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

var spBusHourMemory map[string]*domains.SpecialBusinessHour

const spBusUrl = "/store/special_hour"

func SetupSpecialBusinessHourRouter() *gin.Engine {
	r := gin.Default()
	businessHoursRepo := memory.NewBusinessHoursMemoryRepository()
	businessHoursMemory = businessHoursRepo.GetMemory()

	spBusinessHourRepo := memory.NewSpecialBusinessHourMemoryRepository()
	spBusinessHourRepo.Reset()
	spBusHourMemory = spBusinessHourRepo.GetMemory()
	specialHour := r.Group("/store/special_hour")
	{
		useCase := storeUseCase.NewSpecialBusinessHoursUseCase(businessHoursRepo, spBusinessHourRepo)
		handler := storeHandler.NewSpecialBusinessHourHandler(useCase)
		specialHour.GET("/:id", handler.Get)
		specialHour.GET("/", handler.GetAll)
		specialHour.POST("/", handler.Post)
		specialHour.PUT("/:id", handler.Put)
		specialHour.DELETE("/:id", handler.Delete)
	}
	return r
}

func TestSpecialBusinessHourHandler_GETALL(t *testing.T) {
	r := SetupSpecialBusinessHourRouter()

	schIds := []string{}
	for _, schId := range businessHoursMemory.GetSchedules() {
		schIds = append(schIds, schId.GetId())
	}

	wants := []map[string]interface{}{
		{"name": "特別日程1", "date": "2022/05/06", "start": "08:00", "end": "12:00", "businessHourId": schIds[0]},
		{"name": "特別日程2", "date": "2022/05/08", "start": "11:00", "end": "14:00", "businessHourId": schIds[1]},
	}

	req, _ := http.NewRequest("GET", spBusUrl+"/", nil)
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

func TestSpecialBusinessHourHandler_GET(t *testing.T) {
	r := SetupSpecialBusinessHourRouter()

	schIds := []string{}
	for _, schId := range businessHoursMemory.GetSchedules() {
		schIds = append(schIds, schId.GetId())
	}

	wants := []map[string]interface{}{
		{"name": "特別日程1", "date": "2022/05/06", "start": "08:00", "end": "12:00", "businessHourId": schIds[0]},
		{"name": "特別日程2", "date": "2022/05/08", "start": "11:00", "end": "14:00", "businessHourId": schIds[1]},
	}

	for id := range spBusHourMemory {
		req, _ := http.NewRequest("GET", spBusUrl+"/"+id, nil)
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
}

func TestSpecialBusinessHourHandler_GET_NotFound(t *testing.T) {

	r := SetupSpecialBusinessHourRouter()
	req, _ := http.NewRequest("GET", spBusUrl+"/1234", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSpecialBusinessHourHandler_POST_Create(t *testing.T) {
	r := SetupSpecialBusinessHourRouter()

	schIds := []string{}
	for _, schId := range businessHoursMemory.GetSchedules() {
		schIds = append(schIds, schId.GetId())
	}

	bodies := []map[string]interface{}{
		{"name": "特別日程3", "date": "2023/07/06", "start": "10:00", "end": "12:00", "businessHourId": schIds[1]},
		{"name": "特別日程4", "date": "2023/07/08", "start": "19:00", "end": "22:00", "businessHourId": schIds[2]},
	}
	wants := []map[string]interface{}{
		{"name": "特別日程3", "date": "2023/07/06", "start": "10:00", "end": "12:00", "businessHourId": schIds[1]},
		{"name": "特別日程4", "date": "2023/07/08", "start": "19:00", "end": "22:00", "businessHourId": schIds[2]},
	}
	for index, body := range bodies {
		jBytes, err := json.Marshal(body)
		assert.NoError(t, err, "init json is failed")

		req, _ := http.NewRequest("POST", spBusUrl+"/", bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var idResponse map[string]string
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &idResponse)

		fmt.Println("idResponse", w.Body)
		id := idResponse["id"]
		assert.NotEmpty(t, id, "response id should not be empty.")

		// confirm result from result id
		getReq, _ := http.NewRequest("GET", spBusUrl+"/"+id, nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, getReq)

		assert.Equal(t, http.StatusOK, w.Code)
		fmt.Println("response", w.Body)

		var response map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		AssertMaps(t, response, wants[index])
	}
}

type specilaBusinessHourErrorData struct {
	name string
	args map[string]interface{}
	want int
}

func getSpecialBusinessHourCommonErrorData(busHourIds []string) []specilaBusinessHourErrorData {
	var items = []specilaBusinessHourErrorData{
		{name: "lack name", args: map[string]interface{}{
			"date": "2024/11/12", "businessHourId": busHourIds[0], "start": "9:00", "end": "10:30",
		}, want: 2},
		{name: "empty name", args: map[string]interface{}{
			"name": "", "date": "2024/11/12", "businessHourId": busHourIds[0], "start": "9:00", "end": "10:30",
		}, want: 2},
		{name: "over limit name(31)", args: map[string]interface{}{
			"name": "1234567890123456789012345678901", "date": "2024/11/12", "businessHourId": busHourIds[0], "start": "9:00", "end": "10:30",
		}, want: 2},
		{name: "lack date", args: map[string]interface{}{
			"name": "1234", "businessHourId": busHourIds[0], "start": "9:00", "end": "10:30",
		}, want: 2},
		{name: "incorrect format date", args: map[string]interface{}{
			"name": "1234", "date": "20241112", "businessHourId": busHourIds[0], "start": "9:00", "end": "10:30",
		}, want: 2},
		{name: "lack start", args: map[string]interface{}{
			"name": "1234", "date": "2024/11/12", "businessHourId": busHourIds[0], "end": "10:30",
		}, want: 2},
		{name: "incorrect format start", args: map[string]interface{}{
			"name": "1234", "date": "2024/11/12", "businessHourId": busHourIds[0], "start": "900", "end": "10:30",
		}, want: 2},
		{name: "lack end", args: map[string]interface{}{
			"name": "1234", "date": "2024/11/12", "businessHourId": busHourIds[0], "start": "10:30",
		}, want: 2},
		{name: "incorrect end start", args: map[string]interface{}{
			"name": "1234", "date": "2024/11/12", "businessHourId": busHourIds[0], "start": "09:00", "end": "1030",
		}, want: 2},
		{name: "start > end", args: map[string]interface{}{
			"name": "1234", "date": "2024/11/12", "businessHourId": busHourIds[0], "start": "09:00", "end": "08:30",
		}, want: 2},
		{name: "start > end(+59)", args: map[string]interface{}{
			"name": "1234", "date": "2024/11/12", "businessHourId": busHourIds[0], "start": "09:00", "end": "09:59",
		}, want: 2},
		{name: "lack businessHourId", args: map[string]interface{}{
			"name": "1234", "date": "2024/11/12", "start": "9:00", "end": "10:30",
		}, want: 2},
		{name: "empty businessHourId", args: map[string]interface{}{
			"name": "1234", "date": "2024/11/12", "businessHourId": "", "start": "9:00", "end": "10:30",
		}, want: 2},
	}
	return items
}

func getSpecialBusinessHourPostErrorData(busHourIds []string) []specilaBusinessHourErrorData {
	var commonError = getSpecialBusinessHourCommonErrorData(busHourIds)
	var items = []specilaBusinessHourErrorData{
		{name: "overlap date and hourId", args: map[string]interface{}{
			"name": "1234", "date": "2022/05/06", "businessHourId": busHourIds[0], "start": "9:10", "end": "10:30",
		}, want: 2},
		{name: "overlap time with other hourId", args: map[string]interface{}{
			"name": "1234", "date": "2022/05/06", "businessHourId": busHourIds[1], "start": "10:30", "end": "14:30",
		}, want: 2},
	}
	commonError = append(commonError, items...)
	return commonError
}

func TestSpecialBusinessHourHandler_POST_BadRequest(t *testing.T) {
	r := SetupSpecialBusinessHourRouter()

	schIds := []string{}
	for _, schId := range businessHoursMemory.GetSchedules() {
		schIds = append(schIds, schId.GetId())
	}

	inputs := getSpecialBusinessHourPostErrorData(schIds)
	assert.NotEqual(t, 0, len(inputs), "input data is empty")

	for _, tt := range inputs {
		fmt.Println("case:", tt.name)
		jBytes, err := json.Marshal(tt.args)
		assert.NoError(t, err, "init json is failed")

		req, _ := http.NewRequest("POST", spBusUrl+"/", bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		fmt.Println("body", w.Body)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// confirm stock is not added
		req, _ = http.NewRequest("GET", spBusUrl+"/", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "failed to get result to confirm.")

		var response []map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		assert.Equal(t, tt.want, len(spBusHourMemory), "response length is incorrect")
	}
}

func TestSpecialBusinessHourHandler_Put(t *testing.T) {
	r := SetupSpecialBusinessHourRouter()

	busIds := []string{}
	for busId := range spBusHourMemory {
		busIds = append(busIds, busId)
	}

	schIds := []string{}
	for _, schId := range businessHoursMemory.GetSchedules() {
		schIds = append(schIds, schId.GetId())
	}

	bodies := []map[string]interface{}{
		{"name": "特別日程3", "date": "2023/07/06", "start": "10:00", "end": "12:00", "businessHourId": schIds[1]},
		{"name": "特別日程4", "date": "2023/07/08", "start": "19:00", "end": "22:00", "businessHourId": schIds[2]},
	}
	wants := []map[string]interface{}{
		{"name": "特別日程3", "date": "2023/07/06", "start": "10:00", "end": "12:00", "businessHourId": schIds[1]},
		{"name": "特別日程4", "date": "2023/07/08", "start": "19:00", "end": "22:00", "businessHourId": schIds[2]},
	}
	for index, body := range bodies {
		jBytes, err := json.Marshal(body)
		assert.NoError(t, err, "init json is failed")

		req, _ := http.NewRequest("PUT", spBusUrl+"/"+busIds[0], bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		fmt.Println("body", w.Body)
		assert.Equal(t, http.StatusOK, w.Code)

		var idResponse map[string]string
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &idResponse)

		// confirm result from result id
		getReq, _ := http.NewRequest("GET", spBusUrl+"/"+busIds[0], nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, getReq)

		assert.Equal(t, http.StatusOK, w.Code)
		fmt.Println("response", w.Body)

		var response map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)

		AssertMaps(t, response, wants[index])
	}
}

func getSpecialBusinessHourPutErrorData(busHourIds []string) []specilaBusinessHourErrorData {
	var commonError = getSpecialBusinessHourCommonErrorData(busHourIds)
	var items = []specilaBusinessHourErrorData{
		{name: "overlap date and hourId", args: map[string]interface{}{
			"name": "1234", "date": "2022/05/08", "businessHourId": busHourIds[1], "start": "9:10", "end": "10:30",
		}, want: 2},
		{name: "overlap time with other hourId", args: map[string]interface{}{
			"name": "1234", "date": "2022/05/08", "businessHourId": busHourIds[0], "start": "10:30", "end": "14:30",
		}, want: 2},
	}
	commonError = append(commonError, items...)
	return commonError
}

func TestSpecialBusinessHourHandler_PUT_BadRequest(t *testing.T) {
	r := SetupSpecialBusinessHourRouter()

	busIds := []string{}
	for busId := range spBusHourMemory {
		busIds = append(busIds, busId)
	}

	schIds := []string{}
	for _, schId := range businessHoursMemory.GetSchedules() {
		schIds = append(schIds, schId.GetId())
	}

	wants := []map[string]interface{}{
		{"name": "特別日程1", "date": "2022/05/06", "start": "08:00", "end": "12:00", "businessHourId": schIds[0]},
		{"name": "特別日程2", "date": "2022/05/08", "start": "11:00", "end": "14:00", "businessHourId": schIds[1]},
	}

	inputs := getSpecialBusinessHourPutErrorData(schIds)
	assert.NotEqual(t, 0, len(inputs), "input data is empty")

	for _, tt := range inputs {
		fmt.Println("case:", tt.name)
		jBytes, err := json.Marshal(tt.args)
		assert.NoError(t, err, "init json is failed")

		req, _ := http.NewRequest("PUT", spBusUrl+"/"+busIds[0], bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		// confirm stock is not added
		req, _ = http.NewRequest("GET", spBusUrl+"/"+busIds[0], nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "failed to get result to confirm.")

		var response map[string]interface{}
		_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)
		fmt.Println("response", w.Body)

		AssertMaps(t, response, wants[0])
	}
}

func TestSpecialBusinessHourHandler_PUT_NotFound(t *testing.T) {
	r := SetupSpecialBusinessHourRouter()

	schIds := []string{}
	for _, schId := range businessHoursMemory.GetSchedules() {
		schIds = append(schIds, schId.GetId())
	}

	body := []map[string]interface{}{
		{"name": "特別日程1", "date": "2022/05/06", "start": "08:00", "end": "12:00", "businessHourId": schIds[0]},
	}

	jBytes, err := json.Marshal(body[0])
	assert.NoError(t, err, "init json is failed")

	req, _ := http.NewRequest("PUT", spBusUrl+"/1234", bytes.NewBuffer(jBytes))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	fmt.Println("body", w.Body)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSpecialBusinessHourHandler_DELETE(t *testing.T) {
	r := SetupSpecialBusinessHourRouter()

	busIds := []string{}
	for busId := range spBusHourMemory {
		busIds = append(busIds, busId)
	}

	req, _ := http.NewRequest("DELETE", spBusUrl+"/"+busIds[0], nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// confirm result
	getReq, _ := http.NewRequest("GET", spBusUrl+"/"+busIds[0], nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, getReq)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSpecialBusinessHourHandler_DELETE_NotFound(t *testing.T) {
	r := SetupSpecialBusinessHourRouter()

	req, _ := http.NewRequest("DELETE", spBusUrl+"/1234", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
