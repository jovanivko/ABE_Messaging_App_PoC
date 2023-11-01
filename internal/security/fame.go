package security

import (
	"abeProofOfConcept/internal/logger"
	"abeProofOfConcept/pkg/store"
	"abeProofOfConcept/pkg/utils"
	"context"
	"github.com/fentec-project/gofe/abe"
)

type FAME = abe.FAME
type FAMEPubKey = abe.FAMEPubKey
type FAMESecretKey = abe.FAMESecKey

type MSP = abe.MSP

type FameCipher = abe.FAMECipher

var (
	fame          *FAME
	FamePublicKey *FAMEPubKey
	fmsk          *FAMESecretKey
)

func init() {
	fame = abe.NewFAME()
}

func InitFame() error {
	fameKeys, err := store.GetFameKeys(context.Background())
	if err != nil {
		FamePublicKey, fmsk, err = fame.GenerateMasterKeys()
		if err != nil {
			logger.Logger.Errorln("Couldn't generate master keys for FAME!", err)
			return err
		}
		pk, err := utils.EncodeToBytes(FamePublicKey)
		if err != nil {
			logger.Logger.Errorln("Couldn't encode public key for FAME!", err)
			return err
		}
		msk, err := utils.EncodeToBytes(fmsk)
		if err != nil {
			logger.Logger.Errorln("Couldn't encode secret key for FAME!", err)
			return err
		}
		_, err = store.CreateFameKeys(context.Background(), pk, msk)
		if err != nil {
			logger.Logger.Errorln("Couldn't store keys for FAME!", err)
			return err
		}
		return nil
	}
	FamePublicKey = &FAMEPubKey{}
	fmsk = &FAMESecretKey{}

	if err = utils.DecodeToObject(fameKeys.PublicKey, FamePublicKey); err != nil {
		logger.Logger.Errorln("Couldn't decode FAME public key!", err)
		return err
	}
	if err = utils.DecodeToObject(fameKeys.MasterSecretKey, fmsk); err != nil {
		logger.Logger.Errorln("Couldn't decode FAME secret key!", err)
		return err
	}
	return nil
}

func DecryptFameCipher(fameCipher *FameCipher, email string, position string, department string) (string, error) {
	attributes := []string{"email:" + email, "role:" + position, "department:" + department}
	keys, err := fame.GenerateAttribKeys(attributes, fmsk)
	if err != nil {
		logger.Logger.Errorln("Couldn't generate attribute keys for FAME!", err)
		return "", err
	}
	plaintext, err := fame.Decrypt(fameCipher, keys, FamePublicKey)
	if err != nil {
		logger.Logger.Errorln("Couldn't decrypt FAME cipher!", err)
		return "", err
	}
	return plaintext, err
}
