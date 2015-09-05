
package main

import (
  "time"
  "bytes"
  "net"
  "net/http"
  "os/exec"
  "errors"
  "strings"
  "fmt"
  "log"
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
  Hostnames    []string  `json:"hostnames"`;
  IpAddresses  []string  `json:"ipAddresses"`;
  Uname        string    `json:"uname"`;
  UnameA       string    `json:"unameA"`;
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
  log.Printf("Reported %s status to %s\n", this.moduleType, this.url);
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

  data, err := exec.Command("hostname").Output();
  if err == nil {
    hostInfo.Hostnames = append(hostInfo.Hostnames, strings.TrimSpace(string(data)));
  }

  data, err = exec.Command("uname").Output();
  if err == nil {
    hostInfo.Uname = strings.TrimSpace(string(data));
  }

  data, err = exec.Command("uname", "-a").Output();
  if err == nil {
    hostInfo.UnameA = strings.TrimSpace(string(data));
  }

  hostInfo.IpAddresses = this.GetIpAddresses();

  for _,ipAddress := range hostInfo.IpAddresses {
    hostnames, err := net.LookupAddr(ipAddress)
    if err == nil {
      for _,hostname := range hostnames {
        hostInfo.Hostnames = append(hostInfo.Hostnames, hostname);
      }
    }
  }

  return hostInfo;
}

func (this *Reporter) GetIpAddresses() ([]string) {
  ipAddresses := make([]string, 0);
  interfaces, err := net.Interfaces()
  if err != nil {
    log.Printf("Warning could not get net.Interfaces information: %s\n", err);
    return ipAddresses
  }
  for _, i := range interfaces {
    addrs, err := i.Addrs()
    if err == nil {
      for _, addr := range addrs {
        var ip net.IP
        switch v := addr.(type) {
          case *net.IPNet:
            ip = v.IP
          case *net.IPAddr:
            ip = v.IP
        }
        ipnet := addr.(*net.IPNet);
        if (!ipnet.IP.IsLoopback()) {
          ipAddresses = append(ipAddresses, ip.String());
        }
      }
    }
  }
  return ipAddresses;
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
      log.Printf("Warning, could not report: %v\n", err);
    }
    time.Sleep(time.Duration(numSeconds) * time.Second);
  }
  return nil
}

func (this *Reporter) Stop() () {
  this.running = false;
}
