
package main

import (
  "io/ioutil"
  "gopkg.in/yaml.v2"
  "encoding/json"
)

type Config struct {
  Hostname string `json:"hostname"`;
  Server struct {
    Bind string    `json:"bind"`;
    Port int       `json:"port"`;
  }               `json:"server"`;
  Publish struct {
    Host string    `json:"host"`;
    Port int       `json:"port"`;
    Frequency string;
  }               `json:"publish"`;
  Netstat struct {
    IncludeLocal bool `json:"includeLocal"`;
  }
}

func (config *Config) LoadFile(path string) error {
  data, err := ioutil.ReadFile(path)
  if err != nil {
    return err
  }
  err = yaml.Unmarshal(data, config)
  if err != nil {
    return err
  }
  return nil
}

func (config *Config) ToJSON() ([]byte, error) {
  return json.Marshal(config)
}
