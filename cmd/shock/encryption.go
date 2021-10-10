/*

*/
package main

import (
	"io"
	"log"
	"os"

	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
)

// convenience to use a known token as the state value
func encryptState(plain string) string {
	// encrypt plain value 
	phrase := os.Getenv("AMBER_PHRASE")
	buf := encrypt([]byte(plain), phrase)
	return hex.EncodeToString(buf)
}
func decryptState(ciphertxt string) string {
	// take hex val and decrypt to plain txt
	phrase := os.Getenv("AMBER_PHRASE")
	data, err := hex.DecodeString(ciphertxt)
	if err != nil {
		log.Fatal(err)
	}
	buf := decrypt(data, phrase)
	return string(buf)
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err)
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Fatal(err)
	}
	return plaintext
}


