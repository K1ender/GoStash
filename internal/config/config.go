package config

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

type Config struct {
	Host string `cfg:"host,default:localhost"`
	Port int    `cfg:"port,default:8080"`
}

func LoadConfig(from string) *Config {
	var cfg Config

	switch from {
	case "config":
		panic("not implemented")
	case "cli":
		getter := NewCLIGetter()
		getter.Run()
		load(&cfg, getter)
	}

	return &cfg
}

type Getter interface {
	Get(string) any
}

func load(cfg *Config, getter Getter) {
	typ := reflect.TypeOf(*cfg)
	val := reflect.ValueOf(cfg).Elem()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("cfg")
		if tag == "" {
			continue
		}
		f := val.FieldByName(field.Name)
		if !f.CanSet() {
			continue
		}

		parts := strings.Split(tag, ",")
		tag = parts[0]
		val := getter.Get(tag)
		isZeroValue := val == nil || reflect.ValueOf(val).IsZero()

		if slices.Contains(parts, "required") && isZeroValue {
			panic(fmt.Sprintf("missing required config: %s", tag))
		} else if isZeroValue {
			for _, part := range parts[1:] {
				if after, ok := strings.CutPrefix(part, "default:"); ok {
					switch f.Kind() {
					case reflect.String:
						val = after
					case reflect.Int:
						intVal, err := strconv.Atoi(after)
						if err != nil {
							panic(fmt.Sprintf("invalid default int value for %s: %s", tag, after))
						}
						val = intVal
					}
					break
				}
			}
		}

		switch f.Kind() {
		case reflect.String:
			f.SetString(val.(string))
		case reflect.Int:
			intVal := val.(int)
			f.SetInt(int64(intVal))
		}
	}
}
