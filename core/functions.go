package core

import "log"

func ErrHandler(err error) {
	if err != nil {
		log.Println(err)
	}
	return
}