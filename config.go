package wenex

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"strings"
	"sync"
)

// DefaultConfig returns default configuration options:
//  server.http.listen:   ":http"
//  server.timeout.read:  "30s"
//  server.timeout.write: "30s"
//  server.timeout.idle:  "30s"
//  logger.defaultName:   "wenex"
//  logger.namePrefix:    "log/"
//  logger.usePrefix:     "[!] "
//  logger.useFlag:       log.LstdFlags
func DefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"server.http.listen":   ":http",
		"server.timeout.read":  "30s",
		"server.timeout.write": "30s",
		"server.timeout.idle":  "30s",
		"logger.defaultName":   "wenex",
		"logger.namePrefix":    "log/",
		"logger.usePrefix":     "[!] ",
		"logger.useFlag":       log.LstdFlags,
	}
}

func NewConfig(name string, defaultConfig map[string]interface{}) (*Config, error) {
	if path := path.Dir(name); path != "" {
		if err := os.MkdirAll(path, 0755); err != nil {
			return nil, err
		}
	}

	file, err := os.OpenFile(name+".conf", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		if os.IsPermission(err) {
			if file, err = os.OpenFile(name+".conf", os.O_RDONLY, 0); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if fi.Size() == 0 {
		if _, err = file.WriteString("{}"); err != nil {
			return nil, err
		}
	}

	config := Config{
		file:    file,
		decoder: json.NewDecoder(file),
		encoder: json.NewEncoder(file),
	}

	config.encoder.SetIndent("", "    ")

	if err = config.Load(); err != nil {
		return nil, err
	}

	if err = config.loadDefault(defaultConfig); err != nil {
		return nil, err
	}

	return &config, nil
}

// Config struct
type Config struct {
	mutex   sync.Mutex
	file    *os.File
	decoder *json.Decoder
	encoder *json.Encoder
	data    map[string]interface{}
}

func (c *Config) loadDefault(defaultConfig map[string]interface{}) error {
	if defaultConfig != nil {
		var needSave bool

		for key, value := range defaultConfig {
			if c.Get(key) == nil {
				c.Set(key, value)
				needSave = true
			}
		}

		if needSave {
			if err := c.Save(); err != nil {
				return err
			}
		}
	}

	return nil
}

// Load loads configuration from file.
func (c *Config) Load() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, err := c.file.Seek(0, 0); err != nil {
		return err
	}

	return c.decoder.Decode(&c.data)
}

// Save saves the current configuration to file.
func (c *Config) Save() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, err := c.file.Seek(0, 0); err != nil {
		return err
	}

	if err := c.file.Truncate(0); err != nil {
		return err
	}

	return c.encoder.Encode(c.data)
}

// Set sets the value to config.
// The key washes to be separated by a dot symbol.
// For example:
//  wnx.Config.Set("key1", "value1")
//  wnx.Config.Set("key2.key3", "value2")
func (c *Config) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	data := c.data
	path := strings.Split(key, ".")
	var i int

	for i < len(path)-1 {
		tmp, ok := data[path[i]]
		if !ok {
			tmp = make(map[string]interface{})
			data[path[i]] = tmp
		}

		if data, ok = tmp.(map[string]interface{}); !ok {
			data = make(map[string]interface{})
		}

		i++
	}

	switch value.(type) {
	case int:
		value = float64(value.(int))
	case int8:
		value = float64(value.(int8))
	case int16:
		value = float64(value.(int16))
	case int32:
		value = float64(value.(int32))
	case int64:
		value = float64(value.(int64))
	case uint:
		value = float64(value.(uint))
	case uint8:
		value = float64(value.(uint8))
	case uint16:
		value = float64(value.(uint16))
	case uint32:
		value = float64(value.(uint32))
	case uint64:
		value = float64(value.(uint64))
	case float32:
		value = float64(value.(float32))
	}

	data[path[i]] = value
}

// Get is a general method for getting the value
// from config as interface{} type.
// The key washes to be separated by a dot symbol.
// For example:
//  wnx.Config.Get("key1")
//  wnx.Config.Get("key2.key3")
func (c *Config) Get(key string) interface{} {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	data := c.data
	path := strings.Split(key, ".")
	var i int

	for i < len(path)-1 {
		tmp, ok := data[path[i]]
		if !ok {
			return nil
		}

		if data, ok = tmp.(map[string]interface{}); !ok {
			return nil
		}

		i++
	}

	return data[path[i]]
}

// Bool returns a value as a boolean type or error.
func (c *Config) Bool(key string) (bool, error) {
	valueInterface := c.Get(key)
	if valueInterface == nil {
		return false, ErrConfigValueNotFound
	}

	valueBool, ok := c.Get(key).(bool)
	if !ok {
		return false, ErrConfigValueMismatched
	}

	return valueBool, nil
}

// MustBool returns a value as a boolean type or runtime panic.
func (c *Config) MustBool(key string) bool {
	valueBool, err := c.Bool(key)
	if err != nil {
		panic(err)
	}

	return valueBool
}

// Float64 returns a value as a float64 type or error.
func (c *Config) Float64(key string) (float64, error) {
	valueInterface := c.Get(key)
	if valueInterface == nil {
		return 0, ErrConfigValueNotFound
	}

	valueFloat64, ok := c.Get(key).(float64)
	if !ok {
		return 0, ErrConfigValueMismatched
	}

	return valueFloat64, nil
}

// MustFloat64 returns a value as a float64 type or runtime panic.
func (c *Config) MustFloat64(key string) float64 {
	valueFloat64, err := c.Float64(key)
	if err != nil {
		panic(err)
	}

	return valueFloat64
}

// String returns a value as a string type or error.
func (c *Config) String(key string) (string, error) {
	valueInterface := c.Get(key)
	if valueInterface == nil {
		return "", ErrConfigValueNotFound
	}

	valueString, ok := valueInterface.(string)
	if !ok {
		return "", ErrConfigValueMismatched
	}

	return valueString, nil
}

// MustString returns a value as a string type or runtime panic.
func (c *Config) MustString(key string) string {
	valueString, err := c.String(key)
	if err != nil {
		panic(err)
	}

	return valueString
}

// Slice returns a value as a []interface{} type or error.
func (c *Config) Slice(key string) ([]interface{}, error) {
	valueInterface := c.Get(key)
	if valueInterface == nil {
		return nil, ErrConfigValueNotFound
	}

	valueSlice, ok := c.Get(key).([]interface{})
	if !ok {
		return nil, ErrConfigValueMismatched
	}

	return valueSlice, nil
}

// MustSlice returns a value as a []interface{} type or runtime panic.
func (c *Config) MustSlice(key string) []interface{} {
	valueSlice, err := c.Slice(key)
	if err != nil {
		panic(err)
	}

	return valueSlice
}

// Map returns a value as a map[string]interface{} type or error.
func (c *Config) Map(key string) (map[string]interface{}, error) {
	valueInterface := c.Get(key)
	if valueInterface == nil {
		return nil, ErrConfigValueNotFound
	}

	valueMap, ok := c.Get(key).(map[string]interface{})
	if !ok {
		return nil, ErrConfigValueMismatched
	}

	return valueMap, nil
}

// MustMap returns a value as a map[string]interface{} type or runtime panic.
func (c *Config) MustMap(key string) map[string]interface{} {
	valueMap, err := c.Map(key)
	if err != nil {
		panic(err)
	}

	return valueMap
}
