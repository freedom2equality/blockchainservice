package nat

import (
	"fmt"
	"net"
)

func DiscoverPMP() (nat NAT, err error) {
	return nil, nil
}

type pmp struct {
	gw net.IP
}

func (n *pmp) String() string {
	return fmt.Sprintf("NAT-PMP(%v)", n.gw)
}

func (n *pmp) GetExternalAddress() (addr net.IP, err error) {
	// todo
	return nil, nil
}

func (n *pmp) AddPortMapping(protocol string, externalPort, internalPort int, description string, timeout int) (mappedExternalPort int, err error) {
	// todo
	return 0, nil
}

func (n *pmp) DeletePortMapping(protocol string, externalPort, internalPort int) (err error) {
	// todo
	return nil
}
