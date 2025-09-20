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

// Run parses command-line flags for "host" and "port", and stores their values
// in the CLIGetter's args map. It uses the standard flag package to define and
// parse the flags, then iterates over the set flags to populate the args map
// with their names and values.
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
