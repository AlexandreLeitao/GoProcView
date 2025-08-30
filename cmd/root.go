package cmd

import (
	"flag"
	"fmt"
	"time"

	fetch "github.com/AlexandreLeitao/GoProcView/internal/fetch"
	"github.com/AlexandreLeitao/GoProcView/internal/render"
)

func Execute() {
	// Define a command-line flag
	help := flag.Bool("h", false, "Help for GoProcView")
	command := flag.String("c", "procview", "The command to be executed by GoProcView")
	pId := flag.Uint64("p", 0, "Process ID to fetch info for (0 for none)")
	mode := flag.String("m", "graphic", "Mode of operation: simple or graphic")
	watch := flag.Bool("w", false, "Watch mode - refresh every 1 seconds")
	flag.Parse()

	if *help == false {
		switch *command {
		case "procview":
			if *watch {
				for {
					procData := fetch.FetchProcData(*pId)
					if *mode == "graphic" {
						render.GraphicRender(procData)
					} else if *mode == "simple" {
						render.SimpleRender(procData)
					}
					fmt.Println("--------------------------------------------------")
					time.Sleep(1 * time.Second)
				}
			}

			procData := fetch.FetchProcData(*pId)
			if *mode == "graphic" {
				render.GraphicRender(procData)
			} else if *mode == "simple" {
				render.SimpleRender(procData)
			}
		case "sysinfo":
			hwData, _ := fetch.GetHardwareInfo()
			render.RenderHardwareInfo(hwData)
		case "help":
			*help = true
		default:
			fmt.Printf("Unknown command: %s\n", *command)
		}
	}
	if *help {
		fmt.Println("Usage of GoProcView:")
		flag.PrintDefaults()
		return
	}
}
