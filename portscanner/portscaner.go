package portscanner

import (
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/JoseCarlosGarcia95/go-port-scanner/portscanner/services"
	"github.com/thediveo/netdb"
)

// IsPortOpen checks if a port is open or not.
func IsPortOpen(host, protocol string, port uint32) bool {
	conn, err := net.DialTimeout(protocol, host+":"+strconv.Itoa(int(port)), 10 * time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// PortRange returns a slice of ports that are open.
func PortRange(host, protocol string, start, end, workers uint32) []uint32 {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var ports = make(chan uint32, end-start+1)

	var openedPorts []uint32

	for i := start; i <= end; i++ {
		ports <- i
	}

	toProcess := uint32(0)
	for i := uint32(0); i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for toProcess < end-start+1 {
				atomic.AddUint32(&toProcess, 1)

				port := <-ports
				if IsPortOpen(host, protocol, port) {
					mutex.Lock()
					openedPorts = append(openedPorts, port)
					mutex.Unlock()
				}
			}

		}()
	}

	wg.Wait()
	close(ports)

	return openedPorts
}

// Port2ServiceByFingerprint returns the service name of a port.
func Port2ServiceByFingerprint(host, protocol string, port uint32) string {
	if services.IsPrometheusExporter(host, port) {
		return "prometheus_exporter"
	}

	if services.IsSSH(host, port) {
		return "ssh"
	}

	return ""
}

// Port2Service returns the service name of a port.
func Port2Service(host, protocol string, port uint32, fingerprint bool) string {
	proto := netdb.ServiceByPort(int(port), protocol)

	if proto == nil {
		if fingerprint {
			return Port2ServiceByFingerprint(host, protocol, port)
		}
		return ""
	}
	return proto.Name
}
