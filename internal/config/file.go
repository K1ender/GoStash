package config

import (
	"os"
	"path"
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

func (f *FileGetter) Load(filePath string) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	filepath := path.Join(cwd, filePath)
	file, err := os.ReadFile(filepath)

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
