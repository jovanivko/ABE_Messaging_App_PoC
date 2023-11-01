package security

import (
	"abeProofOfConcept/internal/logger"
	"abeProofOfConcept/pkg/store"
	"abeProofOfConcept/pkg/utils"
	"context"
	"fmt"
	"github.com/fentec-project/gofe/abe"
	"github.com/fentec-project/gofe/data"
	"strconv"
)

type GPSW = abe.GPSW
type GPSWPubKey = abe.GPSWPubKey
type GPSWCipher = abe.GPSWCipher

const MaxNumberOfAttributes = 50

var (
	gpsw          *GPSW
	gpswPublicKey *GPSWPubKey
	gpswSecretKey data.Vector
)

func init() {
	gpsw = abe.NewGPSW(MaxNumberOfAttributes)
}

func InitGpsw() error {
	gpswKeys, err := store.GetGpswKeys(context.Background())
	if err != nil {
		gpswPublicKey, gpswSecretKey, err = gpsw.GenerateMasterKeys()
		if err != nil {
			logger.Logger.Errorln("Couldn't generate master keys for GPSW!", err)
			return err
		}
		pk, err := utils.EncodeToBytes(gpswPublicKey)
		if err != nil {
			logger.Logger.Errorln("Couldn't encode public key for GPSW!", err)
			return err
		}
		msk, err := utils.EncodeToBytes(gpswSecretKey)
		if err != nil {
			logger.Logger.Errorln("Couldn't encode secret key for GPSW!", err)
			return err
		}
		_, err = store.CreateGpswKeys(context.Background(), pk, msk)
		if err != nil {
			logger.Logger.Errorln("Couldn't store keys for GPSW!", err)
			return err
		}
		return nil
	}
	gpswPublicKey = &GPSWPubKey{}
	gpswSecretKey = data.Vector{}

	if err = utils.DecodeToObject(gpswKeys.PublicKey, gpswPublicKey); err != nil {
		logger.Logger.Errorln("Couldn't decode GPSW public GPSW!", err)
		return err
	}
	if err = utils.DecodeToObject(gpswKeys.MasterSecretKey, &gpswSecretKey); err != nil {
		logger.Logger.Errorln("Couldn't decode GPSW secret key!", err)
		return err
	}
	return nil
}

func DecryptGpswCipher(gpswCipher *GPSWCipher, email string, position string, department string) (
	res string,
	err error,
) {
	defer func() {
		if r := recover(); r != nil {
			logger.Logger.Errorln("Recovered from panic in DecryptGpswCipher:", r)
			res = ""
			err = fmt.Errorf("Recovered from panic in DecryptGpswCipher")
		}
	}()
	policy := strconv.Itoa(AttributesToNum["email:"+email]) + " OR (" + strconv.Itoa(AttributesToNum["role:"+position]) +
		" AND " + strconv.Itoa(AttributesToNum["department:"+department]) + ")"
	msp, err := abe.BooleanToMSP(policy, true)
	if err != nil {
		logger.Logger.Errorln("Couldn't generate MSP for GPSW!", err)
		return "", err
	}
	abeKey, err := gpsw.GeneratePolicyKey(msp, gpswSecretKey)
	if err != nil {
		logger.Logger.Errorln("Failed to generate decryption keys for GPSW!", err)
	}

	res, err = gpsw.Decrypt(gpswCipher, abeKey)
	if err != nil {
		logger.Logger.Errorln("Couldn't decrypt FAME cipher!", err)
		return "", err
	}
	return
}

func EncryptSensitiveInfo(email string, department string, salary int, address string) ([]byte, []byte, error) {

	attrs := []int{
		AttributesToNum["email:"+email], AttributesToNum["department:"+department],
		AttributesToNum["role:manager"],
	}
	salaryString := strconv.Itoa(salary)
	salaryCipher, err := gpsw.Encrypt(salaryString, attrs, gpswPublicKey)
	if err != nil {
		logger.Logger.Errorln("Couldn't encrypt salary!", err)
		return nil, nil, err
	}
	addressCipher, err := gpsw.Encrypt(address, attrs, gpswPublicKey)
	if err != nil {
		logger.Logger.Errorln("Couldn't encrypt address!", err)
		return nil, nil, err
	}
	salaryBytes, err := utils.EncodeToBytes(salaryCipher)
	if err != nil {
		logger.Logger.Errorln("Couldn't encode salary cipher!", err)
		return nil, nil, err
	}
	addressBytes, err := utils.EncodeToBytes(addressCipher)
	if err != nil {
		logger.Logger.Errorln("Couldn't encode address cipher!", err)
		return nil, nil, err
	}
	return salaryBytes, addressBytes, nil
}

var AttributesToNum = map[string]int{
	"role:manager":                0,
	"role:senior":                 1,
	"role:junior":                 2,
	"department:engineering":      3,
	"department:marketing":        4,
	"department:finance":          5,
	"email:emilijaivko@gmail.com": 6,
	"email:helenadp@gmail.com":    7,
	"email:elenav@gmail.com":      8,
	"email:jovanivko@gmail.com":   9,
	"email:matejam@gmail.com":     10,
	"email:markol@gmail.com":      11,
	"email:miljkoca@gmail.com":    12,
	"email:teodorad@gmail.com":    13,
	"email:anastasijad@gmail.com": 14,
	"email:filipi@gmail.com":      15,
	"email:ivanm@gmail.com":       16,
	"email:aleksandram@gmail.com": 17,
}
