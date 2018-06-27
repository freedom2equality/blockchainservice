package p2p

// Manage handles peer connections and exposes an API to receive incoming messages on `Business`
type Manage struct {
	reactors map[string]Reactor
}
