package store

import (
	"abeProofOfConcept/pkg/store/models"
	"context"
)

type Key = models.MasterKey

func GetFameKeys(ctx context.Context) (*Key, error) {
	key := new(Key)
	err := DB.NewSelect().Model(key).Where("scheme = 'fame'").Scan(ctx)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func GetGpswKeys(ctx context.Context) (*Key, error) {
	key := new(Key)
	err := DB.NewSelect().Model(key).Where("scheme = 'gpsw'").Scan(ctx)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func CreateFameKeys(ctx context.Context, publicKey []byte, msk []byte) (*Key, error) {
	key := &Key{
		Scheme:          "fame",
		PublicKey:       publicKey,
		MasterSecretKey: msk,
	}
	_, err := DB.NewInsert().Model(key).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func CreateGpswKeys(ctx context.Context, publicKey []byte, secretKey []byte) (*Key, error) {
	key := &Key{
		Scheme:          "gpsw",
		PublicKey:       publicKey,
		MasterSecretKey: secretKey,
	}
	_, err := DB.NewInsert().Model(key).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return key, nil
}
