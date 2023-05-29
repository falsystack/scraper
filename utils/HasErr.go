package utils

import "log"

func HasErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
