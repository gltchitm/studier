package id

import (
	"crypto/rand"
	"math/big"
)

const (
	IdTypeRegular = "regular"
	IdTypeToken   = "token"
	IdTypeTicket  = IdTypeToken
)

const idTypeRegularAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const idTypeRegularLength = 16

const idTypeTokenAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const idTypeTokenLength = 64

func GenerateIdCustom(alphabet string, length int) (string, error) {
	result := ""
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		result += string(alphabet[num.Int64()])
	}

	return result, nil
}

func GenerateId(idType string) (string, error) {
	var alphabet string
	var length int

	if idType == IdTypeRegular {
		alphabet = idTypeRegularAlphabet
		length = idTypeRegularLength
	} else if idType == IdTypeToken {
		alphabet = idTypeTokenAlphabet
		length = idTypeTokenLength
	}

	return GenerateIdCustom(alphabet, length)
}
