package config

import (
	"os"
	"path/filepath"
	"strings"
)

type FileGetter struct {
	data map[string]any
}

func NewFileGetter() *FileGetter {
	return &FileGetter{
		data: make(map[string]any),
	}
}

func (f *FileGetter) Get(key string) any {
	val := f.data[key]
	return val
}

// Load reads a configuration file from the specified filePath, parses its contents,
// and stores key-value pairs into the FileGetter's data map. The method resolves
// the absolute path using the current working directory and the provided filePath.
// It ignores empty lines and lines starting with '#' or '//' (treated as comments).
// Each non-comment, non-empty line is expected to be in the format "key = value".
// Lines not matching this format are skipped. If any error occurs while obtaining
// the working directory or reading the file, the method panics.
func (f *FileGetter) Load(filePath string) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fullPath := filepath.Join(cwd, filePath)
	file, err := os.ReadFile(fullPath)

	if err != nil {
		panic(err)
	}

	for line := range strings.SplitSeq(string(file), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// skip comments
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		kv := strings.SplitN(line, "=", 2)
		if len(kv) != 2 {
			continue
		}

		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])
		f.data[key] = val
	}
}
