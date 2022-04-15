package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/huyungtang/go-lib/strings"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Encrypt
// ****************************************************************************************************************************************
func Encrypt(str string, cost int) (enc string, err error) {
	switch {
	case cost > 24:
		cost = 32
	case cost > 16:
		cost = 24
	default:
		cost = 16
	}

	cipherText := make([]byte, aes.BlockSize+len(str))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	var block cipher.Block
	key := strings.Random(cost)
	if block, err = aes.NewCipher([]byte(key)); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], []byte(str))

	return strings.Format("$a%d$%s%s", cost, key, base64.URLEncoding.EncodeToString(cipherText)), nil
}

// Decrypt
// ****************************************************************************************************************************************
func Decrypt(hash string) (str string, err error) {
	var cost int
	switch strings.Find(`^\$a(16|24|32)\$`, hash) {
	case "$a16$":
		cost = 16
	case "$a24$":
		cost = 24
	case "$a32$":
		cost = 32
	default:
		return hash, ErrEncryptFormat
	}

	hash = hash[5:]
	if len(hash) < cost {
		return hash, ErrEncryptFormat
	}

	var cipherText []byte
	if cipherText, err = base64.URLEncoding.DecodeString(hash[cost:]); err != nil {
		return
	}

	var block cipher.Block
	if block, err = aes.NewCipher([]byte(hash[0:cost])); err != nil {
		return
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

// EncryptBase64
// ****************************************************************************************************************************************
func EncryptBase64(str string) string {

	return base64.StdEncoding.EncodeToString([]byte(str))
}

// DecryptBase64
// ****************************************************************************************************************************************
func DecryptBase64(str string) string {
	bs, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return str
	}

	return string(bs)
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
