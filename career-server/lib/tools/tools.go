package tools

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CreateDir(dir string) bool {
	_, err := os.Stat(dir)
	if err == nil {
		return true
	}
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return false
	}
	return true
}
func DeleteDir(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		// log.Fatal(dir + "delete error")
	}
	return
}
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		// log.Fatal("error")
	}
	return strings.Replace(dir, "\\", "/", -1)
}
func GetParentDir() string {
	path, _ := os.Getwd()
	pos := strings.LastIndex(path, "/")
	runes := []rune(path)
	parPath := string(runes[:pos])
	return parPath
}

/*
data: byte类型
生成对应的base64编码
*/
func Base64Encode(data []byte) string {
	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	base64 := b64.EncodeToString(data)
	return base64
}

/*
data: byte类型
对base64编码的数据进行解码
*/
func Base64Decode(data string) ([]byte, error) {
	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	bStr, err := b64.DecodeString(data)
	if err != nil {
		return bStr, err
	}
	return bStr, err
}

/*
data: byte类型
生成对应的md5值
*/
func Md5Sum(data []byte) string {
	mdFive := md5.New()
	mdFive.Write(data)
	md5Str := hex.EncodeToString(mdFive.Sum(nil))
	return md5Str
}

/*
生成随机数，以当前时间戳（ns）作为源，保证每次的随机数不同
*/
func RandU32() []byte {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	bt := make([]byte, 5)
	binary.BigEndian.PutUint32(bt, rand.Uint32())
	return bt
}

/*
生成随机数，以当前时间戳（ns）作为源，保证每次的随机数不同
*/
func RandInt(n int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	bt := rand.Intn(n)
	return bt
}

/**
 * [UnZip解压zip文件]
 * @Author   liupeng17
 * @DateTime 2019-04-17 16:42
 * @param    [string]     zipFile    [要解压的文件]
 * @param    [string]     dest       [解压到路径]
 * @return   [error]      err        [返回值错误信息]
 */
func UnZip(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		filename := dest + file.Name
		fmt.Println("filename=", filename)
		rs := []rune(filename)
		localPath := string(rs[0:strings.LastIndex(filename, "/")])
		err = os.MkdirAll(localPath, 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		w.Close()
		rc.Close()
	}
	return nil
}

func AesDecrypt(decodeStr, key, iv string) ([]byte, error) {
	//先解密base64
	decodeBytes, err := Base64Decode(decodeStr)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	origData := make([]byte, len(decodeBytes))

	blockMode.CryptBlocks(origData, decodeBytes)
	origData = pKCS7UnPadding(origData)
	return origData, nil
}
func AesEncrypt(origData, key []byte, iv string) (string, error) {
	b64Str := ""
	block, err := aes.NewCipher(key)
	if err != nil {
		return b64Str, err
	}
	blockSize := block.BlockSize()
	origData = pKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	b64Str = Base64Encode(crypted)
	return b64Str, nil
}

func pKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
func SendHttpRequest(ctx context.Context, urlStr string, method string, headers map[string]string, body string, params url.Values) (int64, string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return 0, "", err
	}

	// 如果有查询参数，添加到URL中
	if params != nil {
		u.RawQuery = params.Encode()
	}
	// 处理body
	payload := strings.NewReader(body)

	client := &http.Client{}

	// 创建HTTP请求
	req, err := http.NewRequest(method, u.String(), payload)
	if err != nil {
		return 0, "", err
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// 发送请求并获取响应
	res, err := client.Do(req)
	if err != nil {
		return int64(res.StatusCode), "", err
	}
	defer res.Body.Close()

	// 读取响应体
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return int64(res.StatusCode), "", err
	}

	return int64(res.StatusCode), string(bodyBytes), nil
}
