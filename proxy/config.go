package proxy

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"regexp"
)

type Config struct {
	Entries []string `yaml:"block"`
	log Logger
}

func NewConfig() *Config {
	return &Config{Entries: nil, log:Logger{}}
}

func (c *Config) LoadBlockedList(filename string) []regexp.Regexp {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		c.log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	regexps := []regexp.Regexp{}
	for _, condition := range c.Entries {
		//log.Printf("%s", condition)
		r := regexp.MustCompile(condition)
		regexps = append(regexps, *r)
	}
	return regexps
}
