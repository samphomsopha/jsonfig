package jsonfig

import (
	"encoding/json"
	//"fmt"
	"github.com/spf13/cast"
	"io/ioutil"
	"os"
	"reflect"
	"time"
)

type Config struct {
	configFiles  []string
	configValues map[string]interface{}
	keyDelim     string
}

//load config values from this json file
func (c *Config) LoadFile(fname string) error {
	if _, err := os.Stat(fname); err != nil {
		return err
	} else {
		c.configFiles = append(c.configFiles, fname)
		byt, err := ioutil.ReadFile(fname)

		if err != nil {
			return err
		}

		return c.setValues(byt)
	}

	return nil
}

//load config in order, last in overrides first in values
func (c *Config) LoadFiles(files []string) error {
	for _, filename := range files {
		if err := c.LoadFile(filename); err != nil {
			return err
		}
	}
	return nil
}

//load configuration from json
func (c *Config) AddValues(b []byte) error {
	c.configFiles = append(c.configFiles, "setValues")
	return c.setValues(b)
}

//returns Config.configValues
func (c *Config) RawValues() map[string]interface{} {
	return c.configValues
}

//get Config value by key
func (c *Config) Get(key string) interface{} {
	val := c.get(key)

	switch val.(type) {
	case bool:
		return cast.ToBool(val)
	case string:
		return cast.ToString(val)
	case int64, int32, int16, int8, int:
		return cast.ToInt(val)
	case float64, float32:
		return cast.ToFloat64(val)
	case time.Time:
		return cast.ToTime(val)
	case time.Duration:
		return cast.ToDuration(val)
	case []string:
		return val
	}

	return val
}

func (c *Config) get(key string) interface{} {
	return c.find(key)
}

func (c *Config) GetString(key string) string {
	return cast.ToString(c.get(key))
}

func (c *Config) GetBool(key string) bool {
	return cast.ToBool(c.get(key))
}

func (c *Config) GetInt(key string) int {
	return cast.ToInt(c.get(key))
}

func (c *Config) GetFloat64(key string) float64 {
	return cast.ToFloat64(c.get(key))
}

func (c *Config) GetTime(key string) time.Time {
	return cast.ToTime(c.get(key))
}

func (c *Config) GetDuration(key string) time.Duration {
	return cast.ToDuration(c.get(key))
}

//returns loaded files
func (c *Config) LoadedFiles() []string {
	return c.configFiles
}

func New() *Config {
	c := new(Config)
	c.keyDelim = "."
	c.configValues = make(map[string]interface{})
	return c
}

func (c *Config) setValues(b []byte) error {
	var newdata map[string]interface{}
	if err := json.Unmarshal(b, &newdata); err != nil {
		return err
	}
	c.expandKeysValues(newdata, "")

	return nil
}

func (c *Config) SetValue(k string, v interface{}) {
	c.configValues[k] = v
}

func (c *Config) expandKeysValues(vd map[string]interface{}, kkl string) {
	defMap := make(map[string]interface{})
	kl := ""
	for k, v := range vd {
		if kkl != "" {
			kl = kkl + "." + k
		} else {
			kl = k
		}
		if reflect.TypeOf(v) == reflect.TypeOf(defMap) {
			c.expandKeysValues(cast.ToStringMap(v), kl)
		} else {
			var val interface{} = v
			c.configValues[kl] = val
		}
	}
}

func (c *Config) find(key string) interface{} {
	val := c.configValues[key]
	return val
}
