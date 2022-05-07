package tests

import "testing"

func AssertMaps(t *testing.T, got map[string]interface{}, want map[string]interface{}) {
	for k, v := range want {
		AssertMap(t, got, k, v)
	}
}

func AssertMap(t *testing.T, maps map[string]interface{}, key string, want interface{}) {
	value, exists := maps[key]
	if !exists {
		t.Errorf("failed to get data from map. key=%s", key)
	}
	var got interface{}
	switch a := value.(type) {
	case float64:
		got = int(a)
	default:
		got = value
	}
	if got != want {
		t.Errorf("got value is not same from want. key=%s got=%s, want=%s", key, value, want)
	}
}