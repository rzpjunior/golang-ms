package utils

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func RandomInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

func ReplaceSpace(str string) string {
	return strings.ReplaceAll(str, " ", "")
}

func ReplaceUnderscore(str string) string {
	return strings.ReplaceAll(str, " ", "_")
}

func RemoveUnderscore(str string) string {
	return strings.ReplaceAll(str, "_", " ")
}

func ReplaceCodeString(message string, replacer map[string]interface{}) string {
	for k, v := range replacer {
		replacerString := strings.NewReplacer(k, v.(string))
		message = replacerString.Replace(message)
	}

	return message
}

func GenerateRandomDoc(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	t := time.Now().Unix()
	s := strconv.FormatInt(t, 10)
	return s + string(b)
}

func GenerateUUID() string {
	return uuid.New().String()
}

func GenerateUnixTime() string {
	unixTime := time.Now().Unix()          // Get Unix time in seconds
	return strconv.FormatInt(unixTime, 10) // Convert the int64 Unix timestamp to a string
}