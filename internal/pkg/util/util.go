package util

import "log"

func PanicIfError(err error, message string) {
	if err != nil {
		log.Fatalln(message, "-", err)
	}
}
