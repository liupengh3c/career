package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// 固定 16 字节密钥（AES-128）
// 你可以改成 AES-256（32字节）
var AES_KEY = []byte("0123456789abcdef")

// 固定 IV（必须 16 字节）
var AES_IV = []byte("abcdef9876543210")

// PKCS7 Padding
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytesRepeat(byte(padding), padding)
	return append(data, padtext...)
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("invalid padding")
	}
	pad := int(data[length-1])
	if pad > length || pad == 0 {
		return nil, errors.New("invalid padding")
	}
	return data[:length-pad], nil
}

func bytesRepeat(b byte, count int) []byte {
	out := make([]byte, count)
	for i := range out {
		out[i] = b
	}
	return out
}

// 加密 (返回固定长度 Base64URL)
func AesEncryptFixedBase64(plaintext []byte) (string, error) {
	// 初始化aes加密器
	block, err := aes.NewCipher(AES_KEY)
	if err != nil {
		return "", err
	}

	plaintext = pkcs7Pad(plaintext, aes.BlockSize)

	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, AES_IV)
	mode.CryptBlocks(ciphertext, plaintext)

	// 使用 Base64URL，长度可控、不会出现 “/” 或 “+”
	return base64.RawURLEncoding.EncodeToString(ciphertext), nil
}

// 解密
func AesDecryptFixedBase64(ciphertextBase64 string) ([]byte, error) {
	ciphertext, err := base64.RawURLEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(AES_KEY)
	if err != nil {
		return nil, err
	}
	fmt.Println(len(ciphertext))
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext not aligned")
	}

	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, AES_IV)
	mode.CryptBlocks(plaintext, ciphertext)

	return pkcs7Unpad(plaintext)
}

// ---- 生成 clip_id：构造固定长度明文 ------------------

func buildClipPlaintext(carID string, startMS int64, endMS int64,
	clipType int, version string) string {
	// 20+1+13+1+13+1+5+1+22=77
	return fmt.Sprintf(
		"%15s$%013d$%013d$%05d$%22s",
		carID, startMS, endMS, clipType, version)
}

// ---- Demo 测试 ----

func main() {
	plaintext := buildClipPlaintext("JME1503", 1763259582000, 1763259602000, 2, "adfm_attach")
	fmt.Println("len", len(plaintext), plaintext)
	enc, err := AesEncryptFixedBase64([]byte(plaintext))
	if err != nil {
		panic(err)
	}

	fmt.Println("Encrypted:", enc)
	fmt.Println("Encrypted Length:", len(enc)) // 固定长度！

	dec, err := AesDecryptFixedBase64(enc)
	if err != nil {
		panic(err)
	}
	origin := strings.Split(string(dec), "$")
	fmt.Println("len", len(origin))

	fmt.Println("Decrypted:", strings.Trim(string(dec), " "))
	car := origin[0]
	start, _ := strconv.Atoi(origin[1])
	end, _ := strconv.Atoi(origin[2])
	clipType, _ := strconv.Atoi(origin[3])
	fmt.Println(car, start, end, clipType)
}
