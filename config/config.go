package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type TargetInfo struct {
	Name string `json:"name"`
	Mac  string `json:"mac"`
}

type TargetMap map[string]string

type Config struct {
	Url      string       `json:"url"`
	User     string       `json:"user"`
	Password string       `json:"password"`
	Targets  []TargetInfo `json:"targets"`
}

func Default() Config {
	var cfg = Config{
		"user.ddns.net:1222",
		"admin",
		"admin",
		[]TargetInfo{
			{"PC", "11:22:33:44:55:66"},
			{"Printer", "AA:BB:CC:DD:EE:FF"},
		},
	}
	return cfg
}

func (c *Config) FindTarget(name string) (string, error) {
	for _, t := range c.Targets {
		if t.Name == name {
			return t.Mac, nil
		}
	}
	return "", errors.New("target not found")
}

func Load(fileName string) (Config, error) {
	var cfg Config
	file, err := os.Open(fileName)
	if err != nil {
		return cfg, err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (c *Config) Save(fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	return err
}
