package evaluation

import "github.com/elrond-go/cmd/vat/core/scan"

type NodeStatus string

const (
	Evaluated   NodeStatus = "evaluated"
	Sweeped     NodeStatus = "sweeped"
	Unknown     NodeStatus = "unknown"
	NotAssessed NodeStatus = "not-assessed"
)

func (result EvaluationResult) Status() NodeStatus {
	return NodeStatus(result.EvaluationType)
}

type Node struct {
	Address string
	Ports   []scan.Port
	Status  string
}

func NewNode(address string, portSlice []scan.Port) Node {
	return Node{
		Address: address,
		Ports:   portSlice,
		Status:  string(NotAssessed),
	}
}
