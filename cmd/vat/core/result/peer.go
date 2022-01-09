package result

type Peer struct {
	ID           uint
	Protocol     string
	Address      string
	Standard     string
	Ports        []Port
	Architecture string
	Status       string
	AnalysisType string
	Evaluation   Rating
}

func NewPeer(id uint, address string, ports []Port, peerStatus string, state string, analysisType string) Peer {
	var emptyRating Rating
	emptyRating.State = state
	emptyRating.Value = 100

	return Peer{
		ID:           id,
		Protocol:     "",
		Address:      address,
		Standard:     "",
		Ports:        ports,
		Architecture: "",
		Status:       peerStatus,
		AnalysisType: analysisType,
		Evaluation:   emptyRating,
	}
}
