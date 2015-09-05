
package main

import (
  "time"
  "bytes"
  "net/http"
  "errors"
  "os"
  "fmt"
  "encoding/json"
)

type Info interface {
  Encode() ([]byte, error);
}

type Module interface {
  Monitor() (Info, error);
}

type Reporter struct {
  config *Config;
  module Module;
  running bool;
  moduleType string;
  url string;
}

type HostInfo struct {
  Hostnames []string  `json:"hostnames"`;
  Uname     string    `json:"uname"`;
  UnameA    string    `json:"unameA"`;
}

func (this *Reporter) Report() (error) {
  info, err := this.module.Monitor();
  if err != nil {
    return err;
  }
  reportData, err := info.Encode();
  if err != nil {
    return err
  }
  hostData, err := json.Marshal(this.GetHostInfo());
  envelope := fmt.Sprintf(`{"host":%s,"report":%s}`, hostData, reportData);
  req, err := http.NewRequest("POST", this.url, bytes.NewBuffer([]byte(envelope)))
  req.Header.Set("Content-Type", "application/json")
  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    return err
  }
  fmt.Fprintf(os.Stdout, "Reported %s status to %s\n", this.moduleType, this.url);
  defer resp.Body.Close()
  if resp.StatusCode != 200 {
    return errors.New("Invalid response from Atoll server: " + resp.Status)
  }
  return nil
}

func (this *Reporter) GetHostInfo() (HostInfo) {
  hostInfo := HostInfo{};
  hostInfo.Hostnames = make([]string, 0);
  if len(this.config.Hostname) > 0 {
    hostInfo.Hostnames = append(hostInfo.Hostnames, this.config.Hostname);
  }
  // TODO get all hostnames and IP addresses for host
  return hostInfo;
}

func (this *Reporter) Start() (error) {
  numSeconds, err := IntervalToSeconds(this.config.Publish.Frequency);
  if numSeconds == 0 {
    return errors.New("Need an interval of at least 1 second")
  }
  if err != nil {
    return errors.New("Invalid interval configured")
  }
  for this.running == true {
    err := this.Report();
    if err != nil {
      fmt.Fprintf(os.Stderr, "Warning, could not report: %v\n", err);
    }
    time.Sleep(time.Duration(numSeconds) * time.Second);
  }
  return nil
}

func (this *Reporter) Stop() () {
  this.running = false;
}
