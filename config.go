package jsonfig

import (
	"encoding/json"
	//"fmt"
	"github.com/spf13/cast"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

type Config struct {
	configFiles  []string
	configValues map[string]interface{}
	keyDelim     string
}

/*
load config values this json file
*/
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

/*
load config in order, last in overrides first in values
*/
func (c *Config) LoadFiles(files []string) error {
	for _, filename := range files {
		if err := c.LoadFile(filename); err != nil {
			return err
		}
	}
	return nil
}

/*
append json object to config
*/
func (c *Config) AddValues(b []byte) error {
	c.configFiles = append(c.configFiles, "setValues")
	return c.setValues(b)
}

/*
return Config.configValues
*/
func (c *Config) RawValues() map[string]interface{} {
	return c.configValues
}

/*
get Config value by key
*/
func (c *Config) Get(key string) interface{} {

	return c.find(key)
}

/*
returns loaded files
*/
func (c *Config) LoadedFiles() []string {
	return c.configFiles
}

func New() *Config {
	c := new(Config)
	c.keyDelim = "."
	return c
}

func (c *Config) setValues(b []byte) error {
	if err := json.Unmarshal(b, &c.configValues); err != nil {
		return err
	}

	return nil
}

func (c *Config) find(key string) interface{} {
	path := strings.Split(key, c.keyDelim)
	val := c.configValues[key]

	if val == nil {
		source := c.find(path[0])
		if source == nil {
			return nil
		}

		if reflect.TypeOf(source).Kind() == reflect.Map {
			val = c.searchMap(cast.ToStringMap(source), path[1:])
		}
	}

	return val
}

func (c *Config) searchMap(smap map[string]interface{}, path []string) interface{} {
	key := path[0]
	val := smap[key]
	if len(path) != 1 {
		return c.searchMap(cast.ToStringMap(val), path[1:])
	}
	return val
}
