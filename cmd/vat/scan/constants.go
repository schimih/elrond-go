package scan

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
	NMAP_TCP_WEB            NmapCommand = "-Pn -sS -p80,8080,280,443" // added: http-mgmt (280), https (443)
	NMAP_TCP_SSH            NmapCommand = "-Pn -p22"
	NMAP_TCP_FULL           NmapCommand = "-Pn -sS -A -p-"
	NMAP_TCP_STANDARD       NmapCommand = "--randomize-hosts -Pn -sS -A -T4 -g53 --top-ports 1000"
	NMAP_TCP_REQ1           NmapCommand = "-Pn -sS -p22,80,8080,280,443,37373-38383"
	NMAP_TCP_SSH_RUN        NmapCommand = "-p 22 --script=ssh-run /	--script-args="
	NMAP_TCP_SSH_ALGOS      NmapCommand = "--script ssh2-enum-algos"
	NMAP_TCP_SSH_BRUTE      NmapCommand = "-p 22 --script ssh-brute --script-args userdb=users.lst,passdb=pass.lst /	--script-args ssh-brute.timeout=4s"
)

const SSH_ARGS = "ssh-run.cmd=ls -l /, ssh-run.username=myusername, ssh-run.password=mypassword"
