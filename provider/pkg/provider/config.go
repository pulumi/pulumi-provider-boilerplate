package provider

import "os"

type xyzConfig struct {
	Config map[string]string
}

func (c *xyzConfig) getConfig(configName, envVarName string) string {
	if configVal, ok := c.Config[configName]; ok {
		return configVal
	}

	return os.Getenv(envVarName)
}
