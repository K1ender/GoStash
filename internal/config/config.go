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
	Port int    `cfg:"port,default:19201"`

	ConfigPath string
}

type Arg func(cfg *Config)

func WithConfigPath(path string) Arg {
	return func(cfg *Config) {
		cfg.ConfigPath = path
	}
}

func LoadConfig(from string, args ...Arg) *Config {
	var cfg Config
	for _, arg := range args {
		arg(&cfg)
	}

	switch from {
	case "config":
		getter := NewFileGetter()
		getter.Load(cfg.ConfigPath)
		load(&cfg, getter)
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

// load populates the fields of the provided Config struct pointer (cfg) using values
// obtained from the given Getter interface. It uses reflection to iterate over the struct
// fields, reading the "cfg" struct tag to determine the configuration key, and optional
// modifiers such as "required" and "default:<value>".
//
// For each field with a "cfg" tag:
//   - If the tag contains "required" and the value is missing, it panics.
//   - If the value is missing but a "default" is specified, it uses the default value.
//   - Supports string and int field types, converting values as needed.
//   - Panics if a required value is missing, or if a value cannot be converted to the
//     appropriate type.
//
// Example tag: `cfg:"my_key,required,default:42"`
//
// Panics on missing required fields, invalid default values, or unsupported field types.
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
			panic(fmt.Sprintf("missing required configuration field: %s", tag))
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
					default:
						panic("unsupported field type for default value")
					}
					break
				}
			}
		}

		switch f.Kind() {
		case reflect.String:
			f.SetString(fmt.Sprint(val))
		case reflect.Int:
			var intVal int64
			var err error
			switch v := val.(type) {
			case int:
				intVal = int64(v)
			case string:
				intVal, err = strconv.ParseInt(v, 10, 64)
				if err != nil {
					panic(fmt.Sprintf("could not parse int for %s: %v", tag, val))
				}
			default:
				panic(fmt.Sprintf("unsupported type %T for int field %s", val, tag))
			}
			f.SetInt(intVal)
		default:
			panic("unsupported field type")
		}
	}
}
