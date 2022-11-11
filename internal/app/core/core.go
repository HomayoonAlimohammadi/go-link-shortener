package core

import (
	"math/rand"
	"strings"
	"time"
)

var (
	validCharacters string
)

func init() {
	validCharacters = "abcdefghijklmnopqrstuvwxyz"
	validCharacters += strings.ToUpper(validCharacters)
	validCharacters += "0123456789"
}

func GenerateToken(maxTokenLength int) string {
	rand.Seed(time.Now().UnixNano())
	var token string
	for i := 0; i < maxTokenLength; i++ {
		randomNumber := rand.Intn(len(validCharacters))
		token += string(validCharacters[randomNumber])
	}
	return token
}
