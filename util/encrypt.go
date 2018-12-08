package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"github.com/mohemohe/butimili-api/configs"
	"io"
)

func Encrypt(s string) *string {
	plainText := []byte(s)
	cipherText := make([]byte, aes.BlockSize+len(s))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		Logger().Error(err)
		return nil
	}
	block, err := getBlock()
	if err != nil {
		Logger().Error(err)
		return nil
	}
	encryptStream := cipher.NewCTR(block, iv)
	encryptStream.XORKeyStream(cipherText[aes.BlockSize:], plainText)
	result := hex.EncodeToString(cipherText)
	return &result
}

func Decrypt(s string) (result *string) {
	defer func() {
		err := recover()
		if err != nil {
			result = nil
		}
		return
	}()

	plainText, err := hex.DecodeString(s)
	if err != nil {
		Logger().Error(err)
		return nil
	}

	cipherText := plainText[aes.BlockSize:]
	iv := plainText[:aes.BlockSize]
	block, err := getBlock()
	if err != nil {
		Logger().Error(err)
		return nil
	}

	decryptedText := make([]byte, len(cipherText))
	decryptStream := cipher.NewCTR(block, iv)
	decryptStream.XORKeyStream(decryptedText, cipherText)

	rawText := string(decryptedText)
	return &rawText
}

func getBlock() (cipher.Block, error) {
	secret := configs.GetEnv().Encrypt.Secret
	key := []byte(secret)
	return aes.NewCipher(key)
}
