package wenex

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
)

// GetDefaultConfig returns default configuration options
func GetDefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"server.http.listen":   ":3000",
		"server.timeout.read":  "30s",
		"server.timeout.write": "30s",
		"log.filePrefix":       "",
	}
}

func newConfig(name string) (*Config, error) {
	file, err := os.OpenFile(name+".conf", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
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

	return &config, nil
}

type Config struct {
	mutex   sync.Mutex
	file    *os.File
	decoder *json.Decoder
	encoder *json.Encoder
	data    map[string]interface{}
}

func (c *Config) Load() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, err := c.file.Seek(0, 0); err != nil {
		return err
	}

	if err := c.decoder.Decode(&c.data); err != nil {
		return err
	}

	return nil
}

func (c *Config) Save() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, err := c.file.Seek(0, 0); err != nil {
		return err
	}

	if err := c.file.Truncate(0); err != nil {
		return err
	}

	if err := c.encoder.Encode(c.data); err != nil {
		return err
	}

	return nil
}

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

func (c *Config) Bool(key string) (bool, error) {
	value, ok := c.Get(key).(bool)
	if !ok {
		return false, ErrGetFromConfig
	}

	return value, nil
}

func (c *Config) MustBool(key string) bool {
	value, ok := c.Get(key).(bool)
	if !ok {
		panic(ErrGetFromConfig)
	}

	return value
}

func (c *Config) Float64(key string) (float64, error) {
	value, ok := c.Get(key).(float64)
	if !ok {
		return 0, ErrGetFromConfig
	}

	return value, nil
}

func (c *Config) MustFloat64(key string) float64 {
	value, ok := c.Get(key).(float64)
	if !ok {
		panic(ErrGetFromConfig)
	}

	return value
}

func (c *Config) String(key string) (string, error) {
	value, ok := c.Get(key).(string)
	if !ok {
		return "", ErrGetFromConfig
	}

	return value, nil
}

func (c *Config) MustString(key string) string {
	value, ok := c.Get(key).(string)
	if !ok {
		panic(ErrGetFromConfig)
	}

	return value
}

func (c *Config) Slice(key string) ([]interface{}, error) {
	value, ok := c.Get(key).([]interface{})
	if !ok {
		return nil, ErrGetFromConfig
	}

	return value, nil
}

func (c *Config) MustSlice(key string) []interface{} {
	value, ok := c.Get(key).([]interface{})
	if !ok {
		panic(ErrGetFromConfig)
	}

	return value
}

func (c *Config) Map(key string) (map[string]interface{}, error) {
	value, ok := c.Get(key).(map[string]interface{})
	if !ok {
		return nil, ErrGetFromConfig
	}

	return value, nil
}

func (c *Config) MustMap(key string) map[string]interface{} {
	value, ok := c.Get(key).(map[string]interface{})
	if !ok {
		panic(ErrGetFromConfig)
	}

	return value
}
