package collector

import (
    "io/ioutil"
    "path/filepath"
    "gopkg.in/yaml.v2"
)

type kstatConfig struct {
	KstatModules []struct {
		ID         string `yaml:"id"`
		KstatNames []struct {
			ID         string `yaml:"id"`
			KstatStats []struct {
				ID		string `yaml:"id"`
				Suffix  	string `yaml:"suffix"`
				Type 		string `yaml:"type"`
				Help 		string `yaml:"help"`
				ScaleFactor 	float64	`yaml:"scale_factor"`
			} `yaml:"kstat_stats"`
		} `yaml:"kstat_names"`
	} `yaml:"kstat_modules"`
}

func (config *kstatConfig) init() error {
	var (
		err error
	)
	filename, _ := filepath.Abs("./collector/kstat_config.yml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil { return err }

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil { return err }

	return nil
}
