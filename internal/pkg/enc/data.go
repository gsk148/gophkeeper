package enc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

var (
	ErrDecryption = errors.New("enc: failed to decrypt data")
	ErrDataLength = errors.New("enc: the data length is too short for encryption")
)

var secret = []byte("f91j&famF*kf_PgjJ1Yfv$_0f1A8BB#2")

// EncryptData transforms an original slice of bytes into an encoded one.
func EncryptData(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, ErrDataLength
	}

	c, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

// DecryptData transforms an encrypted slice of bytes into an original one.
func DecryptData(data []byte) ([]byte, error) {
	c, err := aes.NewCipher(secret)
	if err != nil {
		return nil, ErrDecryption
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, ErrDecryption
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, ErrDataLength
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	res, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrDecryption
	}
	return res, nil
}
