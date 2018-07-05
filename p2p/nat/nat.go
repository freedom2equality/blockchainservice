package nat

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/log"
)

// NAT is an interface representing a NAT traversal options for example UPNP or
// NAT-PMP. It provides methods to query and manipulate this traversal to allow
// access to services.
type NAT interface {
	// Get the external address from outside the NAT.
	GetExternalAddress() (addr net.IP, err error)
	// Add a port mapping for protocol ("udp" or "tcp") from external port to
	// internal port with description lasting for timeout.
	AddPortMapping(protocol string, externalPort, internalPort int, description string, timeout int) (mappedExternalPort int, err error)
	// Remove a previously added port mapping from external port to
	// internal port.
	DeletePortMapping(protocol string, externalPort, internalPort int) (err error)

	// Should return name of the method. This is used for logging.
	String() string
}

// Parse parses a NAT interface description.
// The following formats are currently accepted.
// Note that mechanism names are not case-sensitive.
//
//     "" or "none"         return nil
//     "extip:77.12.33.4"   will assume the local machine is reachable on the given IP
//     "any"                uses the first auto-detected mechanism
//     "upnp"               uses the Universal Plug and Play protocol
//     "pmp"                uses NAT-PMP with an auto-detected gateway address
//     "pmp:192.168.0.1"    uses NAT-PMP with the given gateway address
func Parse(spec string) (NAT, error) {
	var (
		parts = strings.SplitN(spec, ":", 2)
		mech  = strings.ToLower(parts[0])
		ip    net.IP
	)
	if len(parts) > 1 {
		ip = net.ParseIP(parts[1])
		if ip == nil {
			return nil, errors.New("invalid IP address")
		}
	}
	switch mech {
	case "", "none", "off":
		return nil, nil
	case "any", "auto", "on":
		return Any(), nil
	case "extip", "ip":
		if ip == nil {
			return nil, errors.New("missing IP address")
		}
		return ExtIP(ip), nil
	case "upnp":
		return UPnP(), nil
	case "pmp", "natpmp", "nat-pmp":
		return PMP(ip), nil
	default:
		return nil, fmt.Errorf("unknown mechanism %q", parts[0])
	}
}

const (
	mapTimeout        = 20 * time.Minute
	mapUpdateInterval = 15 * time.Minute
)

// Map adds a port mapping on m and keeps it alive until c is closed.
// This function is typically invoked in its own goroutine.
func Map(m NAT, c chan struct{}, protocol string, extport, intport int, name string) {
	//log := log.New("proto", protocol, "extport", extport, "intport", intport, "interface", m)
	refresh := time.NewTimer(mapUpdateInterval)
	defer func() {
		refresh.Stop()
		log.Debug("Deleting port mapping")
		m.DeletePortMapping(protocol, extport, intport)
	}()
	if _, err := m.AddPortMapping(protocol, extport, intport, name, int(mapTimeout/time.Second)); err != nil {
		log.Debug("Couldn't add port mapping", "err", err)
	} else {
		log.Info("Mapped network port")
	}
	for {
		select {
		case _, ok := <-c:
			if !ok {
				return
			}
		case <-refresh.C:
			log.Trace("Refreshing port mapping")
			if _, err := m.AddPortMapping(protocol, extport, intport, name, int(mapTimeout/time.Second)); err != nil {
				log.Debug("Couldn't add port mapping", "err", err)
			}
			refresh.Reset(mapUpdateInterval)
		}
	}
}

func ExtIP(ip net.IP) NAT {
	if ip == nil {
		panic("IP must not be nil")
	}
	return extIP(ip)
}

type extIP net.IP

func (n extIP) GetExternalAddress() (addr net.IP, err error) { return net.IP(n), nil }
func (n extIP) String() string                               { return fmt.Sprintf("ExtIP(%v)", net.IP(n)) }

// These do nothing.
func (extIP) AddPortMapping(protocol string, externalPort, internalPort int, description string, timeout int) (mappedExternalPort int, err error) {
	return 0, nil
}
func (extIP) DeletePortMapping(protocol string, externalPort, internalPort int) (err error) {
	return nil
}

// Any returns a port mapper that tries to discover any supported
// mechanism on the local network.
func Any() NAT {
	// TODO: attempt to discover whether the local machine has an
	// Internet-class address. Return ExtIP in this case.
	return startautodisc("UPnP or NAT-PMP", func() (NAT, error) {
		type NATERR struct {
			nat NAT
			err error
		}
		found := make(chan NATERR, 2)
		go func() {
			nat, err := DiscoverUPnP()
			found <- NATERR{nat: nat, err: err}
		}()
		go func() {
			nat, err := DiscoverUPnP()
			found <- NATERR{nat: nat, err: err}
		}()
		for i := 0; i < cap(found); i++ {
			if nat := <-found; nat.nat != nil {
				return nat.nat, nat.err
			}
		}
		return nil, nil
	})
}

// UPnP returns a port mapper that uses UPnP. It will attempt to
// discover the address of your router using UDP broadcasts.
func UPnP() NAT {
	return startautodisc("UPnP", DiscoverUPnP)
}

// PMP returns a port mapper that uses NAT-PMP. The provided gateway
// address should be the IP of your router. If the given gateway
// address is nil, PMP will attempt to auto-discover the router.
func PMP(gateway net.IP) NAT {
	//if gateway != nil {
	//	return &pmp{gw: gateway, c: natpmp.NewClient(gateway)}
	//}
	return startautodisc("NAT-PMP", DiscoverPMP)
}

type autodisc struct {
	what string // type of interface being autodiscovered
	once sync.Once
	doit func() (NAT, error)

	mu    sync.Mutex
	found NAT
}

func startautodisc(what string, doit func() (NAT, error)) NAT {
	// TODO: monitor network configuration and rerun doit when it changes.
	return &autodisc{what: what, doit: doit}
}

func (n *autodisc) AddPortMapping(protocol string, externalPort, internalPort int, description string, timeout int) (mappedExternalPort int, err error) {
	if err := n.wait(); err != nil {
		return 0, err
	}
	return n.found.AddPortMapping(protocol, externalPort, internalPort, description, timeout)
}

func (n *autodisc) DeletePortMapping(protocol string, externalPort, internalPort int) (err error) {
	if err := n.wait(); err != nil {
		return err
	}
	return n.found.DeletePortMapping(protocol, externalPort, internalPort)
}

func (n *autodisc) GetExternalAddress() (net.IP, error) {
	if err := n.wait(); err != nil {
		return nil, err
	}
	return n.found.GetExternalAddress()
}

func (n *autodisc) String() string {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.found == nil {
		return n.what
	} else {
		return n.found.String()
	}
}

// wait blocks until auto-discovery has been performed.
func (n *autodisc) wait() error {
	n.once.Do(func() {
		n.mu.Lock()
		nat, err := n.doit()
		n.found = nil
		if err == nil {
			n.found = nat
		}

		n.mu.Unlock()
	})
	if n.found == nil {
		return fmt.Errorf("no %s router discovered", n.what)
	}
	return nil
}
