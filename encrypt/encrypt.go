package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"hash"
	"io"

	"github.com/huyungtang/go-lib/strings"
	"golang.org/x/crypto/bcrypt"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

const (
	DefaCost encryptCost = 16
	MidCost  encryptCost = 24
	MaxCost  encryptCost = 32

	preLen = 5
)

var (
	errInvalidFormat = errors.New("invalid format of encrypted string")
)

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Encrypt
// ****************************************************************************************************************************************
func Encrypt(str string, cost encryptCost) (enc string, err error) {
	ct := make([]byte, aes.BlockSize+len(str))
	iv := ct[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	var blk cipher.Block
	key := strings.Random(cost)
	if blk, err = aes.NewCipher([]byte(key)); err != nil {
		return
	}

	cipher.NewCFBEncrypter(blk, iv).
		XORKeyStream(ct[aes.BlockSize:], []byte(str))

	return strings.Format("$a%02d$%s%s", cost, key, base64.URLEncoding.EncodeToString(ct)), nil
}

// Decrypt
// ****************************************************************************************************************************************
func Decrypt(ori string) (str string, err error) {
	var cs string
	if cs = strings.Find(ori, `^\$a(16|24|32)\$`); cs == "" {
		return ori, errInvalidFormat
	}

	var cv int
	switch cs {
	case "$a16$":
		cv = int(DefaCost)
	case "$a24$":
		cv = int(MidCost)
	case "$a32$":
		cv = int(MaxCost)
	}

	if len(ori) <= (cv + preLen) {
		return ori, errInvalidFormat
	}

	key := ori[preLen : cv+preLen]
	str = ori[cv+preLen:]

	var ct []byte
	if ct, err = base64.URLEncoding.DecodeString(str); err != nil {
		return ori, err
	}

	var blk cipher.Block
	if blk, err = aes.NewCipher([]byte(key)); err != nil {
		return ori, err
	}

	iv := ct[:aes.BlockSize]
	ct = ct[aes.BlockSize:]

	cipher.NewCFBDecrypter(blk, iv).
		XORKeyStream(ct, ct)

	return string(ct), nil
}

// Base64
// ****************************************************************************************************************************************
func Base64(str string) string {

	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Base64Decrypt
// ****************************************************************************************************************************************
func Base64Decrypt(str string) (s string, err error) {
	var bs []byte
	if bs, err = base64.StdEncoding.DecodeString(str); err != nil {
		return
	}

	return string(bs), nil
}

// BCrypt
// ****************************************************************************************************************************************
func BCrypt(str string, cost int) (enc string, err error) {
	switch {
	case cost < bcrypt.MinCost:
		cost = bcrypt.MinCost
	case cost > bcrypt.MaxCost:
		cost = bcrypt.MaxCost
	}

	var bs []byte
	if bs, err = bcrypt.GenerateFromPassword([]byte(str), cost); err != nil {
		return
	}

	return string(bs), nil
}

// BCryptVerify
// ****************************************************************************************************************************************
func BCryptVerify(hash, str string) (err error) {

	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
}

// MD5
// ****************************************************************************************************************************************
func MD5(str string) (enc string, err error) {

	return cryptoEncrypt(md5.New(), str)
}

// MD5Verify
// ****************************************************************************************************************************************
func MD5Verify(hash, str string) (err error) {

	return cryptoEncryptVerify(md5.New(), hash, str)
}

// SHA256
// ****************************************************************************************************************************************
func SHA256(str string) (enc string, err error) {

	return cryptoEncrypt(sha256.New(), str)
}

// SHA256Verify
// ****************************************************************************************************************************************
func SHA256Verify(hash, str string) (err error) {

	return cryptoEncryptVerify(sha256.New(), hash, str)
}

// SHA512
// ****************************************************************************************************************************************
func SHA512(str string) (enc string, err error) {

	return cryptoEncrypt(sha512.New(), str)
}

// SHA512Verify
// ****************************************************************************************************************************************
func SHA512Verify(hash, str string) (err error) {

	return cryptoEncryptVerify(sha512.New(), hash, str)
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// encryptCost ****************************************************************************************************************************
type encryptCost = int

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// cryptoEncrypt **************************************************************************************************************************
func cryptoEncrypt(hash hash.Hash, str string) (enc string, err error) {
	if _, err = hash.Write([]byte(str)); err != nil {
		return
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// cryptoEncryptVerify ********************************************************************************************************************
func cryptoEncryptVerify(hash hash.Hash, hstr, str string) (err error) {
	if str, err = cryptoEncrypt(hash, str); err != nil {
		return
	}

	if hstr != str {
		err = errors.New("encryption validity failed")
	}

	return
}
