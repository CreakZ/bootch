package cfg

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	MemcachedHost   string `yaml:"memcached_host"`
	MemcachedPort   int32  `yaml:"memcached_port"`
	CacheExpiration int32  `yaml:"cache_expiration"`
}

func MustInitConfig() Config {
	path := "./cfg.yaml"
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var cfg Config
	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		panic(err)
	}
	return cfg
}
