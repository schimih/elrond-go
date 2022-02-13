package utils

import (
	"errors"
)

// ErrHandshakeFail signals that a handshake could not be established with target
var ErrHandshakeFail = errors.New("handshake failed")

// ErrNoEvaluationType signals that no evaluation type has been provided to the evaluator
var ErrNoEvaluationType = errors.New("No Evaluation Type given - Please choose evaluation type")

// ErrNoEvaluationType signals that no evaluation type has been provided to the evaluator
var ErrNoFormatterType = errors.New("No Formatter Type given - Please choose formatter type")
