package config

import "flag"

type CLIGetter struct {
	args map[string]any
}

func NewCLIGetter() *CLIGetter {
	return &CLIGetter{
		args: make(map[string]any),
	}
}

func (c *CLIGetter) Run() {
	flag.String("host", "", "server host")
	flag.Int("port", 0, "server port")

	flag.Parse()

	flag.Visit(func(f *flag.Flag) {
		c.args[f.Name] = f.Value.String()
	})
}

func (c *CLIGetter) Get(key string) any {
	val := c.args[key]
	return val
}
