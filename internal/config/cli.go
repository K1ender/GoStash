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
	hostPtr := flag.String("host", "", "server host")
	portPtr := flag.Int("port", 0, "server port")

	flag.Parse()

	c.args["host"] = *hostPtr
	c.args["port"] = *portPtr
}

func (c *CLIGetter) Get(key string) any {
	val, ok := c.args[key]
	if !ok {
		return nil
	}

	return val
}
