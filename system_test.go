package main

import (
  "testing"
  "github.com/stretchr/testify/assert"
  "encoding/json"
  //"fmt"
)

func TestSystemMonitor(t *testing.T) {
  config := Config{};
  var err error;
  err = config.LoadFile("./atoll.yml")
  if err != nil {
    t.Fatalf("Did not expect error %v", err)
  }
  system := System{config: &config}
  var info Info;
  info, err = system.Monitor();
  if err != nil {
    t.Fatalf("Did not expect error %v", err)
  }
  result := info.(SystemInfo)
  assert.EqualValues(t, true, result.Cpu.NumCores > 0)
  assert.EqualValues(t, true, result.Memory.TotalKb > 0)
  assert.EqualValues(t, true, result.Memory.FreeKb > 0)
  data, err := json.Marshal(info)
  if err != nil {
    t.Fatalf("Did not expect error %v", err)
  }
  t.Logf("JSON %s", data)
}

var darwinSysctlOutput = map[string]string{
  "hw.ncpu": "4\n",
  "vm.loadavg": "{ 1.95 2.03 2.24 }\n",
  "hw.memsize": "8589934592\n",
}

func TestSystemParseSysctl(t *testing.T) {
  system := System{}
  info := SystemCpuInfo{}
  err := system.parseNumCpuSysctl(&info, darwinSysctlOutput["hw.ncpu"])
  if err != nil {
    t.Fatalf("Did not expect error %v", err)
  }
  assert.EqualValues(t, info.NumCores, 4)
  err = system.parseLoadAverageSysctl(&info, darwinSysctlOutput["vm.loadavg"])
  if err != nil {
    t.Fatalf("Did not expect error %v", err)
  }
  assert.EqualValues(t, 1.95, info.MinuteAverage)
  assert.EqualValues(t, 2.03, info.FiveMinuteAverage)
  assert.EqualValues(t, 2.24, info.FifteenMinuteAverage)
  memInfo := SystemMemoryInfo{}
  err = system.parseMemsizeSysctl(&memInfo, darwinSysctlOutput["hw.memsize"])
  if err != nil {
    t.Fatalf("Did not expect error %v", err)
  }
  assert.EqualValues(t, 8589934592/1024, memInfo.TotalKb)
}

var linuxProcOutput = map[string]string{
  "/proc/meminfo": `MemTotal:        2050412 kB
MemFree:         1290092 kB
MemAvailable:    1706844 kB
Buffers:           75988 kB
Cached:           434844 kB
SwapCached:         2820 kB
Active:           407008 kB
Inactive:         181656 kB
Active(anon):      79872 kB
Inactive(anon):    14496 kB
Active(file):     327136 kB
Inactive(file):   167160 kB
Unevictable:           0 kB
Mlocked:               0 kB
SwapTotal:       1461256 kB
SwapFree:        1346204 kB
Dirty:                 0 kB
Writeback:             0 kB
AnonPages:         75092 kB
Mapped:            37488 kB
Shmem:             16536 kB
Slab:              86432 kB
SReclaimable:      70160 kB
SUnreclaim:        16272 kB
KernelStack:        2336 kB
PageTables:         1472 kB
NFS_Unstable:          0 kB
Bounce:                0 kB
WritebackTmp:          0 kB
CommitLimit:     2486460 kB
Committed_AS:     714808 kB
VmallocTotal:   34359738367 kB
VmallocUsed:        9752 kB
VmallocChunk:   34359697344 kB
AnonHugePages:     57344 kB
HugePages_Total:       0
HugePages_Free:        0
HugePages_Rsvd:        0
HugePages_Surp:        0
Hugepagesize:       2048 kB
DirectMap4k:       44992 kB
DirectMap2M:     2052096 kB
`,
  "/proc/loadavg": "0.00 0.01 0.05 1/143 2058\n",
  "/proc/cpuinfo": `processor	: 0
vendor_id	: GenuineIntel
cpu family	: 6
model		: 61
model name	: Intel(R) Core(TM) M-5Y31 CPU @ 0.90GHz
stepping	: 4
microcode	: 0x19
cpu MHz		: 1099.999
cache size	: 4096 KB
physical id	: 0
siblings	: 4
core id		: 0
cpu cores	: 4
apicid		: 0
initial apicid	: 0
fpu		: yes
fpu_exception	: yes
cpuid level	: 20
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx rdtscp lm constant_tsc rep_good nopl xtopology nonstop_tsc pni ssse3 sse4_1 sse4_2 hypervisor lahf_lm
bugs		:
bogomips	: 2199.99
clflush size	: 64
cache_alignment	: 64
address sizes	: 39 bits physical, 48 bits virtual
power management:

processor	: 1
vendor_id	: GenuineIntel
cpu family	: 6
model		: 61
model name	: Intel(R) Core(TM) M-5Y31 CPU @ 0.90GHz
stepping	: 4
microcode	: 0x19
cpu MHz		: 1099.999
cache size	: 4096 KB
physical id	: 0
siblings	: 4
core id		: 1
cpu cores	: 4
apicid		: 1
initial apicid	: 1
fpu		: yes
fpu_exception	: yes
cpuid level	: 20
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx rdtscp lm constant_tsc rep_good nopl xtopology nonstop_tsc pni ssse3 sse4_1 sse4_2 hypervisor lahf_lm
bugs		:
bogomips	: 2199.99
clflush size	: 64
cache_alignment	: 64
address sizes	: 39 bits physical, 48 bits virtual
power management:

processor	: 2
vendor_id	: GenuineIntel
cpu family	: 6
model		: 61
model name	: Intel(R) Core(TM) M-5Y31 CPU @ 0.90GHz
stepping	: 4
microcode	: 0x19
cpu MHz		: 1099.999
cache size	: 4096 KB
physical id	: 0
siblings	: 4
core id		: 2
cpu cores	: 4
apicid		: 2
initial apicid	: 2
fpu		: yes
fpu_exception	: yes
cpuid level	: 20
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx rdtscp lm constant_tsc rep_good nopl xtopology nonstop_tsc pni ssse3 sse4_1 sse4_2 hypervisor lahf_lm
bugs		:
bogomips	: 2199.99
clflush size	: 64
cache_alignment	: 64
address sizes	: 39 bits physical, 48 bits virtual
power management:

processor	: 3
vendor_id	: GenuineIntel
cpu family	: 6
model		: 61
model name	: Intel(R) Core(TM) M-5Y31 CPU @ 0.90GHz
stepping	: 4
microcode	: 0x19
cpu MHz		: 1099.999
cache size	: 4096 KB
physical id	: 0
siblings	: 4
core id		: 3
cpu cores	: 4
apicid		: 3
initial apicid	: 3
fpu		: yes
fpu_exception	: yes
cpuid level	: 20
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx rdtscp lm constant_tsc rep_good nopl xtopology nonstop_tsc pni ssse3 sse4_1 sse4_2 hypervisor lahf_lm
bugs		:
bogomips	: 2199.99
clflush size	: 64
cache_alignment	: 64
address sizes	: 39 bits physical, 48 bits virtual
power management:
`,
}

func TestSystemParseProc(t *testing.T) {
  system := System{}
  info := SystemCpuInfo{}
  err := system.parseCpuinfoProc(&info, linuxProcOutput["/proc/cpuinfo"])
  if err != nil {
    t.Fatalf("Did not expect error %v", err)
  }
  assert.EqualValues(t, 4, info.NumCores)
  err = system.parseLoadAverageProc(&info, linuxProcOutput["/proc/loadavg"])
  if err != nil {
    t.Fatalf("Did not expect error %v", err)
  }
  assert.EqualValues(t, 0.00, info.MinuteAverage)
  assert.EqualValues(t, 0.01, info.FiveMinuteAverage)
  assert.EqualValues(t, 0.05, info.FifteenMinuteAverage)
  memInfo := SystemMemoryInfo{}
  err = system.parseMeminfoProc(&memInfo, linuxProcOutput["/proc/meminfo"])
  if err != nil {
    t.Fatalf("Did not expect error %v", err)
  }
  assert.EqualValues(t, 2050412, memInfo.TotalKb)
  assert.EqualValues(t, 1290092, memInfo.FreeKb)
}

const darwinVmstatOutput = `Mach Virtual Memory Statistics: (page size of 4096 bytes)
Pages free:                               96872.
Pages active:                            640667.
Pages inactive:                          203733.
Pages speculative:                         5172.
Pages throttled:                              0.
Pages wired down:                        840711.
Pages purgeable:                         114491.
"Translation faults":                1662495560.
Pages copy-on-write:                   30324937.
Pages zero filled:                    443026576.
Pages reactivated:                      8849762.
Pages purged:                          28834888.
File-backed pages:                       156440.
Anonymous pages:                         693132.
Pages stored in compressor:             1368944.
Pages occupied by compressor:            309683.
Decompressions:                        32534922.
Compressions:                          45427471.
Pageins:                               78086671.
Pageouts:                               1961153.
Swapins:                               22136101.
Swapouts:                              23355472.
`;

func TestSystemMemoryParseVmstat(t *testing.T) {
  system := System{}
  info := SystemMemoryInfo{}
  info.TotalKb = 8 * 1024 * 1024 * 1024
  err := system.parseVmstat(&info, darwinVmstatOutput)
  if err != nil {
    t.Fatalf("Did not expect error %v", err)
  }
  assert.EqualValues(t, 0x1ffa59578, info.FreeKb)
}
