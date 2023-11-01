package client

import (
	"abeProofOfConcept/internal/logger"
	"abeProofOfConcept/pkg/utils"
	"github.com/fentec-project/gofe/abe"
)

const MaxGpswAttributes = 10

var (
	FamePublicKey *abe.FAMEPubKey
	fame          *abe.FAME
)

func init() {
	fame = abe.NewFAME()
}

func EncryptFame(msg string, policy string) ([]byte, error) {
	msp, err := abe.BooleanToMSP(policy, false)
	if err != nil {
		logger.Logger.Errorln(
			"message", "Error generating MSP for FAME",
			"error", err,
		)
		return nil, err
	}
	cipher, err := fame.Encrypt(msg, msp, FamePublicKey)
	if err != nil {
		logger.Logger.Errorln(
			"message", "Error encrypting message with FAME",
			"error", err,
		)
		return nil, err
	}
	bytes, err := utils.EncodeToBytes(cipher)
	if err != nil {
		logger.Logger.Errorln(
			"message", "Error encoding the FAME cipher",
			"error", err,
		)
		return nil, err
	}
	return bytes, nil
}
