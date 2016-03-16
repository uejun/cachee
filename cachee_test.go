package cachee

import (
	"testing"
	"time"
)

func clear() {
	cache = make(map[string]interface{})
}

func TestSetGet(t *testing.T) {
	var tests = []struct {
		k      string
		v      interface{}
		expire time.Duration
	}{
		{k: "num", v: 1, expire: time.Second},
		{k: "str", v: "a", expire: time.Second},
	}

	for _, test := range tests {
		Set(test.k, test.v, test.expire)
		v, ok := Get(test.k)
		if !ok {
			t.Errorf("No such key %s", test.k)
		}
		switch tt := v.(type) {
		case int:
			if tt != test.v {
				t.Errorf("Got value != %d", test.v)
			}
		case string:
			if tt != test.v {
				t.Errorf("Got value != %s", test.v)
			}
		default:
			t.Error("Not ready type")
		}

		time.Sleep(2 * time.Second)

		if _, ok := Get(test.k); ok {
			t.Errorf("Expire doesn't work for %s", test.k)
		}
	}

	clear()

}

func TestGetWithoutSet(t *testing.T) {
	if v, ok := Get("NotSet"); v != nil || ok {
		t.Error("The value for not setted key found.")
	}
}

func TestDelete(t *testing.T) {
	cache["key"] = "value"
	Delete("key")
	if v, ok := cache["key"]; v == "value" || ok {
		t.Error("Delete(key) failed.")
	}

	clear()
}

func TestGetIfNotSet(t *testing.T) {
	if _, ok := Get("NotYetSet"); ok {
		t.Fatal("The key 'NotYetSet' is already used.")
	}
	if v, ok := GetIfNotSet("NotYetSet", "val", time.Second); ok {
		t.Error("GetIfNotSet('NotYetSet') returns true, expected false")
	} else if v != "val" {
		t.Errorf("GetIfNotSet('NotYetSet') returns %s, expected 'val'", v)
	}

	Set("AlreadySet", "val", time.Second)
	if v, ok := GetIfNotSet("AlreadySet", "val", time.Second); !ok {
		t.Error("GetIfNotSet('AlreadySet') returns false, expected true")
	} else if v != "val" {
		t.Errorf("GetIfNotSet('AlreadySet') returns %s, expected 'val'", v)
	}

	clear()
}

func TestKeys(t *testing.T) {
	tests := map[string]string{"a": "A", "b": "B", "c": "C"}
	for k, v := range tests {
		Set(k, v, 2*time.Second)
	}
	got := Keys()
	if len(got) != len(tests) {
		t.Error("len(%v) does not match expected len %d.", got, len(tests))
	}
	for _, k := range got {
		if _, ok := tests[k]; !ok {
			t.Errorf("The key %s is not expected.", k)
		}
	}
	clear()
}

func TestValues(t *testing.T) {
	tests := map[string]string{"a": "A", "b": "B", "c": "C"}
	for k, v := range tests {
		Set(k, v, time.Second)
	}
	got := Values()
	if len(got) != len(tests) {
		t.Error("len(%v) does not match expected len %d.", got, len(tests))
	}
	for _, val := range got {
		switch val {
		case "A", "B", "C":
		default:
			t.Error("Values() has a not expected value")
		}
	}
	clear()
}
