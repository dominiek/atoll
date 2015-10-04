
package main

import (
  "os/exec"
  "strings"
  "errors"
  "github.com/jeffail/gabs"
)

type PluginInfo struct {
  raw []byte;
  id string;
  name string;
  reportKeys []string;
  reportData []byte;
}

type Plugin struct {
  config *Config;
  path string;
};

func (this PluginInfo) Encode() ([]byte, error) {
  return this.reportData, nil;
}

func (this PluginInfo) GetType() (string) {
  return this.id;
}

func (this Plugin) Monitor() (Info, error) {
  data, err := this.run(this.path);
  var result PluginInfo;
  if err != nil {
    return result, err;
  }
  result.raw = []byte(data);

  jsonParsed, err := gabs.ParseJSON(result.raw)
  if err != nil {
    return result, err;
  }

  result.id = jsonParsed.Path("id").Data().(string);
  if len(result.id) == 0 {
    return result, errors.New("Invalid plugin response, need id")
  }

  result.name = jsonParsed.Path("name").Data().(string);
  if len(result.name) == 0 {
    return result, errors.New("Invalid plugin response, need name")
  }

  children, err := jsonParsed.S("report").ChildrenMap();
  if err != nil {
    return result, errors.New("Invalid plugin response, need report object")
  }
  result.reportKeys = make([]string, 0);
  for key, _ := range children {
    result.reportKeys = append(result.reportKeys, key);
  }
  result.reportData = []byte(jsonParsed.S("report").String());

  info := Info(result);
  return info, err
}

func (this *Plugin) run(commandStr string) (string, error) {
  var (
    data []byte
    err error
  )
  command := strings.Split(commandStr, " ");
  cmd := exec.Command(command[0], command[1:]...);
  data, err = cmd.Output();
  return string(data), err;
}
