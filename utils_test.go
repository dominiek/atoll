package main

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestIntervalToSeconds(t *testing.T) {
  var seconds float64
  var err error
  seconds, err = IntervalToSeconds("5s")
  assert.Equal(t, err, nil)
  assert.Equal(t, seconds, 5.0)
  seconds, err = IntervalToSeconds("5seconds")
  assert.Equal(t, err, nil)
  assert.Equal(t, seconds, 5.0)
  seconds, err = IntervalToSeconds("10minutes")
  assert.Equal(t, err, nil)
  assert.Equal(t, seconds, 10.0*60.0)
  seconds, err = IntervalToSeconds("1h")
  assert.Equal(t, err, nil)
  assert.Equal(t, seconds, 60.0*60.0)
  seconds, err = IntervalToSeconds("2bla")
  assert.NotEqual(t, err, nil)
  seconds, err = IntervalToSeconds("bla")
  assert.NotEqual(t, err, nil)
  seconds, err = IntervalToSeconds("0.5s")
  assert.Equal(t, err, nil)
  assert.Equal(t, seconds, 0.5)
}
