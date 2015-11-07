
package main

import (
  "time"
  "fmt"
  "errors"
  "log"
)

type Daemon struct {
  config *Config;
  running bool;
}

func (this *Daemon) Start() (error) {
  url := fmt.Sprintf("http://%s:%d/1/report/%s", this.config.Publish.Host, this.config.Publish.Port, this.config.Publish.ApiKey)
  log.Printf("Publish URL: %s\n", url);
  log.Printf("Publish Frequency: %s\n", this.config.Publish.Frequency);

  netstat := Netstat{config: this.config}
  reporters := make([]Reporter, 0);
  reporters = append(reporters, Reporter{this.config, netstat, true, "netstat", url});

  for _, command := range this.config.Plugins {
    plugin := Plugin{this.config, command};
    reporters = append(reporters, Reporter{this.config, plugin, true, "plugin", url});
  }

  this.running = true;
  publishFrequency, err := IntervalToSeconds(this.config.Publish.Frequency);
  if publishFrequency == 0 {
    return errors.New("Need a publish frequency of at least 1 second")
  }
  if err != nil {
    return errors.New("Invalid reporting interval configured")
  }

  netstatCheckFrequency := 1.0;
  if (len(this.config.Netstat.CheckFrequency) > 0) {
    netstatCheckFrequency, err = IntervalToSeconds(this.config.Netstat.CheckFrequency);
    if err != nil {
      return errors.New("Invalid netstat check interval configured")
    }
  }

  lastReportTs := time.Now();
  for this.running == true {
    if (time.Now().Sub(lastReportTs) >= (time.Duration(netstatCheckFrequency) * time.Second)) {
      _, err = netstat.Monitor();
      if err != nil {
        log.Printf("Warning, could not monitor netstat: %v\n", err);
      }
    }
    if (time.Now().Sub(lastReportTs) >= (time.Duration(publishFrequency) * time.Second)) {
      lastReportTs = time.Now();
      for _, reporter := range reporters {
        err := reporter.Report();
        if err != nil {
          log.Printf("Warning, could not report: %v\n", err);
        }
      }
      netstat.ClearResults();
    }
    time.Sleep(time.Duration(100) * time.Millisecond);
  }
  return nil
}

func (this *Daemon) Stop() () {
  this.running = false;
}
