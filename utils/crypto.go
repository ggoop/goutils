package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
)

func getKey(key string) []byte {
	//AES-128。key长度：16, 24, 32 bytes 对应 AES-128, AES-192, AES-256
	l := 32
	if len(key) >= 32 {
		l = 32
	} else if len(key) >= 24 {
		l = 24
	} else {
		l = 16
	}
	ctx := md5.New()
	ctx.Write([]byte(key))
	return ([]byte(hex.EncodeToString(ctx.Sum(nil))))[:l]
}
func Encrypt(text, key string) (string, error) {
	skey := getKey(key)
	var iv = skey[:aes.BlockSize]
	encrypted := make([]byte, len(text))
	block, err := aes.NewCipher(skey)
	if err != nil {
		return "", err
	}
	encrypter := cipher.NewCFBEncrypter(block, iv)
	encrypter.XORKeyStream(encrypted, []byte(text))
	return hex.EncodeToString(encrypted), nil
}
func Decrypt(encrypted, key string) (string, error) {
	skey := getKey(key)
	var err error
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	src, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	var iv = skey[:aes.BlockSize]
	decrypted := make([]byte, len(src))
	var block cipher.Block
	block, err = aes.NewCipher(skey)
	if err != nil {
		return "", err
	}
	decrypter := cipher.NewCFBDecrypter(block, iv)
	decrypter.XORKeyStream(decrypted, src)
	return string(decrypted), nil
}
