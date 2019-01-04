package rest

import (
	"fmt"
	"net"
)

var (
	noFreePortErr = fmt.Errorf("no free port found from ports")
)

// gets a free port from given ports slice.
// if no free port could be found error != nil and port == -1.
func GetFreePort(address string, ports ...int) (int, error) {
	for _, port := range ports {
		err := checkPort(address, port)
		if err != nil {
			continue
		}
		return port, nil
	}
	return -1, noFreePortErr
}

// gets a free port from given port range including.
// if no free port could be found error != nil and port == -1.
func GetFreePortInRange(address string, minPort int, maxPort int) (int, error) {
	for port := minPort; port <= maxPort; port++ {
		err := checkPort(address, port)
		if err != nil {
			continue
		}
		return port, nil
	}
	return -1, noFreePortErr
}

// checks for a single port being available in the
// network. if not error != nil
func checkPort(address string, port int) error {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return err
	}

	n, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	defer func() {
		_ = n.Close()
	}()
	return err
}
