
package main

import (
  "time"
  "bytes"
  "net/http"
  "errors"
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
  hostnames []string;
  uname string;
  unameA string;
}

func (this *Reporter) Report() (error) {
  info, err := this.module.Monitor();
  if err != nil {
    return err;
  }
  data, err := info.Encode();
  if err != nil {
    return err
  }
  req, err := http.NewRequest("POST", this.url, bytes.NewBuffer(data))
  req.Header.Set("Content-Type", "application/json")
  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    return err
  }
  defer resp.Body.Close()
  if resp.StatusCode != 200 {
    return errors.New("Invalid response from Atoll server: " + resp.Status)
  }
  return nil
}

func (this *Reporter) GetHostInfo() (HostInfo) {
  hostInfo := HostInfo{};
  hostInfo.hostnames = make([]string, 0);
  if len(this.config.HOSTNAME) > 0 {
    hostInfo.hostnames = append(hostInfo.hostnames, this.config.HOSTNAME)
  }
  // TODO get all hostnames and IP addresses for host
  return hostInfo;
}

func (this *Reporter) Start() (error) {
  for this.running == true {
     println("Hello")
     time.Sleep(3000 * time.Millisecond);
  }
  return nil
}

func (this *Reporter) Stop() () {
  this.running = false;
}