package main

import (
  "testing"
  "github.com/stretchr/testify/assert"
  "net/http/httptest"
  "net/http"
  "fmt"
  "io/ioutil"
  "github.com/jeffail/gabs"
)

func TestReporterWithNetstat(t *testing.T) {
  config := Config{};
  config.Hostname = "0.localhost"

  var requestBody []byte
  handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var err error
    requestBody, err = ioutil.ReadAll(r.Body)
    assert.Equal(t, err, nil)
    fmt.Fprintln(w, `{}`)
  })
  ts := httptest.NewServer(handler)
  defer ts.Close();

  url := ts.URL
  reporter := Reporter{&config, Netstat{config: &config}, true, "netstat", url};
  err := reporter.Report();

  t.Logf("Mock URL for reporter: %v", url)
  t.Logf("Data sent: %s", requestBody)
  jsonParsed, err := gabs.ParseJSON(requestBody)
  assert.Equal(t, err, nil)

  hostnames, _ := jsonParsed.S("host").S("hostnames").Children()
  assert.Equal(t, len(hostnames) > 0, true)

  children, _ := jsonParsed.S("report").S("outgoing").ChildrenMap()
  var keys = []string{}
  for key := range children {
    keys = append(keys, key)
  }
  assert.Equal(t, len(keys) > 0, true)
}

func TestReporterGetHostInfo(t *testing.T) {
  config := Config{};
  config.Hostname = "0.localhost"
  reporter := Reporter{&config, Netstat{config: &config}, true, "netstat", "http://localhost:47011"};
  hostInfo := reporter.GetHostInfo();
  t.Logf("hostInfo: %v", hostInfo)
  assert.Equal(t, len(hostInfo.Uname) > 1, true)
  assert.Equal(t, len(hostInfo.UnameA) > 1, true)
  assert.Equal(t, len(hostInfo.IpAddresses) > 1, true)
  assert.Equal(t, len(hostInfo.Hostnames) > 1, true)
}
