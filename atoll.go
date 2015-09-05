
package main

import (
  "os"
  "fmt"
  "log"
  "github.com/codegangsta/cli"
  "github.com/VividCortex/godaemon"
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
    cli.BoolFlag{
      Name: "detach, d",
      Usage: "Detach process and run as daemon",
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
  app.Commands = []cli.Command{
    {
      Name:      "setup",
      Usage:     "Generate atoll.yml configuration file",
      Flags:     []cli.Flag {
        cli.StringFlag{
          Name: "publish-host",
          Value: "api.atoll.io",
          Usage: "Publish host to configure",
        },
        cli.IntFlag{
          Name: "publish-port",
          Value: 47011,
          Usage: "Publish port to configure",
        },
        cli.StringFlag{
          Name: "frequency",
          Value: "5s",
          Usage: "Update frequency to configure",
        },
        cli.StringFlag{
          Name: "hostname",
          Value: "",
          Usage: "Node hostname to configure",
        },
      },
      Action: func(c *cli.Context) {
        config := Config{}
        config.Publish.Host = c.String("publish-host")
        config.Publish.Port = c.Int("publish-port")
        config.Publish.Frequency = c.String("frequency")
        config.Hostname = c.String("hostname")
        path := c.GlobalString("config")
        log.Printf("Writing configuration to: %s\n", path);
        config.StoreFile(path);
      },
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

    if (c.Bool("detach")) {
      godaemon.MakeDaemon(&godaemon.DaemonAttr{})
    }

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
