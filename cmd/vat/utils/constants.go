package utils

// AnalysisType represents the type of Analysis to be ran
type AnalysisType int

// Enumerates some of the analysis that can be ran
// TCP_ELROND - Scan only ELROND Ports
//
// TCP_WEB - Scan only WEB Ports
//
// TCP_SSH - Scan only SSH Ports
//
// TCP_ALL - Scan all ports
//
// TCP_STANDARD - Scan top 1000 ports
//
// TCP_REQ1 - Scan ports as per requirement
//
// TCP_SSH_ALGOS - Scan installed algos on SSH ports
//
// TCP_POLITE_REQ1 - Creates a scanner that uses ssh.Dial function to check access to a SSH port
//
// TCP_BRUTE_REQ1 - Uses nmap scanner to check access to a SSH port
const (
	TCP_ELROND      AnalysisType = iota
	TCP_WEB         AnalysisType = 1
	TCP_SSH         AnalysisType = 2
	TCP_ALL         AnalysisType = 3
	TCP_STANDARD    AnalysisType = 4
	TCP_REQ1        AnalysisType = 5
	TCP_SSH_ALGOS   AnalysisType = 6
	TCP_POLITE_REQ1 AnalysisType = 7
	TCP_BRUTE_REQ1  AnalysisType = 8
)

// EvaluationType represents the type of Evaluation that can be ran.
//
// For now ONLY "portStatusEvaluation" evaluation is implemented. Other implementations
// like (fingerprint, dns, bruteforce etc) will follow.
type EvaluationType int

// EvaluationType represents the type of Evaluation that can be ran.
//
// For now ONLY "portStatusEvaluation" evaluation is implemented. Other implementations
// like (fingerprint, dns, bruteforce etc) will follow.
//
// Req1_Evaluation - PortStatusEvaluation + POLITE - if SSH port is open -> try introducing Username & Pass
//
// The evaluationType will be controlled by "manager".
const (
	NoEvaluation                EvaluationType = iota
	PortStatusEvaluation        EvaluationType = 1
	Polite_PortAndSshEvaluation EvaluationType = 2
	Brute_PortAndSshEvaluation  EvaluationType = 3
	// DNS etc
)

// PortStatus represents a port's status.
type PortStatus string

// Enumerates the different possible state values of a port
const (
	Open       PortStatus = "open"
	Reset      PortStatus = "reset"
	Closed     PortStatus = "closed"
	Filtered   PortStatus = "filtered"
	Unfiltered PortStatus = "unfiltered"
)

// PortType represents a port's type.
type PortType int

// Enumerates the different possible types of a port (elrond 37373-38383, web- 80,8080, ssh - 22, etc)
const (
	ElrondPort    PortType = iota
	WebPort       PortType = 1
	SshPort       PortType = 2
	OutsideElrond PortType = 3
	Unknown       PortType = 4
)

// TargetStatus represents a target's status
type TargetStatus string

// Enumerates the different possible state values of a target
const (
	NEW       TargetStatus = "NEW"
	SCANNED   TargetStatus = "SCANNED"
	EXPIRED   TargetStatus = "EXPIRED"
	EVALUATED TargetStatus = "EVALUATED"
)

// TargetStatus represents a scanner's status
type ScannerStatus int

// Enumerates the different possible state values of a scanner
const (
	NULL        ScannerStatus = iota
	NOT_STARTED ScannerStatus = 1
	IN_PROGRESS ScannerStatus = 2
	FAILED      ScannerStatus = 3
	DONE        ScannerStatus = 4
	FINISHED    ScannerStatus = 5
)

type Risk int

const (
	HighRisk   Risk = 25
	MediumRisk Risk = 10
	SmallRisk  Risk = 5
	NoRisk     Risk = 0
)

type Judgement string

const (
	JudgementSshUserPermited Judgement = "5 points deducted - SMALL RISK: SSH User permitted"
	JudgementSshPwdPermited  Judgement = "25 points deducted - HIGH RISK - SSH PasswordAuthentication active"
	JudgementDummyPermited   Judgement = "50 points deducted - ALERT - Dummy values for test have been accepted"
	JudgementMediumRisk      Judgement = "10 points deducted - MEDIUM RISK - Port Outside Elrond Permited ports Open"
	JudgementWeb             Judgement = "5 points deducted - SMALL RISK - Tcp-Web Port Open"
	JudgementSsh             Judgement = "5 points deducted - SMALL RISK - SSH Port Open"
	JudgementNoRisk          Judgement = "0 points deducted - NO RISK - Elrond Port Open"
	JudgementFromPort        Judgement = "0 points deducted - "
)

type SecureLevel int

const (
	HIGH  SecureLevel = 0
	MID   SecureLevel = 1
	LOW   SecureLevel = 2
	ALERT SecureLevel = 3
)

type OutputType int

const (
	Table    OutputType = 0
	JSON     OutputType = 1
	XML      OutputType = 2
	GIN      OutputType = 3
	JustScan OutputType = 4
)

const JsonFilePath = "./cmd/vat/AnalysisResult.json"
const XMLFilePath = "./cmd/vat/AnalysisResult.xml"
