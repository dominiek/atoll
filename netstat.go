
package main

import (
  "os/exec"
  "strings"
  "encoding/json"
)

type NetstatHost string;
type NetstatPort string;
type NetstatAddress struct {
  Host NetstatHost `json:"host"`;
  Port NetstatPort `json:"port"`;
}
type NetstatAddressPair struct {
  Local NetstatAddress   `json:"local"`;
  Remote NetstatAddress  `json:"remote"`;
}
type NetstatConnection struct {
  Host NetstatHost `json:"host"`;
  Count uint64     `json:"count"`;
}

type NetstatConnections map[NetstatHost]NetstatConnection;
type NetstatServices map[NetstatPort]NetstatConnections;

type NetstatInfo struct {
  Outgoing NetstatServices `json:"outgoing"`;
  Incoming NetstatServices `json:"incoming"`;
}

type Netstat struct {
  config *Config;
};

func (this NetstatInfo) Encode() ([]byte, error) {
  return json.Marshal(this)
}

func (this NetstatInfo) GetType() (string) {
  return "netstat";
}

func (this Netstat) Monitor() (Info, error) {
  data, err := this.run();
  var result NetstatInfo;
  if err != nil {
    return result, err;
  }
  result, err = this.parse(data)
  info := Info(result);
  return info, err
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
  result.Outgoing = make(NetstatServices)
  listenAddresses := this.parseAddressPairs(data, "LISTEN")
  for i := 0; len(listenAddresses) > i; i++ {
    result.Incoming[listenAddresses[i].Local.Port] = make(NetstatConnections)
  }
  establishedAddresses := this.parseAddressPairs(data, "ESTABLISHED")
  for i := 0; len(establishedAddresses) > i; i++ {
    localPort := establishedAddresses[i].Local.Port
    if service, ok := result.Incoming[localPort]; ok == true {
      remoteHost := establishedAddresses[i].Remote.Host
      if _, ok := service[remoteHost]; ok == true {
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
      continue
    }
    remotePort := establishedAddresses[i].Remote.Port
    if _, ok := result.Outgoing[remotePort]; ok == false {
      result.Outgoing[remotePort] = make(NetstatConnections);
    }
    remoteHost := establishedAddresses[i].Remote.Host
    if _, ok := result.Outgoing[remotePort][remoteHost]; ok == true {
      result.Outgoing[remotePort][remoteHost] = NetstatConnection{
        remoteHost,
        result.Outgoing[remotePort][remoteHost].Count + 1,
      };
    } else {
      result.Outgoing[remotePort][remoteHost] = NetstatConnection{
        remoteHost,
        1,
      };
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
