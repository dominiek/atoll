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
  assert.Equal(t, "127.0.0.1", config.SERVER.BIND)
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
  assert.Contains(t, string(data), "\"BIND\":")
}
