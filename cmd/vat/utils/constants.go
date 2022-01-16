package utils

// AnalysisType represents the type of Analysis to be ran
type AnalysisType int

// Enumerates some of the commands that can be ran
const (
	TCP_ELROND         AnalysisType = iota
	TCP_OUTSIDE_ELROND AnalysisType = 1
	TCP_WEB            AnalysisType = 2
	TCP_SSH            AnalysisType = 3
	TCP_FULL           AnalysisType = 4
	TCP_STANDARD       AnalysisType = 5
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

// TargetStatus represents a target's status
type TargetStatus string

// Enumerates the different possible state values of a target
const (
	NEW     TargetStatus = "NEW"
	SCANNED TargetStatus = "SCANNED"
	EXPIRED TargetStatus = "EXPIRED"
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

// NmapCommand represents a string nmap command
type NmapCommand string

/*
-Pn --skip the ping test and simply scan every target host provided.
-sS --stealth scan,fastest way to scan ports of the most popular protocol (TCP).
-pn --port to be scanned.
-sC --
*/
const (
	NMAP_TCP_ELROND         NmapCommand = "-Pn -sS -p37373-38383"
	NMAP_TCP_OUTSIDE_ELROND NmapCommand = "-Pn -sS -p-37372,38384-"
	NMAP_TCP_WEB            NmapCommand = "-Pn -p80,8080,280,443" // added: http-mgmt (280), https (443)
	NMAP_TCP_SSH            NmapCommand = "-Pn -p22"
	NMAP_TCP_FULL           NmapCommand = "-Pn -sS -A -p-"
	NMAP_TCP_STANDARD       NmapCommand = "--randomize-hosts -Pn -sS -A -T4 -g53 --top-ports 1000"
)
