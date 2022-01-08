package result

import go_nmap "github.com/lair-framework/go-nmap"

type Peer struct {
	ID           uint
	Protocol     string
	Address      string
	Standard     string
	Ports        []Port
	Architecture string
	Status       string
	TestType     string
	Evaluation   Rating
}

type Rating struct {
	State  string
	Reason []string
	Value  uint
}

type RatingState string

const (
	Assessed RatingState = "assessed"
	Sweeped  RatingState = "sweeped"
	Unknown  RatingState = "unknown"
)

func (rating Rating) Status() RatingState {
	return RatingState(rating.State)
}

type Port struct {
	ID       uint
	Number   int
	Protocol string
	State    string
	Owner    string
}

// PortStatus represents a port's status.
type PortStatus string

// Enumerates the different possible state values.
const (
	Open       PortStatus = "open"
	Closed     PortStatus = "closed"
	Filtered   PortStatus = "filtered"
	Unfiltered PortStatus = "unfiltered"
)

// Status returns the status of a port.
func (p Port) Status() PortStatus {
	return PortStatus(p.State)
}

func (r *ResultsContainer) Process(NmapScanResult []*go_nmap.NmapRun, testType string) *ResultsContainer {
	var tempPort Port
	var tempPeer Peer
	portSlice := make([]Port, 0)
	for _, nmapRun := range NmapScanResult {
		for hidx, host := range nmapRun.Hosts {
			for pidx, port := range host.Ports {
				portSlice = append(portSlice, *tempPort.Add(uint(pidx), port.PortId, port.Protocol, port.State.State, port.Owner.Name))
			}
			r.Results = append(r.Results, *tempPeer.Add(uint(hidx), host.Addresses[0].Addr, portSlice, host.Status.State, "SCANNED", testType))
		}
	}
	return r
}

func (peer *Peer) Add(id uint, address string, ports []Port, status string, state string, test string) *Peer {

	var emptyRating Rating
	emptyRating.State = state
	emptyRating.Value = 100
	peer = &Peer{
		ID:           id,
		Protocol:     "",
		Address:      address,
		Standard:     "",
		Ports:        ports,
		Architecture: "",
		Status:       status,
		TestType:     test,
		Evaluation:   emptyRating,
	}
	return peer
}

func (port *Port) Add(id uint,
	number int,
	protocol string,
	state string,
	owner string) *Port {

	port = &Port{
		ID:       id,
		Number:   number,
		Protocol: protocol,
		State:    state,
		Owner:    owner,
	}

	return port
}
