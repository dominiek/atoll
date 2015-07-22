
package main

import (
  "os"
  "fmt"
  "log"
  "github.com/codegangsta/cli"
)

func fatalError(err error) {
  log.Fatalf("Error: %v", err)
}

func main() {
  app := cli.NewApp()
  app.Name = "atoll"
  app.Usage = "Monitoring agent for Atoll"
  app.Version = "0.0.1"
  app.Flags = []cli.Flag {
    cli.StringFlag{
      Name: "config, c",
      Value: "atoll.yml",
      Usage: "Configuration file to use",
    },
    cli.StringFlag{
      Name: "hostname, hn",
      Value: "",
      Usage: "Set primary hostname for this node",
    },
    cli.StringFlag{
      Name: "frequency, f",
      Value: "",
      Usage: "Set reporting frequency (Default: 5s)",
    },
  }
  app.Action = func(c *cli.Context) {
    var err error;

    config := Config{}
    err = config.LoadFile(c.String("config"))
    if err != nil {
      fatalError(err)
    }

    hostname := c.String("hostname");
    if len(hostname) > 0 {
      config.Hostname = hostname;
    }

    frequency := c.String("frequency");
    if len(frequency) > 0 {
      config.Publish.Frequency = frequency;
    }

    var data []byte;
    data, err = config.ToJSON()
    fmt.Printf("%s\n", data)

    url := fmt.Sprintf("http://%s:%d/1/report", config.Publish.Host, config.Publish.Port)
    log.Printf("Publish URL: %s\n", url);
    log.Printf("Publish Frequency: %s\n", config.Publish.Frequency);
    reporter := Reporter{&config, Netstat{config: &config}, true, "netstat", url};
    err = reporter.Start();
    if err != nil {
      log.Fatalf("Error: %v\n", err)
    }
  }

  app.Run(os.Args)
}
