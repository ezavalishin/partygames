package utils

import (
	"log"
	"os"
	"strconv"
)

func MustGet(k string) string {
	v := os.Getenv(k)

	if v == "" {
		log.Panicln("ENV missing, key: " + k)
	}

	return v
}

func MustGetBool(k string) bool {
	v := os.Getenv(k)
	if v == "" {
		log.Panicln("ENV missing, key: " + k)
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		log.Panicln("ENV err: [" + k + "]\n" + err.Error())
	}
	return b
}
