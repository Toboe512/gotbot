package utils

import (
	"crypto/aes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math"
	"math/rand"
	"time"
)

const blockSize = 16

func StringToSha256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)

	return fmt.Sprintf("%x\n", bs)
}

func EncryptStrAes(pwd string, src string) (res string, slt string, err error) {
	srcB := normalize16Block([]byte(src))
	key, salt := genKey(pwd)

	bc, err := aes.NewCipher(key)

	if err != nil {
		return EmptyStr, EmptyStr, err
	}

	result := make([]byte, len(srcB))
	for i := 0; i < len(srcB); i += blockSize {
		var dst = make([]byte, blockSize)
		bc.Encrypt(dst, srcB[i:i+blockSize])
		copy(result[i:], dst)
	}
	sltStr := base64.StdEncoding.EncodeToString(salt)
	return string(result), sltStr, nil
}

func DecryptStrAes(pwd string, slt string, src string) (res string, err error) {
	srcB := normalize16Block([]byte(src))
	pwdB := []byte(pwd)
	sltB, err := base64.StdEncoding.DecodeString(slt)

	if err != nil {
		return EmptyStr, err
	}
	key := make([]byte, len(pwdB)+len(sltB))
	copy(key, pwdB)
	copy(key[len(pwdB):], sltB)

	bc, err := aes.NewCipher(key)

	if err != nil {
		return EmptyStr, err
	}

	result := make([]byte, len(srcB))
	for i := 0; i < len(srcB); i += blockSize {
		var dst = make([]byte, blockSize)
		bc.Decrypt(dst, srcB[i:i+blockSize])
		copy(result[i:], dst)
	}
	return string(trimArrByEndZero(result)), nil
}

func genKey(pwd string) (res []byte, slt []byte) {
	pwdB := []byte(pwd)
	lenSalt := blockSize - len(pwdB)%blockSize
	if lenSalt == blockSize {
		return pwdB, nil
	}
	resultSize := lenSalt + len(pwdB)
	result := make([]byte, len(pwdB), resultSize)
	copy(result, pwdB)

	salt := make([]byte, 0, lenSalt)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i <= lenSalt-1; i++ {
		rndB := byte(r.Intn(math.MaxUint8))
		salt = append(salt, rndB)
		result = append(result, rndB)
	}

	return result, salt
}

func normalize16Block(src []byte) []byte {
	lenSalt := blockSize - len(src)%blockSize
	if lenSalt == blockSize {
		return src
	}

	resultSize := lenSalt + len(src)
	result := make([]byte, resultSize)
	copy(result, src)

	salt := make([]byte, 0, lenSalt)
	copy(result[lenSalt:], salt)

	return result
}

func trimArrByEndZero(src []byte) []byte {
	var result []byte
	for _, b := range src {
		if b != 0 {
			result = append(result, b)
		}
	}
	return result
}
