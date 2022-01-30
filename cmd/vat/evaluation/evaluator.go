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
	Address        string
	Ports          []scan.Port
	Status         string
	Score          int
	SecurityLevel  utils.SecureLevel
	Judgements     []string
	EvaluationType utils.EvaluationType
	scannerFactory analysis.ScannerFactory
}

func (eT *EvaluatedTarget) Evaluate() EvaluatedTarget {
	eT.Score -= eT.calculateDeductionPoints()
	eT.SecurityLevel = eT.calculateSecurityLevel()
	eT.Status = string(utils.EVALUATED)
	return *eT
}

func (eT *EvaluatedTarget) calculateDeductionPoints() (deductionPoints int) {
	judgement := utils.JudgementFromPort
	deductionPoints = 0

	for _, port := range eT.Ports {
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
	return ((portNumber == 22) && (eT.EvaluationType == utils.Polite_PortAndSshEvaluation)) ||
		((portNumber == 22) && (eT.EvaluationType == utils.Brute_PortAndSshEvaluation))
}

func (eT *EvaluatedTarget) runSpecialCheck() (judgement utils.Judgement) {
	if eT.EvaluationType == utils.Polite_PortAndSshEvaluation {
		return eT.runPoliteSshCheck()
	} else {
		return eT.runBruteForceSshCheck()
	}
}

func (eT *EvaluatedTarget) runPoliteSshCheck() (judgement utils.Judgement) {
	s := eT.scannerFactory.CreateScanner(eT.Address, utils.TCP_POLITE_REQ1)

	_, err := s.Scan()
	if err != nil {
		return interpretErrorGravity(err)
	}

	return utils.JudgementDummyPermited
}

func interpretErrorGravity(err error) (judgement utils.Judgement) {
	errSlice := strings.Split(err.Error(), " ")

	if errSlice[9] == "[none]," {
		return utils.JudgementSshUserPermited
	} else if (errSlice[9] == "[none") && (errSlice[10] == "password],") {
		return utils.JudgementSshPwdPermited
	} else {
		return utils.JudgementFromPort
	}
}

// BruteForce was never called => the result interpretation has not been done.
// Has to be checked with the team if permitted.
func (eT *EvaluatedTarget) runBruteForceSshCheck() (judgement utils.Judgement) {
	s := eT.scannerFactory.CreateScanner(eT.Address, utils.TCP_BRUTE_REQ1)

	_, err := s.Scan()
	if err != nil {
		log.Error("Scan failed because %e", err)
		return
	}

	return
}

func (eT *EvaluatedTarget) runJudgement(portJudgement utils.Judgement, checkJudgement utils.Judgement) (synopsis int) {
	if checkJudgement != utils.JudgementFromPort {
		eT.Judgements = append(eT.Judgements, fmt.Sprintf(string(checkJudgement)))
	}

	eT.Judgements = append(eT.Judgements, fmt.Sprintf(string(portJudgement)))

	return extractNumber(string(portJudgement)) + extractNumber(string(checkJudgement))
}

func extractNumber(judgement string) int {
	slice := strings.Split(judgement, "$ ")
	i, _ := strconv.Atoi(slice[0])

	return i
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
