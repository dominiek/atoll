
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

  reporters := make([]Reporter, 0);
  reporters = append(reporters, Reporter{this.config, Netstat{config: this.config}, true, "netstat", url});

  for _, command := range this.config.Plugins {
    plugin := Plugin{this.config, command};
    reporters = append(reporters, Reporter{this.config, plugin, true, "plugin", url});
  }

  this.running = true;
  numSeconds, err := IntervalToSeconds(this.config.Publish.Frequency);
  if numSeconds == 0 {
    return errors.New("Need an interval of at least 1 second")
  }
  if err != nil {
    return errors.New("Invalid interval configured")
  }
  for this.running == true {
    for _, reporter := range reporters {
      err := reporter.Report();
      if err != nil {
        log.Printf("Warning, could not report: %v\n", err);
      }
    }
    time.Sleep(time.Duration(numSeconds) * time.Second);
  }
  return nil
}

func (this *Daemon) Stop() () {
  this.running = false;
}
