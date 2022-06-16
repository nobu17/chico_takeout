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

var businessHoursMemory *domains.BusinessHours

const hoursUrl = "/store/hour"

func SetupHourRouter() *gin.Engine {
	r := gin.Default()
	// hour
	businessHourRepo := memory.NewBusinessHoursMemoryRepository()
	spBusinessHourRepo := memory.NewSpecialBusinessHourMemoryRepository()
	hour := r.Group(hoursUrl)
	{
		businessHourRepo.Reset()
		businessHoursMemory = businessHourRepo.GetMemory()
		useCase := storeUseCase.NewBusinessHoursUseCase(businessHourRepo, spBusinessHourRepo)
		handler := storeHandler.BusinessHoursHandler(useCase)
		hour.GET("/", handler.Get)
		hour.PUT("/:id", handler.Put)
	}
	return r
}

func getAllBusinessHour(t *testing.T, r *gin.Engine) []map[string]interface{} {
	// GET to confirm result
	req, _ := http.NewRequest("GET", "/store/hour/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)
	res, ok := response["schedules"].([]interface{})
	if !ok {
		assert.Fail(t, "failed to cast result")
	}
	results := []map[string]interface{}{}
	for _, ind := range res {
		hours, ok := ind.(map[string]interface{})
		if !ok {
			assert.Fail(t, "failed")
			continue
		}
		results = append(results, hours)
	}
	return results
}

func TestBusinessHoursHandler_GET(t *testing.T) {
	wants := []map[string]interface{}{
		{"name": "morning", "start": "07:00", "end": "09:30", "weekdays": []int{2, 3, 5, 6, 0}},
		{"name": "lunch", "start": "11:30", "end": "15:00", "weekdays": []int{2, 3, 5, 6, 0}},
		{"name": "dinner", "start": "18:00", "end": "21:00", "weekdays": []int{3, 6}},
	}

	r := SetupHourRouter()
	req, _ := http.NewRequest("GET", "/store/hour/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Convert the JSON response to a map
	var response map[string]interface{}
	_ = json.Unmarshal([]byte(w.Body.Bytes()), &response)
	res, ok := response["schedules"].([]interface{})
	if !ok {
		assert.Fail(t, "failed to cast result")
	}
	for index, ind := range res {
		hours, ok := ind.(map[string]interface{})
		if !ok {
			assert.Fail(t, "failed")
			continue
		}
		AssertMaps(t, hours, wants[index])
	}
}

func TestBusinessHoursHandler_PUT(t *testing.T) {
	r := SetupHourRouter()

	type input struct {
		name string
		id   string
		args map[string]interface{}
		want []map[string]interface{}
	}
	inputs := []input{
		{name: "put morning", id: businessHoursMemory.GetSchedules()[0].GetId(),
			args: map[string]interface{}{"name": "morning2", "start": "08:00", "end": "09:00", "weekdays": []int{2, 3, 4}},
			want: []map[string]interface{}{
				{"name": "morning2", "start": "08:00", "end": "09:00", "weekdays": []int{2, 3, 4}},
				{"name": "lunch", "start": "11:30", "end": "15:00", "weekdays": []int{2, 3, 5, 6, 0}},
				{"name": "dinner", "start": "18:00", "end": "21:00", "weekdays": []int{3, 6}},
			}},
		{name: "put lunch", id: businessHoursMemory.GetSchedules()[1].GetId(),
			args: map[string]interface{}{"name": "lunch2", "start": "11:00", "end": "14:30", "weekdays": []int{4}},
			want: []map[string]interface{}{
				{"name": "morning2", "start": "08:00", "end": "09:00", "weekdays": []int{2, 3, 4}},
				{"name": "lunch2", "start": "11:00", "end": "14:30", "weekdays": []int{4}},
				{"name": "dinner", "start": "18:00", "end": "21:00", "weekdays": []int{3, 6}},
			}},
		{name: "put dinner", id: businessHoursMemory.GetSchedules()[2].GetId(),
			args: map[string]interface{}{"name": "dinner2", "start": "17:00", "end": "20:00", "weekdays": []int{6}},
			want: []map[string]interface{}{
				{"name": "morning2", "start": "08:00", "end": "09:00", "weekdays": []int{2, 3, 4}},
				{"name": "lunch2", "start": "11:00", "end": "14:30", "weekdays": []int{4}},
				{"name": "dinner2", "start": "17:00", "end": "20:00", "weekdays": []int{6}},
			}},
		{name: "empty weekend", id: businessHoursMemory.GetSchedules()[2].GetId(),
			args: map[string]interface{}{"name": "dinner2", "start": "17:00", "end": "20:00", "weekdays": []int{}},
			want: []map[string]interface{}{
				{"name": "morning2", "start": "08:00", "end": "09:00", "weekdays": []int{2, 3, 4}},
				{"name": "lunch2", "start": "11:00", "end": "14:30", "weekdays": []int{4}},
				{"name": "dinner2", "start": "17:00", "end": "20:00", "weekdays": []int{}},
			}},
	}

	for _, input := range inputs {
		jBytes, err := json.Marshal(input.args)
		if err != nil {
			assert.Fail(t, "failed to create json", err)
			continue
		}

		req, _ := http.NewRequest("PUT", "/store/hour/"+input.id, bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// GET to confirm result
		results := getAllBusinessHour(t, r)
		for index, result := range results {
			AssertMaps(t, result, input.want[index])
		}
	}
}

func TestBusinessHoursHandler_PUT_BadRequest(t *testing.T) {
	r := SetupHourRouter()

	type input struct {
		name string
		id   string
		args map[string]interface{}
		want []map[string]interface{}
	}

	want := []map[string]interface{}{
		{"name": "morning", "start": "07:00", "end": "09:30", "weekdays": []int{2, 3, 5, 6, 0}},
		{"name": "lunch", "start": "11:30", "end": "15:00", "weekdays": []int{2, 3, 5, 6, 0}},
		{"name": "dinner", "start": "18:00", "end": "21:00", "weekdays": []int{3, 6}},
	}

	inputs := []input{
		{name: "error: empty name", id: businessHoursMemory.GetSchedules()[0].GetId(),
			args: map[string]interface{}{"name": "", "start": "08:00", "end": "09:00", "weekdays": []int{2, 3, 4}},
			want: want,
		},
		{name: "error: start time format", id: businessHoursMemory.GetSchedules()[0].GetId(),
			args: map[string]interface{}{"name": "morning2", "start": "0800", "end": "09:00", "weekdays": []int{2, 3, 4}},
			want: want,
		},
		{name: "error: end time format", id: businessHoursMemory.GetSchedules()[0].GetId(),
			args: map[string]interface{}{"name": "morning2", "start": "08:00", "end": "0900", "weekdays": []int{2, 3, 4}},
			want: want,
		},
		{name: "error: start is greater than end time", id: businessHoursMemory.GetSchedules()[0].GetId(),
			args: map[string]interface{}{"name": "morning2", "start": "08:00", "end": "07:00", "weekdays": []int{2, 3, 4}},
			want: want,
		},
		{name: "error: duplicated weekends", id: businessHoursMemory.GetSchedules()[0].GetId(),
			args: map[string]interface{}{"name": "morning2", "start": "08:00", "end": "09:00", "weekdays": []int{2, 3, 4, 2}},
			want: want,
		},
		{name: "error: overlap time", id: businessHoursMemory.GetSchedules()[0].GetId(),
			args: map[string]interface{}{"name": "morning2", "start": "08:00", "end": "12:00", "weekdays": []int{2, 3, 4}},
			want: want,
		},
		{name: "error: lack name", id: businessHoursMemory.GetSchedules()[0].GetId(),
			args: map[string]interface{}{"start": "08:00", "end": "10:00", "weekdays": []int{2, 3, 4}},
			want: want,
		},
		{name: "error: lack start", id: businessHoursMemory.GetSchedules()[0].GetId(),
			args: map[string]interface{}{"name": "morning2", "end": "10:00", "weekdays": []int{2, 3, 4}},
			want: want,
		},
		{name: "error: lack end", id: businessHoursMemory.GetSchedules()[0].GetId(),
			args: map[string]interface{}{"name": "morning2", "start": "08:00", "weekdays": []int{2, 3, 4}},
			want: want,
		},
		{name: "error: lack weekends", id: businessHoursMemory.GetSchedules()[0].GetId(),
			args: map[string]interface{}{"name": "morning2", "start": "08:00", "end": "09:00"},
			want: want,
		},
	}

	for _, input := range inputs {
		fmt.Println("case:", input.name)
		jBytes, err := json.Marshal(input.args)
		if err != nil {
			assert.Fail(t, "failed to create json", err)
			continue
		}

		req, _ := http.NewRequest("PUT", "/store/hour/"+input.id, bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		// GET to confirm result
		results := getAllBusinessHour(t, r)
		for index, result := range results {
			AssertMaps(t, result, input.want[index])
		}
	}
}

func TestBusinessHoursHandler_PUT_NotFound(t *testing.T) {
	r := SetupHourRouter()

	type input struct {
		name string
		id   string
		args map[string]interface{}
		want []map[string]interface{}
	}

	want := []map[string]interface{}{
		{"name": "morning", "start": "07:00", "end": "09:30", "weekdays": []int{2, 3, 5, 6, 0}},
		{"name": "lunch", "start": "11:30", "end": "15:00", "weekdays": []int{2, 3, 5, 6, 0}},
		{"name": "dinner", "start": "18:00", "end": "21:00", "weekdays": []int{3, 6}},
	}

	inputs := []input{
		{name: "not exists id", id: "12345",
			args: map[string]interface{}{"name": "test", "start": "08:00", "end": "09:00", "weekdays": []int{2, 3, 4}},
			want: want,
		},
	}

	for _, input := range inputs {
		fmt.Println("case:", input.name)
		jBytes, err := json.Marshal(input.args)
		if err != nil {
			assert.Fail(t, "failed to create json", err)
			continue
		}

		req, _ := http.NewRequest("PUT", "/store/hour/"+input.id, bytes.NewBuffer(jBytes))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		// GET to confirm result
		results := getAllBusinessHour(t, r)
		for index, result := range results {
			AssertMaps(t, result, input.want[index])
		}
	}
}
