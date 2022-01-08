package result

type assessor struct {
	assessments []AssessmentRunner
}

func (r *ResultsContainer) Evaluate(container *ResultsContainer) {
	for _, peer := range r.Results {
		peer.TestType = "TCP-ELROND"
		if peer.TestType == "TCP-ELROND" {
			if peer.Evaluation.State != "EVALUATED" {
				peer.RunEvaluation()
			} else {
				log.Info("PEER already evaluated")
			}
		} else {
			log.Info("Please run TCP-ALL testType for full Evaluation")
		}
	}
}

func (p *Peer) RunEvaluation() {
	for _, port := range p.Ports {
		var rating Rating
		if port.Status() == Open {
			if port.PortRange(37373, 38383) {
				p.Evaluation.Process("TCP-ELROND")
			} else {
				if port.Number == 80 {
					p.Evaluation.Process("TCP-WEB")
				}
				if port.Number == 8080 {
					p.Evaluation.Process("TCP-WEB")
				}
				if port.Number == 22 {
					//
					p.Evaluation.Process("TCP-SSH")
				}
			}
		} else {
			rating.State = "EVALUATED"
			rating.Reason = append(rating.Reason, "no open port")
			rating.Value = 999
		}
	}
}
