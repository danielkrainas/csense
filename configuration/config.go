package configuration

import (
	"fmt"
	"io"
	"io/ioutil"
	"reflect"

	cfg "github.com/danielkrainas/gobag/configuration"
)

type LogConfig struct {
	Level     cfg.LogLevel           `yaml:"level,omitempty"`
	Formatter string                 `yaml:"formatter,omitempty"`
	Fields    map[string]interface{} `yaml:"fields,omitempty"`
}

type CORSConfig struct {
	Debug   bool     `yaml:"debug"`
	Origins []string `yaml:"origins"`
	Methods []string `yaml:"methods"`
	Headers []string `yaml:"headers"`
}

type HTTPConfig struct {
	Enabled bool       `yaml:"enabled"`
	Addr    string     `yaml:"addr"`
	Host    string     `yaml:"host"`
	CORS    CORSConfig `yaml:"cors"`
}

type Config struct {
	Log        LogConfig  `yaml:"logging"`
	Containers cfg.Driver `yaml:"containers"`
	HTTP       HTTPConfig `yaml:"http"`
	Storage    cfg.Driver `yaml:"storage"`
}

type v0_1Config Config

func newConfig() *Config {
	config := &Config{
		Log: LogConfig{
			Level:     "debug",
			Formatter: "text",
			Fields:    make(map[string]interface{}),
		},

		Containers: make(cfg.Driver),

		HTTP: HTTPConfig{
			Enabled: true,
			Addr:    ":9181",
			Host:    "localhost",
		},
	}

	return config
}

func Parse(rd io.Reader) (*Config, error) {
	in, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, err
	}

	p := cfg.NewParser("csense", []cfg.VersionedParseInfo{
		{
			Version: cfg.MajorMinorVersion(0, 1),
			ParseAs: reflect.TypeOf(v0_1Config{}),
			ConversionFunc: func(c interface{}) (interface{}, error) {
				if v0_1, ok := c.(*v0_1Config); ok {
					if v0_1.Containers.Type() == "" {
						return nil, fmt.Errorf("no containers configuration provided")
					}

					return (*Config)(v0_1), nil
				}

				return nil, fmt.Errorf("Expected *v0_1Config, received %#v", c)
			},
		},
	})

	config := new(Config)
	err = p.Parse(in, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
