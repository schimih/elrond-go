package evaluation

import (
	"fmt"
)

// SIMPLE Evaluation -> evaluate based on port's status
// COMPLEX Evaluation -> after interpreting a port's status -> do other actions-> ex. fingerprint, bruteforce, dns, domain, etc
const (
	SIMPLE = iota
	COMPLEX
)

type EvaluationResult struct {
	Node           Node
	EvaluationType int
	Score          score
	SecurityLevel  securityLevel
	Reasons        []string
}

type score int
type securityLevel int

const (
	highRisk   score = -25
	mediumRisk score = -10
	smallRisk  score = -5
	noRisk     score = 0
)

const (
	HIGH  securityLevel = 0
	MID   securityLevel = 1
	LOW   securityLevel = 2
	ALERT securityLevel = 3
)

const (
	TCP_ELROND = iota
	TCP_OUTSIDE_ELROND
	TCP_WEB
	TCP_SSH
	TCP_FULL
	TCP_STANDARD
	UNKOWN
)

// Enumerates the different possible state values.
const (
	Open       = "open"
	Closed     = "closed"
	Filtered   = "filtered"
	Unfiltered = "unfiltered"
	Reset      = "reset"
)

func (e *EvaluationResult) Evaluate(node Node) EvaluationResult {
	for _, port := range node.Ports {
		if port.State == Open {
			e.deduct(port.Number)
		}
	}

	return *e
}

func (e *EvaluationResult) deduct(portNumber int) {
	switch getKindByPort(portNumber) {
	case TCP_ELROND:
		e.processEvaluation(TCP_ELROND, portNumber)
	case TCP_WEB:
		e.processEvaluation(TCP_WEB, portNumber)
	case TCP_SSH:
		e.processEvaluation(TCP_SSH, portNumber)
	case TCP_STANDARD:
		e.processEvaluation(TCP_STANDARD, portNumber)
	default:
		return
	}
}

func getKindByPort(portNumber int) int {
	if isPortInRange(portNumber, 37373, 38383) {
		return TCP_ELROND
	} else if (portNumber == 80) || (portNumber == 8080) {
		return TCP_WEB
	} else if portNumber == 22 {
		return TCP_SSH
	} else if isPortInRange(portNumber, 0, 37372) {
		return TCP_STANDARD
	} else if isPortInRange(portNumber, 38384, 65535) {
		return TCP_STANDARD
	} else {
		return 0
	}
}

func isPortInRange(port int, low int, high int) bool {
	if (port > low) && (port < high) {
		return true
	}
	return false
}

func (e *EvaluationResult) processEvaluation(portType int, portNumber int) {
	e.Reasons = append(e.Reasons, fmt.Sprintf("NO RISK - Elrond Port Open - 0 points deducted", "port", portNumber))
	e.Score += evaluateRisk(portType)
	e.SecurityLevel = calculateSecurityLevel(e.Score)
}

func formatReason()

func evaluateRisk(portType int) score {
	switch portType {
	case TCP_ELROND:
		return noRisk
	case TCP_WEB:
		return mediumRisk
	case TCP_SSH:
		return mediumRisk
	case TCP_STANDARD:
		return smallRisk
	default:
		return noRisk
	}
}

func calculateSecurityLevel(score score) securityLevel {
	if score >= 80 {
		return HIGH
	} else if (score >= 60) && (score < 80) {
		return MID
	} else if (score >= 40) && (score < 60) {
		return LOW
	}
	return ALERT
}

// IsInterfaceNil returns true if there is no value under the interface
func (e *EvaluationResult) IsInterfaceNil() bool {
	return e == nil
}
