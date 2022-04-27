package padding

import "bytes"

type PKCS7 struct {
	blockSize int
}

func NewPKCS7(blockSize int) *PKCS7 {
	return &PKCS7{blockSize: blockSize}
}

func (p *PKCS7) UnPad(encryptedData []byte) []byte {
	length := len(encryptedData)
	unpadding := int(encryptedData[length-1])
	return encryptedData[:(length - unpadding)]
}

func (p *PKCS7) Pad(plainData []byte) []byte {
	padding := p.blockSize - len(plainData)%p.blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plainData, padtext...)
}
