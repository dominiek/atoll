
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
  }
  app.Action = func(c *cli.Context) {
    println("Config: ", c.String("config"))
    var err error;

    config := Config{}
    err = config.LoadFile(c.String("config"))
    if err != nil {
      fatalError(err)
    }
    fmt.Printf("%v\n", config)

    var data []byte;

    data, err = config.ToJSON()
    fmt.Printf("%s\n", data)

  }

  app.Run(os.Args)
}
