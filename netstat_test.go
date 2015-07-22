package main

import (
  "testing"
  "github.com/stretchr/testify/assert"
  "encoding/json"
)

func TestNetstatMonitor(t *testing.T) {
  config := Config{};
  var err error;
  err = config.LoadFile("./atoll.yml")
  if err != nil {
    t.Fatalf("Did not expect error %v", err)
  }
  netstat := Netstat{config: &config}
  var info Info;
  info, err = netstat.Monitor();
  if err != nil {
    t.Fatalf("Did not expect error %v", err)
  }
  result := info.(NetstatInfo)
  keys := []NetstatPort{}
  for k := range result.Outgoing {
    keys = append(keys, k)
  }
  assert.Equal(t, len(keys) > 0, true)
}

func TestNetstatExec(t *testing.T) {
  netstat := Netstat{}
  data, err := netstat.run();
  if err != nil {
    t.Fatalf("Did not expect error %v", err)
  }
  assert.Contains(t, string(data), "LISTEN")
}

const darwinOutput string = `Proto Recv-Q Send-Q  Local Address          Foreign Address        (state)        rxbytes    txbytes
tcp4       0      0  192.168.0.8.56830      54.241.191.235.443     ESTABLISHED       5577       1776
tcp4       0      0  192.168.0.8.56826      216.58.216.14.443      ESTABLISHED        494       1552
tcp4       0      0  192.168.0.8.56797      192.168.0.2.445        ESTABLISHED        419        711
tcp4       0      0  127.0.0.1.5556         *.*                    LISTEN               0          0
tcp4       0      0  192.168.0.8.5556       *.*                    LISTEN               0          0
tcp4      37      0  192.168.0.8.56755      54.230.86.171.443      CLOSE_WAIT        6000       1316
tcp4       0      0  192.168.0.8.56746      54.231.0.212.80        ESTABLISHED       1542       1208
tcp4       0      0  192.168.0.8.56726      216.58.216.14.443      ESTABLISHED      15743      17743
tcp4       0      0  192.168.0.8.56625      17.172.232.211.5223    ESTABLISHED       4274       3344
tcp4       0      0  192.168.0.8.56551      74.125.20.189.443      ESTABLISHED       5644       7882
tcp4       0      0  192.168.0.8.56545      74.125.25.188.5228     ESTABLISHED       4871       1303
tcp4       0      0  192.168.0.8.56044      199.59.149.201.443     ESTABLISHED       4812       3644
tcp4      37      0  192.168.0.8.55882      199.47.217.65.443      CLOSE_WAIT     1295480      42617
tcp4      37      0  192.168.0.8.55881      199.47.217.65.443      CLOSE_WAIT     1733722      54489
tcp4      37      0  192.168.0.8.55755      54.192.86.126.443      CLOSE_WAIT        4508       2660
tcp4      37      0  192.168.0.8.55462      199.47.217.2.443       CLOSE_WAIT     2772639      28700
tcp4      37      0  192.168.0.8.55461      199.47.217.2.443       CLOSE_WAIT     4135279      35044
tcp4      37      0  192.168.0.8.55460      54.230.86.65.443       CLOSE_WAIT       17968       6388
tcp4      37      0  192.168.0.8.55233      54.230.86.93.443       CLOSE_WAIT       17840       6440
tcp4       0      0  192.168.0.8.54655      185.11.124.4.443       ESTABLISHED      21452      33703
tcp4      37      0  192.168.0.8.54588      108.160.172.236.443    CLOSE_WAIT        6407       2493
tcp4       0      0  192.168.0.8.54553      104.16.34.27.443       ESTABLISHED    2426679     518832
tcp4       0      0  192.168.0.8.54549      54.236.177.53.443      ESTABLISHED      77357     283733
tcp4       0      0  192.168.0.8.54546      54.236.177.53.443      ESTABLISHED      70872     284710
tcp4       0      0  192.168.0.8.54535      192.168.0.14.8009      ESTABLISHED     230604     257165
tcp4      37      0  192.168.0.8.54527      107.20.249.250.443     CLOSE_WAIT       29475    6962232
tcp4       0      0  192.168.0.8.54526      54.237.50.81.443       ESTABLISHED       7626      10527
tcp4       0      0  192.168.0.8.54515      192.168.0.14.8008      ESTABLISHED      50400      25058
tcp4      37      0  192.168.0.8.54506      108.160.172.225.443    CLOSE_WAIT        6071       8260
tcp4       0      0  192.168.0.8.54502      108.160.170.46.443     ESTABLISHED      35369      90520
tcp4       0      0  192.168.0.8.54450      54.192.86.126.443      ESTABLISHED       5859       1160
tcp4       0      0  192.168.0.8.53454      192.241.163.235.443    ESTABLISHED       4866       4158
tcp4       0      0  192.168.0.8.53453      192.30.252.91.443      ESTABLISHED       8271      14535
tcp4      37      0  192.168.150.132.53041  199.47.217.65.443      CLOSE_WAIT        6199   15255780
tcp4      37      0  192.168.150.132.53040  199.47.217.65.443      CLOSE_WAIT        6736   15259088
tcp4      37      0  192.168.150.132.53039  199.47.217.65.443      CLOSE_WAIT        6199   10367042
tcp4      37      0  192.168.150.132.53037  54.230.86.171.443      CLOSE_WAIT       35518      15330
tcp4      37      0  192.168.150.132.53036  199.47.217.65.443      CLOSE_WAIT        7273   16872882
tcp4      37      0  192.168.150.132.53035  54.230.86.171.443      CLOSE_WAIT       39364      20304
tcp4       0      0  127.0.0.1.26164        127.0.0.1.53034        ESTABLISHED      15084          0
tcp4       0      0  127.0.0.1.53034        127.0.0.1.26164        ESTABLISHED       5883          0
tcp4       0      0  *.17500                *.*                    LISTEN               0          0
tcp4       0      0  127.0.0.1.17603        *.*                    LISTEN               0          0
tcp4       0      0  127.0.0.1.17600        *.*                    LISTEN               0          0
tcp4       0      0  127.0.0.1.26164        *.*                    LISTEN               0          0
tcp4      37      0  192.168.0.8.49297      108.160.172.236.443    CLOSE_WAIT        8749       1697
tcp4       0      0  127.0.0.1.49160        127.0.0.1.49173        ESTABLISHED 1748734920          0
tcp4       0      0  127.0.0.1.49173        127.0.0.1.49160        ESTABLISHED  863929844          0
tcp4       0      0  127.0.0.1.49160        *.*                    LISTEN               0          0
tcp4       0      0  *.631                  *.*                    LISTEN               0          0
tcp6       0      0  *.631                  *.*                    LISTEN               0          0
tcp6       0      0  fe80::1%lo0.8021       *.*                    LISTEN               0          0
tcp4       0      0  127.0.0.1.8021         *.*                    LISTEN               0          0
tcp6       0      0  ::1.8021               *.*                    LISTEN               0          0
tcp4       0      0  *.22                   *.*                    LISTEN               0          0
tcp6       0      0  *.22                   *.*                    LISTEN               0          0
tcp4       0      0  127.0.0.1.631          *.*                    LISTEN               0          0
tcp6       0      0  ::1.631                *.*                    LISTEN               0          0
`;

func TestNetstatParseDarwin(t *testing.T) {
  netstat := Netstat{}
  result, err := netstat.parse(darwinOutput);
  if err != nil {
    t.Fatalf("Did not expect error %v", err)
  }
  assert.EqualValues(t, result.Outgoing["443"]["216.58.216.14"].Count, 2)
  data, err := json.Marshal(result)
  if err != nil {
    t.Fatalf("Did not expect error %v", err)
  }
  t.Logf("JSON %s", data)
  assert.Contains(t, string(data), `127.0.0.1":{"Host":"127.0.0.1","Count":1}`)
}

const linuxOutput string = `Active Internet connections (servers and established)
Proto Recv-Q Send-Q Local Address           Foreign Address         State      
tcp        0      0 0.0.0.0:22              0.0.0.0:*               LISTEN     
tcp        0      0 0.0.0.0:15672           0.0.0.0:*               LISTEN     
tcp        0      0 0.0.0.0:25              0.0.0.0:*               LISTEN     
tcp        0      0 0.0.0.0:25672           0.0.0.0:*               LISTEN     
tcp        0      0 10.45.10.220:57985      10.45.10.220:5672       ESTABLISHED
tcp        0      0 10.45.10.220:49172      10.45.10.222:27018      ESTABLISHED
tcp        0      0 10.45.10.220:46965      10.45.10.222:27018      ESTABLISHED
tcp        0      0 10.45.10.220:40923      10.45.10.193:27018      ESTABLISHED
tcp        0      0 10.45.10.220:46170      10.45.10.222:27018      ESTABLISHED
tcp        0      0 10.45.10.220:57983      10.45.10.198:443       ESTABLISHED
`;

func TestNetstatParseLinux(t *testing.T) {
  netstat := Netstat{}
  result, err := netstat.parse(linuxOutput);
  if err != nil {
    t.Fatalf("Did not expect error %v", err)
  }
  assert.EqualValues(t, result.Outgoing["27018"]["10.45.10.222"].Count, 3)
  data, err := json.Marshal(result)
  if err != nil {
    t.Fatalf("Did not expect error %v", err)
  }
  t.Logf("JSON %s", data)
  assert.Contains(t, string(data), `"22":{}`)
}
