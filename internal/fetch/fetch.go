package fetch

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	utils "github.com/AlexandreLeitao/GoProcView/internal/utils"
)

// MemoryStats holds parsed memory info
type MemoryStats struct {
	Total float64
	Free  float64
	Used  float64
}

type ProcData struct {
	CpuUsage     []float64
	CpuAverage   uint64
	Memory       MemoryStats
	LoadAverages [3]float64
	Uptime       float64
	ProcessInfo  map[int]processData
}

type processData struct {
	Pid        uint64
	Name       string
	CpuPercent float64
}

func FetchProcData(pId uint64) ProcData {
	var data ProcData
	memStats, _ := GetMemoryStats()
	if memStats != nil {
		data.Memory = *memStats
	}
	cpuUsages, cpuAvg, _ := GetCPUStats()
	data.CpuUsage = cpuUsages
	data.CpuAverage = cpuAvg
	loadAvgs, _ := GetLoadAverages()
	data.LoadAverages = loadAvgs
	uptime, _ := GetUptime()
	data.Uptime = uptime

	if pId != 0 {
		procInfo, _ := GetProcessInfo(int(pId))
		if procInfo != nil {
			data.ProcessInfo = make(map[int]processData)
			data.ProcessInfo[int(pId)] = *procInfo
		}
	} else {
		topProcs, _ := GetTopProcessesByCpu(3)
		data.ProcessInfo = make(map[int]processData)
		for _, proc := range topProcs {
			data.ProcessInfo[int(proc.Pid)] = proc
		}
	}

	return data
}

// GetMemoryStats reads /proc/meminfo and parses memory usage
func GetMemoryStats() (*MemoryStats, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var total, free float64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case "MemTotal:":
			val, _ := strconv.ParseFloat(fields[1], 64)
			total = val
		case "MemFree:":
			val, _ := strconv.ParseFloat(fields[1], 64)
			free = val
		}
	}

	used := total - free

	return &MemoryStats{
		Total: total,
		Free:  free,
		Used:  used,
	}, nil
}

// GetCPUStats reads /proc/stat and parses cpu usage
func GetCPUStats() ([]float64, uint64, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()
	var cpuUsages []float64
	var cpuAverage uint64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu") {
			fields := strings.Fields(line)
			if len(fields) < 5 {
				continue
			}
			var total, idle uint64
			for i, field := range fields[1:] {
				value, err := strconv.ParseUint(field, 10, 64)
				if err != nil {
					return nil, 0, err
				}
				total += value
				if i == 3 { // idle is the 4th field
					idle = value
				}
			}
			usage := float64(total-idle) / float64(total) * 100
			cpuUsages = append(cpuUsages, usage)
			cpuAverage += total - idle
		}
	}
	if len(cpuUsages) > 0 {
		var total float64
		for _, usage := range cpuUsages {
			total += usage
		}
		cpuAverage = uint64(total / float64(len(cpuUsages)))
	}

	return cpuUsages, cpuAverage, nil
}

// GetLoadAverages reads /proc/loadavg and parses load averages
func GetLoadAverages() ([3]float64, error) {
	var loadAverages [3]float64

	file, err := os.Open("/proc/loadavg")
	if err != nil {
		return loadAverages, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 3 {
			return loadAverages, nil
		}
		for i := 0; i < 3; i++ {
			value, err := strconv.ParseFloat(fields[i], 64)
			if err != nil {
				return loadAverages, err
			}
			loadAverages[i] = value
		}
	}

	return loadAverages, nil
}

// GetUptime reads /proc/uptime and parses system uptime
func GetUptime() (float64, error) {
	file, err := os.Open("/proc/uptime")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 1 {
			return 0, nil
		}
		val, err := strconv.ParseFloat(fields[0], 64)
		if err != nil {
			return 0, err
		}
		return val, nil
	}

	return 0, nil
}

// GetProcessInfo reads /proc/[pid]/status and parses process info
func GetProcessInfo(pid int) (*processData, error) {
	file, err := os.Open("/proc/" + strconv.Itoa(pid) + "/status")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var name string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) < 2 {
			continue
		}

		if fields[0] == "Name:" {
			name = fields[1]
			break
		}
	}

	return &processData{
		Pid:  uint64(pid),
		Name: name,
	}, nil
}

func GetTopProcessesByCpu(n int) ([]processData, error) {
	type procCpu struct {
		data    processData
		cputime uint64
	}
	var procs []procCpu

	procDir, err := os.Open("/proc")
	if err != nil {
		return nil, err
	}
	defer procDir.Close()

	names, err := procDir.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	// First sample of process cpu times and total cpu time
	procTimes1 := make(map[int]uint64)
	totalCpu1, err := readTotalCpuTime()
	if err != nil {
		return nil, err
	}
	for _, name := range names {
		pid, err := strconv.Atoi(name)
		if err != nil {
			continue
		}
		utime, stime, err := readProcCpuTimes(pid)
		if err != nil {
			continue
		}
		procTimes1[pid] = utime + stime
	}
	// Sleep for sample interval (100ms)
	time.Sleep(100 * time.Millisecond)

	// Second sample
	procTimes2 := make(map[int]uint64)
	totalCpu2, err := readTotalCpuTime()
	if err != nil {
		return nil, err
	}
	for _, name := range names {
		pid, err := strconv.Atoi(name)
		if err != nil {
			continue
		}
		utime, stime, err := readProcCpuTimes(pid)
		if err != nil {
			continue
		}
		procTimes2[pid] = utime + stime
	}

	totalDelta := float64(totalCpu2 - totalCpu1)
	if totalDelta == 0 {
		totalDelta = 1 // avoid div by zero
	}

	// Build slice with cpu percent
	for pid, t2 := range procTimes2 {
		t1, ok := procTimes1[pid]
		if !ok {
			continue
		}
		delta := t2 - t1
		cpuPercent := float64(delta) / totalDelta * 100
		info, err := GetProcessInfo(pid)
		if err != nil {
			continue
		}
		procs = append(procs, procCpu{
			data: processData{
				Pid:        uint64(pid),
				Name:       info.Name,
				CpuPercent: cpuPercent,
			},
			cputime: delta,
		})
	}

	// Sort by cpu percent descending
	sort.Slice(procs, func(i, j int) bool {
		return procs[i].data.CpuPercent > procs[j].data.CpuPercent
	})

	var top []processData
	for i := 0; i < n && i < len(procs); i++ {
		top = append(top, procs[i].data)
	}
	return top, nil
}

func GetHardwareInfo() (map[string]string, error) {
	hwInfo := make(map[string]string)

	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "model name" {
			hwInfo["CPU Name"] = value
			break
		}
	}

	file, err = os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case "MemTotal:":
			val, _ := strconv.ParseFloat(fields[1], 64)
			hwInfo["Total Memory (Gb)"] = fmt.Sprintf("%.2f Gb", utils.CalculateKbtoGb(val))
		}
	}

	cmd := exec.Command("lsblk", "-o", "NAME,SIZE,MODEL", "-dn")
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(stdout), "\n")
	var disks []string
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			disks = append(disks, fmt.Sprintf("%s (%s) - %s", fields[0], fields[2], fields[1]))
		} else if len(fields) == 2 {
			disks = append(disks, fmt.Sprintf("%s - %s", fields[0], fields[1]))
		}
	}
	hwInfo["Disks"] = strings.Join(disks, ", ")

	return hwInfo, nil
}

// Helpers

func readProcCpuTimes(pid int) (utime, stime uint64, err error) {
	statFile := "/proc/" + strconv.Itoa(pid) + "/stat"
	f, err := os.Open(statFile)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		return 0, 0, os.ErrInvalid
	}
	fields := strings.Fields(scanner.Text())
	if len(fields) < 17 {
		return 0, 0, os.ErrInvalid
	}
	utime, _ = strconv.ParseUint(fields[13], 10, 64)
	stime, _ = strconv.ParseUint(fields[14], 10, 64)
	return utime, stime, nil
}

func readTotalCpuTime() (uint64, error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return 0, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)
			var total uint64
			for _, s := range fields[1:] {
				val, _ := strconv.ParseUint(s, 10, 64)
				total += val
			}
			return total, nil
		}
	}
	return 0, os.ErrInvalid
}
