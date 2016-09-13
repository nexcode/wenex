package wenex

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
)

// DefaultConfig is default configuration options
var DefaultConfig = map[string]interface{}{
	"server": map[string]interface{}{
		"http": map[string]interface{}{
			"listen": ":3000",
		},
		"timeout": map[string]interface{}{
			"read":  "30s",
			"write": "30s",
		},
		// "https": map[string]interface{}{
		// 	"listen": ":https",
		// 	"crt":    "file.crt",
		// 	"key":    "file.key",
		// },
	},
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
