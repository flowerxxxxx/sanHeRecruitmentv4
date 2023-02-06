package wechatUtil

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"log"
)

func Dncrypt(rawData, key, iv string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(rawData)
	key_b, err_1 := base64.StdEncoding.DecodeString(key)
	iv_b, _ := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return "", err
	}
	if err_1 != nil {
		return "", err_1
	}
	dnData, err := AesCBCDncrypt(data, key_b, iv_b)
	if err != nil {
		return "", err
	}
	return string(dnData), nil
}

// 解密
func AesCBCDncrypt(encryptData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err)
	}
	blockSize := block.BlockSize()
	if len(encryptData) < blockSize {
		log.Println("AesCBCDncrypt ciphertext too short")
	}
	if len(encryptData)%blockSize != 0 {
		log.Println("AesCBCDncrypt ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encryptData, encryptData)
	// 解填充
	encryptData = PKCS7UnPadding(encryptData)
	return encryptData, nil
}

// 去除填充
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
