package main

import (
  "testing"
  "github.com/stretchr/testify/assert"
  "net/http/httptest"
  "net/http"
  "fmt"
  "io/ioutil"
)

func TestReporterWithNetstat(t *testing.T) {
  config := Config{};

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
  assert.Equal(t, err, nil)
}

func TestReporterGetHostInfo(t *testing.T) {
  config := Config{};
  config.HOSTNAME = "0.localhost"
  reporter := Reporter{&config, Netstat{config: &config}, true, "netstat", "http://localhost:47011"};
  hostInfo := reporter.GetHostInfo();
  assert.Equal(t, len(hostInfo.hostnames), 1)
}