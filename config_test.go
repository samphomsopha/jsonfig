package jsonfig

import (
	//"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

const defaultJson string = `{"env": "default", 
								"noride": "noway",
								"version": 10.0, 
								"level": 100, 
								"disqus": {
									"api_key": "12030v003", 
									"api_version": 1.0, 
									"user": {
										"name": "Mally", 
										"addresses": {
											"ship" : "1000 Ship Me Lane", 
											"bill" : "1000 Bill Me Lane"
											}
										} 
									} 
								}`

const testJson string = `{"env": "test", "version": 10.20, "disqus": {"api_key": "92838vk92kv", "api_version": 1.0, "user": {"name": "John Lee", "addresses": {"ship" : "1029 Ship Me Lane", "bill" : "1029 Bill Me Lane"} } } }`

func TestNew(t *testing.T) {
	//test creation of Config struct
	c := New()
	if reflect.TypeOf(c) != reflect.TypeOf(&Config{}) {
		t.Error("Expected", reflect.TypeOf(&Config{}), "got", reflect.TypeOf(c))
	}
}

func TestAddValues(t *testing.T) {
	c := New()
	assert.Nil(t, c.AddValues([]byte(testJson)))

}

func TestLoadFile(t *testing.T) {
	c := New()
	assert.Nil(t, c.LoadFile("test.json"))
}

func TestLoadFiles(t *testing.T) {
	c := New()
	loadTheses := []string{"default.json", "test.json"}
	assert.Nil(t, c.LoadFiles(loadTheses))
	files := c.LoadedFiles()
	assert.Equal(t, loadTheses, files)
}

func TestLoadedFiles(t *testing.T) {
	c := New()
	c.LoadFile("default.json")
	c.LoadFile("test.json")
	c.AddValues([]byte(testJson))
	files := c.LoadedFiles()
	testcase := []string{"default.json", "test.json", "setValues"}
	assert.Equal(t, testcase, files)
}

func TestGet(t *testing.T) {
	conf := New()
	err := conf.AddValues([]byte(defaultJson))
	if err != nil {
		t.Error("Expected", nil, "got", err.Error())
	}

	strcases := []struct {
		in, want string
	}{
		{"env", "default"},
		{"noride", "noway"},
		{"disqus.api_key", "12030v003"},
		{"disqus.user.name", "Mally"},
		{"disqus.user.addresses.ship", "1000 Ship Me Lane"},
	}

	for _, c := range strcases {
		got := conf.Get(c.in)
		if got != c.want {
			t.Errorf("Reading Key (%q) ==  Value (%q), want %q", c.in, got, c.want)
		}
	}
}
