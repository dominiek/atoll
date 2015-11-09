
package main

import (
  "os/exec"
  "errors"
  "strings"
  "strconv"
  "encoding/json"
  "regexp"
  //"log"
)

type SystemCpuInfo struct {
  NumCores              uint16  `json:"numCores"`;
  MinuteAverage         float64 `json:"minuteAverage"`;
  FiveMinuteAverage     float64 `json:"fiveMinuteAverage"`;
  FifteenMinuteAverage  float64 `json:"fifteenMinuteAverage"`;
}

type SystemMemoryInfo struct {
  TotalKb uint64 `json:"totalKb"`;
  FreeKb  uint64 `json:"freeKb"`;
}

type SystemInfo struct {
  Cpu SystemCpuInfo `json:"cpu"`;
  Memory SystemMemoryInfo `json:"memory"`;
}

type System struct {
  config *Config;
  currentInfo SystemInfo;
};

func (this SystemInfo) Encode() ([]byte, error) {
  return json.Marshal(this)
}

func (this SystemInfo) GetType() (string) {
  return "system";
}

func (this System) Monitor() (Info, error) {
  this.ClearResults();
  uname, _ := this.getUname();
  var data string;
  var err error;
  switch (uname) {
    case "Darwin":
      data, err = this.execCommand("sysctl -n hw.ncpu")
      if err != nil {
        return this.currentInfo, nil
      }
      err = this.parseNumCpuSysctl(&this.currentInfo.Cpu, data)
      if err != nil {
        return this.currentInfo, nil
      }
      data, err = this.execCommand("sysctl -n vm.loadavg")
      if err != nil {
        return this.currentInfo, nil
      }
      err = this.parseLoadAverageSysctl(&this.currentInfo.Cpu, data)
      if err != nil {
        return this.currentInfo, nil
      }
      data, err = this.execCommand("sysctl -n hw.memsize")
      if err != nil {
        return this.currentInfo, nil
      }
      err = this.parseMemsizeSysctl(&this.currentInfo.Memory, data)
      if err != nil {
        return this.currentInfo, nil
      }
      data, err = this.execCommand("vm_stat")
      if err != nil {
        return this.currentInfo, nil
      }
      err = this.parseVmstat(&this.currentInfo.Memory, data)
      if err != nil {
        return this.currentInfo, nil
      }
    break
    default:
      data, err = this.execCommand("cat /proc/cpuinfo")
      if err != nil {
        return this.currentInfo, nil
      }
      err = this.parseCpuinfoProc(&this.currentInfo.Cpu, data)
      if err != nil {
        return this.currentInfo, nil
      }
      data, err = this.execCommand("cat /proc/loadavg")
      if err != nil {
        return this.currentInfo, nil
      }
      err = this.parseLoadAverageProc(&this.currentInfo.Cpu, data)
      if err != nil {
        return this.currentInfo, nil
      }
      data, err = this.execCommand("cat /proc/meminfo")
      if err != nil {
        return this.currentInfo, nil
      }
      err = this.parseMeminfoProc(&this.currentInfo.Memory, data)
      if err != nil {
        return this.currentInfo, nil
      }
    break
  }
  info := Info(this.currentInfo);
  return info, err
}

func (this System) ClearResults() () {
  this.currentInfo = SystemInfo{
    SystemCpuInfo{},
    SystemMemoryInfo{},
  }
}

func (this System) execCommand(commandStr string) (string, error) {
  command := strings.Split(commandStr, " ");
  cmd := exec.Command(command[0], command[1:]...)
  data, err := cmd.Output()
  if err != nil {
    return "", err
  }
  return strings.TrimSpace(string(data)), nil
}

func (this System) getUname() (string, error) {
  data, err := exec.Command("uname").Output();
  if err != nil {
    return "", err
  }
  return strings.TrimSpace(string(data)), nil
}

func (this System) parseNumCpuSysctl(cpuInfo *SystemCpuInfo, data string) (error) {
  data = strings.TrimSpace(data)
  cpuCores, err := strconv.ParseUint(data, 10, 16)
  if err != nil {
    return err
  }
  cpuInfo.NumCores = uint16(cpuCores)
  return nil
}

func (this System) parseMemsizeSysctl(info *SystemMemoryInfo, data string) (error) {
  data = strings.TrimSpace(data)
  total, err := strconv.ParseUint(data, 10, 64)
  if err != nil {
    return err
  }
  info.TotalKb = total / 1024;
  return nil
}

func (this System) parseLoadAverageSysctl(cpuInfo *SystemCpuInfo, data string) (error) {
  data = strings.Replace(data, "{", "", -1)
  data = strings.Replace(data, "}", "", -1)
  data = strings.Replace(data, "\t", " ", -1)
  data = strings.Replace(data, "  ", " ", -1)
  data = strings.TrimSpace(data)
  averages := strings.Split(data, " ")
  var err error;
  if (len(averages) < 3) {
    return errors.New("Invalid load average format")
  }
  cpuInfo.MinuteAverage, err = strconv.ParseFloat(averages[0], 64)
  if err != nil {
    return err
  }
  cpuInfo.FiveMinuteAverage, err = strconv.ParseFloat(averages[1], 64)
  if err != nil {
    return err
  }
  cpuInfo.FifteenMinuteAverage, err = strconv.ParseFloat(averages[2], 64)
  if err != nil {
    return err
  }
  return nil
}

func (this System) parseCpuinfoProc(cpuInfo *SystemCpuInfo, data string) (error) {
  processors := regexp.MustCompile("processor\\s*\\:").Split(data, -1)
  cpuInfo.NumCores = uint16(len(processors) - 1)
  return nil
}

func (this System) parseLoadAverageProc(cpuInfo *SystemCpuInfo, data string) (error) {
  data = strings.Replace(data, "\t", " ", -1)
  data = strings.Replace(data, "  ", " ", -1)
  data = strings.TrimSpace(data)
  averages := strings.Split(data, " ")
  var err error;
  cpuInfo.MinuteAverage, err = strconv.ParseFloat(averages[0], 64)
  if err != nil {
    return err
  }
  cpuInfo.FiveMinuteAverage, err = strconv.ParseFloat(averages[1], 64)
  if err != nil {
    return err
  }
  cpuInfo.FifteenMinuteAverage, err = strconv.ParseFloat(averages[2], 64)
  if err != nil {
    return err
  }
  return nil
}

func (this System) parseMeminfoProc(memoryInfo *SystemMemoryInfo, data string) (error) {
  data = strings.TrimSpace(data);
  lines := strings.Split(data, "\n");
  for i := 0; len(lines) > i; i++ {
    item := regexp.MustCompile(":\\s+").Split(lines[i], -1)
    if len(item) < 2 {
      continue;
    }
    key := strings.ToLower(item[0])
    value := regexp.MustCompile("\\s+").Split(item[1], 2)
    if len(value) < 2 {
      continue;
    }
    if strings.ToLower(value[1]) != "kb" {
      return errors.New("Unsupported unit encountered in /proc/meminfo")
    }
    numKb, _ := strconv.ParseUint(value[0], 10, 64)
    if key == "memtotal" {
      memoryInfo.TotalKb = numKb
    }
    if key == "memfree" {
      memoryInfo.FreeKb = numKb
    }
  }
  return nil
}

func (this System) parseVmstat(memoryInfo *SystemMemoryInfo, data string) (error) {
  data = strings.TrimSpace(data);
  lines := strings.Split(data, "\n");
  var totalUsageKb uint64 = 0
  for i := 1; len(lines) > i; i++ {
    item := regexp.MustCompile(":\\s+").Split(lines[i], -1)
    if len(item) < 2 {
      continue;
    }
    key := strings.ToLower(item[0])
    value := strings.Replace(item[1], ".", "", 1)
    numPages, _ := strconv.ParseUint(value, 10, 64)
    numKb := ((numPages * 4096) / 1024)
    if key == "pages active" {
      totalUsageKb = totalUsageKb + numKb
    }
    if key == "pages wired down" {
      totalUsageKb = totalUsageKb + numKb
    }
  }
  if (totalUsageKb > memoryInfo.TotalKb) {
    return nil
  }
  memoryInfo.FreeKb = memoryInfo.TotalKb - totalUsageKb
  return nil
}
