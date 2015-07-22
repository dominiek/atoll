
package main

import(
  "regexp"
  "errors"
  "strconv"
);

type Interval struct {
  name, alias string;
}

var secondsIntervalMap = map[int]Interval{
  (1): Interval{
    "seconds", "s",
  },
  (60): Interval{
    "minutes", "m",
  },
  (60*60): Interval{
    "hours",   "h",
  },
  (24*60*60): Interval{
    "days",    "d",
  },
  (7*24*60*60): Interval{
    "weeks",   "w",
  },
}

func IntervalToSeconds(interval string) (int, error) {
  re := regexp.MustCompile("^([0-9]+)([a-z]+)$");
  result := re.FindAllStringSubmatch(interval, 2);
  if len(result) > 0 {
    if (len(result[0]) != 3) {
      return 0, errors.New("Invalid format for interval, expected <number><unit> (E.g. 5s, 1h)");
    }
    if len(result[0][1]) == 0 {
      return 0, errors.New("Invalid format for interval, expected interval to start with a number");
    }
    number, err := strconv.Atoi(result[0][1]);
    unit := result[0][2];
    if err != nil {
      return 0, err
    }
    seconds := 0
    for intervalSeconds, interval := range secondsIntervalMap {
      if unit == interval.name || unit == interval.alias {
        seconds = intervalSeconds
      }
    }
    if seconds == 0 {
      return 0, errors.New("Invalid unit for interval");
    }
    return (number * seconds), nil
  }
  return 0, errors.New("Invalid format for interval, expected <number><unit> (E.g. 5s, 1h)");
}
