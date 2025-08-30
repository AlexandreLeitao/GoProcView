package render

import (
	"fmt"
	"strings"

	fetch "github.com/AlexandreLeitao/GoProcView/internal/fetch"
	utils "github.com/AlexandreLeitao/GoProcView/internal/utils"
)

func SimpleRender(data fetch.ProcData) {

	cpuRounded := []string{}

	for _, cpu := range data.CpuUsage {
		cpuRounded = append(cpuRounded, fmt.Sprintf("%.2f", cpu)+"%")
	}
	fmt.Println("CPU Usage: ", cpuRounded)
	fmt.Println("CPU Average: ", data.CpuAverage, "%")
	fmt.Println("Memory:", fmt.Sprintf("%.2f", utils.CalculateKbtoGb(data.Memory.Used)), "/", fmt.Sprintf("%.2f", utils.CalculateKbtoGb(data.Memory.Total)), "Gb")
	fmt.Println("Memory Free: ", fmt.Sprintf("%.2f", utils.CalculateKbtoGb(data.Memory.Free)), "Gb")
	fmt.Println("Load Averages: ", data.LoadAverages)

	hours := int(data.Uptime) / 3600
	mins := (int(data.Uptime) % 3600) / 60
	secs := int(data.Uptime) % 60
	fmt.Printf("Uptime: %02d:%02d:%02d (hh:mm:ss)\n", hours, mins, secs)

	if len(data.ProcessInfo) > 0 {
		for pid, proc := range data.ProcessInfo {
			fmt.Printf("Process ID: %d, Name: %s, Cpu: %s", pid, proc.Name, fmt.Sprintf("%.2f", proc.CpuPercent)+"%\n")
		}
	}
}

// Color constants (ANSI escape codes)
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Green  = "\033[32m"
)

// usageBar renders a bar with colors based on usage %
func usageBar(usage float64, length int) string {
	// Pick color
	color := Green
	switch {
	case usage >= 80:
		color = Red
	case usage >= 50:
		color = Yellow
	}

	// Calculate filled vs empty slots
	filled := int((usage / 100.0) * float64(length))
	if filled > length {
		filled = length
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", length-filled)

	return fmt.Sprintf("%s%s%s %5.1f%%", color, bar, Reset, usage)
}

func GraphicRender(data fetch.ProcData) {

	fmt.Println("CPU Average: ", usageBar(float64(data.CpuAverage), 20))
	for i, cpu := range data.CpuUsage {
		separator := ""
		if i < 10 {
			separator = ": "
		} else {
			separator = ":"
		}
		fmt.Println("Core", i, separator, usageBar(cpu, 20))
	}
	fmt.Println("Memory:", usageBar((data.Memory.Used/data.Memory.Total)*100, 20), fmt.Sprintf("%.2f", utils.CalculateKbtoGb(data.Memory.Used)), "/", fmt.Sprintf("%.2f", utils.CalculateKbtoGb(data.Memory.Total)), "Gb")

	hours := int(data.Uptime) / 3600
	mins := (int(data.Uptime) % 3600) / 60
	secs := int(data.Uptime) % 60
	fmt.Printf("Uptime: %02d:%02d:%02d (hh:mm:ss)\n", hours, mins, secs)

	if len(data.ProcessInfo) > 0 {
		for pid, proc := range data.ProcessInfo {
			fmt.Printf("Process ID: %d, Name: %s, Cpu: %s", pid, proc.Name, fmt.Sprintf("%.2f", proc.CpuPercent)+"%\n")
		}
	}
}

func RenderHardwareInfo(hwInfo map[string]string) {
	fmt.Println("Hardware Information:")
	for key, value := range hwInfo {
		fmt.Printf("%-20s : %s\n", key, value)
	}
}
