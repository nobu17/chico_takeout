package tests

import (
	"crypto/rand"
	"testing"
)

func AssertMaps(t *testing.T, got map[string]interface{}, want map[string]interface{}) {
	for k, v := range want {
		AssertMap(t, got, k, v)
	}
}

func AssertMap(t *testing.T, maps map[string]interface{}, key string, want interface{}) {
	value, exists := maps[key]
	if !exists {
		t.Errorf("failed to get data from map. key=%s", key)
		return
	}
	var got interface{}
	switch a := value.(type) {
	case float64:
		got = int(a)
	case map[string]interface{}:
		val, ok := want.(map[string]interface{})
		if !ok {
			t.Errorf("failed to want:%v", want)
		} else {
			AssertMaps(t, a, val)
			return
		}
	default:
		got = value
	}
	if got != want {
		t.Errorf("got value is not same from want. key=%s got=%s, want=%s", key, value, want)
	}
}

func MakeRandomStr(digit uint32) string {
	// scope of generated string
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// create rand
	b := make([]byte, digit)
	if _, err := rand.Read(b); err != nil {
		panic("unexpected error")
	}
	// select string
	var result string
	for _, v := range b {
		result += string(letters[int(v)%len(letters)])
	}
	return result
}