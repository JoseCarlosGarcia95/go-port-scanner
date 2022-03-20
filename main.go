package main

import (
	"os"
	"strconv"

	"github.com/JoseCarlosGarcia95/go-port-scanner/portscanner"
)

func main() {
	// Create a basic CLI that call the portscanner.PortRange function.
	// The CLI will print the ports that are open.
	//
	// Usage:
	//   go run main.go <host> <protocol> <start port> <end port> <workers>
	//

	args := os.Args[1:]

	if len(args) != 5 {
		panic("Invalid number of arguments")
	}

	host := args[0]
	protocol := args[1]
	start, err := strconv.Atoi(args[2])
	if err != nil {
		panic("Invalid start port")
	}

	end, err := strconv.Atoi(args[3])
	if err != nil {
		panic("Invalid end port")
	}

	workers, err := strconv.Atoi(args[4])
	if err != nil {
		panic("Invalid workers")
	}

	ports := portscanner.PortRange(host, protocol, uint32(start), uint32(end), uint32(workers))

	for _, port := range ports {
		servicename := portscanner.Port2Service(host, protocol, port, true)
		println(port, "-", servicename)
	}
}
