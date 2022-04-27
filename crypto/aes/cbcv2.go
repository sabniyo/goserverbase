package aes

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"

	"github.com/sabariramc/goserverbase/crypto"
	"github.com/sabariramc/goserverbase/crypto/padding"
	"github.com/sabariramc/goserverbase/log"
)

var ErrIVLengthMismatch = fmt.Errorf("IV lenght is not matching with block size")

type AESCBCV2 struct {
	padder crypto.Padder
	key    []byte
	log    *log.Logger
	iv     []byte
}

func NewAESCBCV2PKCS7ConstantIV(ctx context.Context, log *log.Logger, key string, iv []byte) (*AESCBCV2, error) {
	keyByte, err := getKey(key)
	if err != nil {
		log.Error(ctx, "Error creating AES CBC V2", err)
		return nil, fmt.Errorf("crypto.aes.NewAESCBCV2PKCS7ConstantIV: %w", err)
	}
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return nil, fmt.Errorf("crypto.aes.ChiperV2: %w", err)
	}
	if len(iv) != block.BlockSize() {
		return nil, fmt.Errorf("crypto.aes.ChiperV2: %w", ErrIVLengthMismatch)
	}
	return NewAESCBCConstantIV(ctx, log, key, dup(iv), padding.NewPKCS7(block.BlockSize()))
}

func NewAESCBCConstantIV(ctx context.Context, log *log.Logger, key string, iv []byte, padder crypto.Padder) (*AESCBCV2, error) {
	keyByte, err := getKey(key)
	if err != nil {
		log.Error(ctx, "Error creating AES CBC V2", err)
		return nil, fmt.Errorf("crypto.aes.NewAESCBCV2ConstantIV: %w", err)
	}
	return &AESCBCV2{key: keyByte, padder: padder, log: log}, nil
}

func (a *AESCBCV2) Encrypt(plainBlob []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, fmt.Errorf("AESCBCV2.Encrypt: %w", err)
	}
	paddedData := a.padder.Pad(plainBlob)
	blockModel := cipher.NewCBCEncrypter(block, a.iv)
	cipherBlob := make([]byte, len(paddedData))
	blockModel.CryptBlocks(cipherBlob, paddedData)
	return cipherBlob, nil
}

func (a *AESCBCV2) EncryptString(plainText string) (string, error) {
	res, err := a.Encrypt([]byte(plainText))
	if err != nil {
		return "", fmt.Errorf("AESCBCV2.EncryptString: %w", err)
	}
	return base64.StdEncoding.EncodeToString(res), nil
}

func (a *AESCBCV2) Decrypt(encryptedData []byte) (plainData []byte, err error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, fmt.Errorf("AESCBCV2.Decrypt: %w", err)
	}
	blockModel := cipher.NewCBCDecrypter(block, a.iv)
	plainData = make([]byte, len(encryptedData))
	blockModel.CryptBlocks(plainData, encryptedData)
	plainData = a.padder.UnPad(plainData)
	return plainData, nil
}

func (a *AESCBCV2) DecryptString(plainText string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(plainText)
	if err != nil {
		return "", fmt.Errorf("AESCBCV2.DecryptString.B64Decode: %w", err)
	}
	res, err := a.Decrypt([]byte(decoded))
	if err != nil {
		return "", fmt.Errorf("AESCBCV2.EncryptString: %w", err)
	}
	return string(res), nil
}

func dup(p []byte) []byte {
	q := make([]byte, len(p))
	copy(q, p)
	return q
}
