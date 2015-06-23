
package main

import (
  "os/exec"
  "log"
  "strings"
)

type NetstatHost string;
type NetstatPort string;
type NetstatAddress struct {
  Host NetstatHost
  Port NetstatPort
}
type NetstatAddressPair struct {
  Local NetstatAddress
  Remote NetstatAddress
}
type NetstatConnection struct {
  Host NetstatHost
  Count uint64
}

type NetstatConnections map[NetstatHost]NetstatConnection;
type NetstatServices map[NetstatPort]NetstatConnections;

type NetstatInfo struct {
  Outgoing NetstatServices
  Incoming NetstatServices
}

type Netstat struct {
  config *Config;
};

func (this *Netstat) Monitor() (NetstatInfo, error) {
  data, err := this.run();
  var result NetstatInfo;
  if err != nil {
    return result, err;
  }
  log.Printf("%v\n", data)
  return result, err;
}

func (this *Netstat) run() (string, error) {
  var (
    commandStr string
    data []byte
    err error
  )
  commandStr, err = this.determineCommand();
  if err != nil {
    return string(data), err
  }
  command := strings.Split(commandStr, " ");
  cmd := exec.Command(command[0], command[1:]...);
  data, err = cmd.Output();
  return string(data), err;
}

func (this *Netstat) determineCommand() (string, error) {
  data, err := exec.Command("uname").Output();
  if err != nil {
    return "", err
  }
  switch strings.TrimSpace(string(data)) {
    case "Darwin":
      return "netstat -na -p tcp -b", nil
    default:
      return "netstat -antl", nil
  }
}

func (this *Netstat) parse(data string) (NetstatInfo, error) {
  result := NetstatInfo{}
  result.Incoming = make(NetstatServices)
  listenAddresses := this.parseAddressPairs(data, "LISTEN")
  for i := 0; len(listenAddresses) > i; i++ {
    result.Incoming[listenAddresses[i].Local.Port] = make(NetstatConnections)
  }
  establishedAddresses := this.parseAddressPairs(data, "ESTABLISHED")
  for i := 0; len(establishedAddresses) > i; i++ {
    localPort := establishedAddresses[i].Local.Port
    if _, ok := result.Incoming[localPort]; ok == true {
      remoteHost := establishedAddresses[i].Remote.Host
      if _, ok := result.Incoming[localPort][remoteHost]; ok == true {
        result.Incoming[localPort][remoteHost] = NetstatConnection{
          remoteHost, 
          result.Incoming[localPort][remoteHost].Count + 1,
        };
      } else {
        result.Incoming[localPort][remoteHost] = NetstatConnection{
          remoteHost, 
          1,
        };
      }
    }
    
  }
  return result, nil
}

func (this *Netstat) parseAddressPairs(data string, state string) []NetstatAddressPair {
  lines := strings.Split(data, "\n");
  pairs := make([]NetstatAddressPair, 0);
  for i := 1; len(lines) > i; i++ {
    line := strings.Fields(lines[i])
    if (len(line) < 4) {
      continue
    }
    if (line[0] == "tcp6") {
      continue
    }
    if line[5] == state {
      localAddress := this.parseAddress(line[3])
      remoteAddress := this.parseAddress(line[4])
      pair := NetstatAddressPair{
        localAddress,
        remoteAddress,
      }
      pairs = append(pairs, pair)
      log.Printf("%v %v %v", state, localAddress.Host, localAddress.Port);
    }
    
  }
  return pairs;
}

func (this *Netstat) parseAddress(data string) NetstatAddress {
  hostPort := strings.Split(data, ":")
  if (len(hostPort) > 1) {
    return NetstatAddress{
      NetstatHost(hostPort[0]), 
      NetstatPort(hostPort[len(hostPort)-1]),
    }
  } else {
    hostPort = strings.Split(data, ".")
    host := strings.Join(hostPort[0:len(hostPort)-1], ".")
    return NetstatAddress{
      NetstatHost(host), 
      NetstatPort(hostPort[len(hostPort)-1]),
    }
  }
}

