package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	key        = []byte("123EFGJA143EDSBE") // 加密的密钥
	iv         = []byte("ABCDEF1234123412") // 偏移
	timeFormat = "2006-01-02 15:04:05"
	nowTime    = time.Now() // 当日时间
)

func SsoAes() {
	var encrypt struct {
		Timestamp int64  `json:"timestamp"`
		AppName   string `json:"app_name"`
	}
	encrypt.Timestamp = nowTime.Unix()
	encrypt.AppName = "tms"
	b, err := json.Marshal(encrypt)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("------------------ CBC模式 --------------------")
	encrypteda := SSOAesEncryptCBC(key, iv, b)
	fmt.Println("SSO加密结果：", encrypteda)
	decrypted, _ := SSOAesDecryptCBC(key, iv, encrypteda)
	fmt.Println("SSO解密结果：", string(decrypted))
}

func TmsAes() {
	fmt.Println("------------------ CBC模式 --------------------")
	b := "13211112222"
	encrypteda := AesEncryptCBC(key, iv, []byte(b))
	fmt.Println("TMS加密结果：", encrypteda)
	//encrypteda = "MTUxYmRhMDRiYzY0ZjA4MTJlNTRjYjYw\n"
	decrypted, _ := AesDecryptCBC(key, iv, encrypteda)
	fmt.Println("TMS解密结果：", decrypted)
}

// AesEncryptCBC AES
func AesEncryptCBC(key, iv []byte, origData []byte) string {
	block, _ := aes.NewCipher(key)               // 创建一个cipher.Block接口
	blockSize := block.BlockSize()               // 获取秘钥块的长度
	origData = PKCS7Padding(origData, blockSize) // 补全码
	blockMode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(origData))   // 创建数组
	blockMode.CryptBlocks(encrypted, origData) // 加密
	dataAes := hex.EncodeToString(encrypted)
	dataBase := base64.StdEncoding.EncodeToString([]byte(dataAes))
	return dataBase
}

// AesDecryptCBC AES
func AesDecryptCBC(key, iv []byte, encrypted string) (string, error) {
	if len(encrypted) == 0 {
		return "", errors.New(" str decode fail: string is empty")
	}
	dataBase, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", errors.New(" base64 decode fail")
	}
	dataAes, err := hex.DecodeString(string(dataBase))
	if err != nil {
		return "", errors.New(" hex decode fail")
	}
	block, err := aes.NewCipher(key) // 分组秘钥
	if err != nil {
		return "", errors.New(" get block fail")
	}
	blockMode := cipher.NewCBCDecrypter(block, iv) // 加密模式
	decrypted := make([]byte, len(dataAes))        // 创建数组
	blockMode.CryptBlocks(decrypted, dataAes)      // 解密
	decrypted = PKCS7UnPadding(decrypted)          // 去除补全码
	return string(decrypted), nil
}

// AesEncryptCBC SSO--AES加密
func SSOAesEncryptCBC(key, iv []byte, encryptStr []byte) string {
	block, _ := aes.NewCipher(key)                  // 创建一个cipher.Block接口
	blockSize := block.BlockSize()                  // 获取秘钥块的长度
	origData := PKCS7Padding(encryptStr, blockSize) // 补全码
	blockMode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(origData))   // 创建数组
	blockMode.CryptBlocks(encrypted, origData) // 加密
	dataBase := base64.StdEncoding.EncodeToString(encrypted)
	return dataBase
}

// SSOAesDecryptCBC SSO--AES解密
func SSOAesDecryptCBC(key, iv []byte, encrypted string) (string, error) {
	if len(encrypted) == 0 {
		return "", errors.New(" str decode fail: string is empty")
	}
	dataBase, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", errors.New(" base64 decode fail")
	}
	block, _ := aes.NewCipher(key) // 分组秘钥
	if err != nil {
		return "", errors.New(" get block fail")
	}
	blockMode := cipher.NewCBCDecrypter(block, iv) // 加密模式
	decrypted := make([]byte, len(dataBase))       // 创建数组
	blockMode.CryptBlocks(decrypted, dataBase)     // 解密
	decrypted = PKCS7UnPadding(decrypted)          // 去除补全码
	return string(decrypted), nil
}

//使用PKCS7进行填充
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
