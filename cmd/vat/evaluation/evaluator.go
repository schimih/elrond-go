package evaluation

import (
	"fmt"

	"github.com/elrond-go/cmd/vat/scan"
	"github.com/elrond-go/cmd/vat/utils"
)

type EvaluatedTarget struct {
	Address       string
	Ports         []scan.Port
	Status        string
	Score         int
	SecurityLevel utils.SecureLevel
	Judgements    []string
}

func (eT *EvaluatedTarget) Evaluate() EvaluatedTarget {
	deduction := 0
	for _, port := range eT.Ports {
		if port.State == utils.Open {
			deduction += int(port.RiskValue)
			eT.addJudgement(port)
		}
	}

	// deduct based on the identified risk
	eT.Score -= deduction
	eT.SecurityLevel = eT.calculateSecurityLevel()
	eT.Status = string(utils.EVALUATED)
	return *eT
}

func (eT *EvaluatedTarget) addJudgement(port scan.Port) {
	eT.Judgements = append(eT.Judgements, fmt.Sprintf("%s - %d points deducted", port.RiskReason, port.RiskValue))
}

func (eT *EvaluatedTarget) calculateSecurityLevel() utils.SecureLevel {
	if eT.Score >= 80 {
		return utils.HIGH
	} else if (eT.Score >= 60) && (eT.Score < 80) {
		return utils.MID
	} else if (eT.Score >= 40) && (eT.Score < 60) {
		return utils.LOW
	}
	return utils.ALERT
}

// IsInterfaceNil returns true if there is no value under the interface
func (eT *EvaluatedTarget) IsInterfaceNil() bool {
	return eT == nil
}
