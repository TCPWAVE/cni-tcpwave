package main

import(
  "fmt"
  "net"
  "sync"
  "log"
  "github.com/containernetworking/plugins/pkg/ns"
)

// InterfaceInfo Net Interface Info
type InterfaceInfo struct {
	iface *net.Interface
	wg    sync.WaitGroup
}

// GetMacAddress extract MAC Address
func GetMacAddress(netns string, ifaceName string, logger *log.Logger) (string) {
  logger.Println("Computing MAC address")
	var err error
  var mac string
	ifaceInfo := &InterfaceInfo{}
	if netns == "" {
		ifaceInfo.iface, err = net.InterfaceByName(ifaceName)
		if err != nil {
			return ""
		}
		mac = ifaceInfo.iface.HardwareAddr.String()
	} else {
		errCh := make(chan error, 1)
		ifaceInfo.wg.Add(1)
		go func() {
			errCh <- ns.WithNetNSPath(netns, func(_ ns.NetNS) error {
				defer ifaceInfo.wg.Done()

				ifaceInfo.iface, err = net.InterfaceByName(ifaceName)
				if err != nil {
					return fmt.Errorf("error looking up interface '%s': '%s'", ifaceName, err)
				}
				return nil
			})
		}()

		if err = <-errCh; err != nil {
			fmt.Printf("%s\n", err)
		} else {
			mac = ifaceInfo.iface.HardwareAddr.String()
		}
	}
  logger.Printf("Computed MAC = %s",mac)
	return mac
}
