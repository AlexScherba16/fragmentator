package cli

import (
	"flag"
	"fmt"
	"fragmentator/internal/constants"
)

type cliParams struct {
	configPath string
}

func (c *cliParams) validateParams() error {
	type paramCheck struct {
		value string
		name  string
	}

	params := []paramCheck{
		{c.ConfigPath(), constants.CliConfigPathParam},
	}

	// Add params validation logic here.
	// Nonempty params are ok for now =)
	for _, item := range params {
		if item.value == "" {
			flag.Usage()
			return fmt.Errorf("%q is required", item.name)
		}
	}
	return nil
}

func (c *cliParams) ConfigPath() string {
	return c.configPath
}

func NewFlags() (cliParams, error) {
	cmd := cliParams{}

	flag.StringVar(&cmd.configPath, constants.CliConfigPathParam, "", "Path to the application config file")

	flag.Parse()

	// Flags validation logic
	if err := cmd.validateParams(); err != nil {
		return cliParams{}, err
	}

	return cmd, nil
}
