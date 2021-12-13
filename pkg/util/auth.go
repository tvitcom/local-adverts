package util

import (
	"github.com/alexedwards/argon2id"
	"golang.org/x/crypto/bcrypt"
	"encoding/base64"
	rnd "math/rand"
	"crypto/sha256"
	"crypto/rand"
	"crypto/md5"
	"fmt"
)

var alfabetHex = []rune("0123456789abcdef")
var alfabetSimple = []rune("123456789abcdefghijkmnopqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ_-+=!@")

func RandomHexString(n int) string { // without O,0,l symbols
    b := make([]rune, n)
    for i := range b {
        b[i] = alfabetHex[rnd.Intn(len(alfabetHex))]
    }
    return string(b)
}

func RandomUsableString(n int) string { // without O,0,l symbols
    b := make([]rune, n)
    for i := range b {
        b[i] = alfabetSimple[rnd.Intn(len(alfabetSimple))]
    }
    return string(b)
}


// Google OAUTH clients functions
func RandToken32bit() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func GetMD5Hash(text string) string {
	data := []byte(text)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func GetSha256(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func GetSaltedSha256(salt, password string) string {
	return GetSha256(salt + password)
}

func IsEqualSaltedSha256(salt, password, hashDatabase string) bool {
	saltedPasswordHash := GetSaltedSha256(salt, password)
	if saltedPasswordHash == hashDatabase {
		return true
	}
	return false
}

func MakeArgonHash(password string) (string, error) {
	params := &argon2id.Params{
		Memory:      128 * 1024,
		Iterations:  10,
		Parallelism: 4,
		SaltLength:  16,
		KeyLength:   32,
	}
	return argon2id.CreateHash(password, params)
}

func MatchArgonHashAndPassword(trypass, passhash string) bool {
	match, err := argon2id.ComparePasswordAndHash(trypass, passhash)
	if err != nil {
		return false
	}
	return match
}

// Generate return a hashed password
func MakeBCryptHash(raw string, cost int) string {
    hash, err := bcrypt.GenerateFromPassword([]byte(raw), cost) 
    if err != nil {
        panic(err)
    }   
    return string(hash)
}

// Verify compares a hashed password with plaintext password
func VerifyBCrypt(testPassword, hash string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(testPassword))
}
