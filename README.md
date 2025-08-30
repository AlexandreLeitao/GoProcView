# GoProcView

GoProcView is a mini system monitor CLI tool written in Go. It provides real-time information about your system’s CPU, memory, load averages, uptime, and process details.

## Features

- **CPU Usage**: Displays per-core and average CPU usage percentages.
- **Memory Stats**: Shows total, free, and used memory.
- **Load Averages**: Reports 1, 5, and 15 minute system load averages.
- **Uptime**: Shows how long the system has been running.
- **Process Info**: Fetches details for a specific process by PID.

## Usage

```sh
go run .
```

```sh
go run . -h

Usage of GoProcView:
  -c string
        The command to be executed by GoProcView (default "procview")
  -h    Help for GoProcView
  -m string
        Mode of operation: simple or graphic (default "graphic")
  -p uint
        Process ID to fetch info for (0 for none)
  -w    Watch mode - refresh every 1 seconds
```

## Command List:
- procview - Info from /proc 
    - -m simple - Simple data for procview
    - -m graphic - Graphical representation of the data
- sysinfo - Info about system.
- help - same as -h



## Project Structure

- `main.go`: Entry point, calls the CLI logic in `cmd/root.go`.
- `cmd/root.go`: CLI command definitions and execution logic.
- `internal/fetch/fetch.go`: System data fetching (CPU, memory, load, uptime, process info).
- `internal/render/render.go`: Data rendering.

## Example Output

```
CPU Usage:  [8.92% 9.25% 8.43% 8.39% 8.17% 8.26% 8.25% 7.17% 6.50% 11.30% 11.10% 10.53% 9.72%]
CPU Average:  8 %
Memory: 11.26 / 15.42 Gb
Memory Free:  4.17 Gb
Load Averages:  [0.86 1.08 1.21]
Uptime: 03:24:40 (hh:mm:ss)
Process ID: 1889, Name: code, Cpu: 0.80%
Process ID: 156957, Name: GoProcView, Cpu: 0.80%
Process ID: 702, Name: Firefox, Cpu: 0.80%
```

```
CPU Average:  █░░░░░░░░░░░░░░░░░░░   8.0%
Core 0 :  █░░░░░░░░░░░░░░░░░░░   9.0%
Core 1 :  █░░░░░░░░░░░░░░░░░░░   9.3%
Core 2 :  █░░░░░░░░░░░░░░░░░░░   8.5%
Core 3 :  █░░░░░░░░░░░░░░░░░░░   8.4%
Core 4 :  █░░░░░░░░░░░░░░░░░░░   8.2%
Core 5 :  █░░░░░░░░░░░░░░░░░░░   8.3%
Core 6 :  █░░░░░░░░░░░░░░░░░░░   8.3%
Core 7 :  █░░░░░░░░░░░░░░░░░░░   7.2%
Core 8 :  █░░░░░░░░░░░░░░░░░░░   6.5%
Core 9 :  ██░░░░░░░░░░░░░░░░░░  11.3%
Core 10 : ██░░░░░░░░░░░░░░░░░░  11.1%
Core 11 : ██░░░░░░░░░░░░░░░░░░  10.6%
Core 12 : █░░░░░░░░░░░░░░░░░░░   9.8%
Memory: ██████████████░░░░░░  72.1% 11.12 / 15.42 Gb
Uptime: 03:15:21 (hh:mm:ss)
Process ID: 154497, Name: GoProcView, Cpu: 1.59%
Process ID: 1889, Name: code, Cpu: 0.79%
Process ID: 142, Name: Firefox, Cpu: 0.00%
```
## Requirements

- Go 1.18+
- Linux (uses `/proc` filesystem)

## License

MIT
