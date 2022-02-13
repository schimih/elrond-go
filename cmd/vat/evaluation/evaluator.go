package evaluation

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/elrond-go/cmd/vat/analysis"
	"github.com/elrond-go/cmd/vat/scan"
	"github.com/elrond-go/cmd/vat/utils"
)

type EvaluatedTarget struct {
	identity       identity
	evaluation     evaluationResult
	scannerFactory analysis.ScannerFactory
}

func (eT *EvaluatedTarget) Evaluate() EvaluatedTarget {
	eT.evaluation.score -= eT.calculateDeductionPoints()
	eT.evaluation.securityLevel = eT.calculateSecurityLevel()
	eT.identity.status = string(utils.EVALUATED)

	return *eT
}

func (eT *EvaluatedTarget) calculateDeductionPoints() (deductionPoints int) {
	judgement := utils.JudgementFromPort
	deductionPoints = 0

	for _, port := range eT.identity.ports {
		if port.State == utils.Open {
			if eT.requestedSpecialCheck(port.Number) {
				judgement = eT.runSpecialCheck()
			} else {
				judgement = utils.JudgementFromPort
			}
			deductionPoints += eT.runJudgement(port.Importance, judgement)
		}
	}

	return
}

func (eT *EvaluatedTarget) requestedSpecialCheck(portNumber int) bool {
	return ((portNumber == 22) && (eT.evaluation.evaluationType == utils.Polite_PortAndSshEvaluation)) ||
		((portNumber == 22) && (eT.evaluation.evaluationType == utils.Brute_PortAndSshEvaluation))
}

func (eT *EvaluatedTarget) runSpecialCheck() (judgement utils.Judgement) {
	if eT.evaluation.evaluationType == utils.Polite_PortAndSshEvaluation {
		return eT.runPoliteSshCheck()
	} else {
		return eT.runBruteForceSshCheck()
	}
}

func (eT *EvaluatedTarget) runPoliteSshCheck() (judgement utils.Judgement) {
	s := eT.scannerFactory.CreateScanner(eT.identity.address, utils.TCP_POLITE_REQ1)

	_, err := s.Scan()
	if err != nil {
		return interpretErrorGravity(err)
	}

	return utils.JudgementDummyPermited
}

func interpretErrorGravity(err error) (judgement utils.Judgement) {
	errSlice := strings.Split(err.Error(), " ")

	for i := range errSlice {
		if errSlice[i] == "[none]," {
			return utils.JudgementSshUserPermited
		}
		if errSlice[i] == "[none" {
			// next word in string
			if errSlice[i+1] == "password]," {

				return utils.JudgementSshPwdPermited
			}
		}
	}

	return utils.JudgementFromPort
}

// BruteForce was never called => the result interpretation has not been done.
// Has to be checked with the team if permitted.
func (eT *EvaluatedTarget) runBruteForceSshCheck() (judgement utils.Judgement) {
	s := eT.scannerFactory.CreateScanner(eT.identity.address, utils.TCP_BRUTE_REQ1)

	_, err := s.Scan()
	if err != nil {
		log.Error("Scan failed because %e", err)
		return
	}

	return
}

func (eT *EvaluatedTarget) runJudgement(portJudgement utils.Judgement, checkJudgement utils.Judgement) (synopsis int) {
	if checkJudgement != utils.JudgementFromPort {
		eT.evaluation.judgements = append(eT.evaluation.judgements, fmt.Sprintf(string(checkJudgement)))
	}

	eT.evaluation.judgements = append(eT.evaluation.judgements, fmt.Sprintf(string(portJudgement)))

	return extractNumber(string(portJudgement)) + extractNumber(string(checkJudgement))
}

func extractNumber(judgement string) int {
	slice := strings.Split(judgement, " points")
	i, _ := strconv.Atoi(slice[0])

	return i
}

func (eT *EvaluatedTarget) calculateSecurityLevel() utils.SecureLevel {
	if eT.evaluation.score >= 80 {
		return utils.HIGH
	}
	if (eT.evaluation.score >= 60) && (eT.evaluation.score < 80) {
		return utils.MID
	}
	if (eT.evaluation.score >= 40) && (eT.evaluation.score < 60) {
		return utils.LOW
	}

	return utils.ALERT
}

func (eT *EvaluatedTarget) GetAddress() (address string) {
	return eT.identity.address
}

func (eT *EvaluatedTarget) GetPortsSlice() (ports []scan.Port) {
	return eT.identity.ports
}

func (eT *EvaluatedTarget) GetScore() (score int) {
	return eT.evaluation.score
}

func (eT *EvaluatedTarget) GetJudgements() (judgements []string) {
	return eT.evaluation.judgements
}

func (eT *EvaluatedTarget) GetSecurityLevel() (securityLevel utils.SecureLevel) {
	return eT.evaluation.securityLevel
}

func (eT *EvaluatedTarget) GetStatus() (status string) {
	return eT.identity.status
}

func (eT *EvaluatedTarget) GetEvaluationType() (evaluationType utils.EvaluationType) {
	return eT.evaluation.evaluationType
}

// IsInterfaceNil returns true if there is no value under the interface
func (eT *EvaluatedTarget) IsInterfaceNil() bool {
	return eT == nil
}
