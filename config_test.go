package main

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestLoadFile(t *testing.T) {
  config := Config{};
  err := config.LoadFile("./atoll.yml")
  if err != nil {
    t.Fatalf("Did not expect error %v", err)
  }
  assert.Equal(t, "api.atoll.io", config.Publish.Host)
}

func TestToJSON(t *testing.T) {
  config := Config{};
  var err error;
  err = config.LoadFile("./atoll.yml")
  if err != nil {
    t.Fatalf("Did not expect error %v", err)
  }

  var data []byte;
  data, err = config.ToJSON()
  assert.Contains(t, string(data), "\"publish\":")
}
