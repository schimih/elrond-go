package result

func (port *Port) PortRange(minPort int, maxPort int) bool {
	if (port.Number < minPort) || (port.Number > maxPort) {
		return true
	}
	return false
}

func (rating *Rating) Process(assessment string) {
	switch assessment {
	case "TCP-ELROND":
		rating.State = "assessed"
		rating.Reason = append(rating.Reason, "Only Elrond Ports Open")
		rating.Value = 100
	case "TCP-WEB":
		rating.State = "assessed"
		rating.Reason = append(rating.Reason, "Only Elrond Ports Open")
		rating.Value = rating.Value - 10
	case "TCP-SSH":
		rating.State = "assessed"
		rating.Reason = append(rating.Reason, "Only Elrond Ports Open")
		rating.Value = rating.Value - 10
	case "TCP-ALL":
		log.Info("Scans everything")

	default:
		return
	}
	return

}
