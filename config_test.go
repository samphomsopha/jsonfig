package jsonfig

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

const defaultJson string = `{"env": "default", 
								"noride": "noway",
								"version": 10.0, 
								"level": 100,
								"build_time": "2012-11-01T22:08:41+00:00",
								"run_duration": "300ms", 
								"disqus": {
									"api_key": "12030v003", 
									"api_version": 1.0, 
									"user": {
										"name": "Mally",
										"city": "Herndon",
										"addresses": {
											"ship" : "1000 Ship Me Lane", 
											"bill" : "1000 Bill Me Lane" } } } }`

const testJson string = `{"env": "test", 
								"version": 10.20, 
								"disqus": {
									"api_key": "92838vk92kv", 
									"api_version": 1.0, 
									"user": {
										"name": "John Lee", 
										"addresses": {
											"ship" : "1029 Ship Me Lane", 
											"bill" : "1029 Bill Me Lane"} } } }`

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
	conf := loadConf()

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
		assert.Equal(t, c.want, got)
	}
}

func TestGetString(t *testing.T) {
	conf := loadConf()
	exp := "default"
	val := conf.GetString("env")
	assert.Equal(t, exp, val)
	if reflect.TypeOf(val) != reflect.TypeOf(exp) {
		t.Error("Expected", reflect.TypeOf(exp), "got", reflect.TypeOf(val))
	}
}

func TestGetInt(t *testing.T) {
	conf := loadConf()
	exp := 100
	val := conf.GetInt("level")
	assert.Equal(t, exp, val)
	if reflect.TypeOf(val) != reflect.TypeOf(exp) {
		t.Error("Expected", reflect.TypeOf(exp), "got", reflect.TypeOf(val))
	}
}

func TestGetFloat64(t *testing.T) {
	conf := loadConf()
	exp := 10.0
	val := conf.GetFloat64("version")
	assert.Equal(t, exp, val)
	if reflect.TypeOf(val) != reflect.TypeOf(exp) {
		t.Error("Expected", reflect.TypeOf(exp), "got", reflect.TypeOf(val))
	}
}

func TestGetTime(t *testing.T) {
	conf := loadConf()
	exp, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	val := conf.GetTime("build_time")
	assert.Equal(t, exp, val)
	if reflect.TypeOf(val) != reflect.TypeOf(exp) {
		t.Error("Expected", reflect.TypeOf(exp), "got", reflect.TypeOf(val))
	}
}

func TestGetDuration(t *testing.T) {
	conf := loadConf()
	exp, _ := time.ParseDuration("300ms")
	val := conf.GetDuration("run_duration")
	assert.Equal(t, exp, val)
	if reflect.TypeOf(val) != reflect.TypeOf(exp) {
		t.Error("Expected", reflect.TypeOf(exp), "got", reflect.TypeOf(val))
	}
}

func TestSetValue(t *testing.T) {
	conf := loadConf()
	conf.SetValue("env", "production")
	conf.SetValue("disqus.user.name", "Alex Morgan")
	conf.SetValue("nokey", "novalue")
	conf.SetValue("nokey.noline", "noline")
	assert.Equal(t, "production", conf.GetString("env"))
	assert.Equal(t, "Alex Morgan", conf.GetString("disqus.user.name"))
	assert.Equal(t, "novalue", conf.GetString("nokey"))
	assert.Equal(t, "noline", conf.GetString("nokey.noline"))
}

func TestNonExistKey(t *testing.T) {
	conf := loadConf()
	assert.Nil(t, conf.Get("nonexist"))
}

func TestOverrides(t *testing.T) {
	conf := loadConf()
	conf.AddValues([]byte(testJson))
	strcases := []struct {
		in, want string
	}{
		{"env", "test"},
		{"noride", "noway"},
		{"disqus.api_key", "92838vk92kv"},
		{"disqus.user.name", "John Lee"},
		{"disqus.user.city", "Herndon"},
		{"disqus.user.addresses.ship", "1029 Ship Me Lane"},
	}

	for _, c := range strcases {
		got := conf.Get(c.in)
		assert.Equal(t, c.want, got)
	}
}

func BenchmarkAddValues(b *testing.B) {
	for n := 0; n < b.N; n++ {
		c := New()
		c.AddValues([]byte(defaultJson))
	}
}

func BenchmarkGetSimple(b *testing.B) {
	c := loadConf()
	for n := 0; n < b.N; n++ {
		c.Get("env")
	}
}

func BenchmarkGetLevel(b *testing.B) {
	c := loadConf()
	for n := 0; n < b.N; n++ {
		c.Get("disqus.user.name")
	}
}

func loadConf() *Config {
	c := New()
	c.AddValues([]byte(defaultJson))

	return c
}
