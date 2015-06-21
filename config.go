
package main

import (
  "io/ioutil"
  "gopkg.in/yaml.v2"
  "encoding/json"
)

type Config struct {
  SERVER struct {
    BIND string;
    PORT int;
  }
  PUBLISH struct {
    HOST string;
    FREQUENCY string;
  }
  NETSTAT struct {
    INCLUDE_LOCAL bool;
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
