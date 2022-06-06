package utils

import (
	"math/rand"
	"time"
)

func RandomString(num int) string {
	var letters = []byte("2wert342342erert5324wertw523retw4524")

	result := make([]byte, num)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
