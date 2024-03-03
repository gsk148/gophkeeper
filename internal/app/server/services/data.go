package services

import (
	"encoding/json"
	"errors"

	"github.com/gsk148/gophkeeper/internal/app/server/storage"
	"github.com/gsk148/gophkeeper/internal/pkg/enc"
)

func StoreSecureDataFromPayload(db storage.IDataRepository, uid string, payload any, t storage.Type) (string, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	encData, err := enc.EncryptData(data)
	if err != nil {
		return "", err
	}

	sd := storage.SecureData{
		UID:  uid,
		Data: encData,
		Type: t,
	}
	return db.StoreData(sd)
}

func GetDataFromBytes(b []byte, t storage.Type) (any, error) {
	db, err := enc.DecryptData(b)
	if err != nil {
		return nil, err
	}

	return getDataOfType(db, t)
}

func getDataOfType(data []byte, t storage.Type) (any, error) {
	var (
		res any
		err error
	)

	switch t {
	case storage.SBinary:
		var d BinaryRes
		err = json.Unmarshal(data, &d)
		res = d
		break
	case storage.SCard:
		var d CardRes
		err = json.Unmarshal(data, &d)
		res = d
		break
	case storage.SPassword:
		var d PasswordRes
		err = json.Unmarshal(data, &d)
		res = d
		break
	case storage.SText:
		var d TextRes
		err = json.Unmarshal(data, &d)
		res = d
		break
	}

	if err != nil {
		return nil, errors.New("failed to decrypt data")
	}
	return res, nil
}
