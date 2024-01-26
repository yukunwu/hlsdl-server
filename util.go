package main

import "log"

func Recover() {
	if r := recover(); r != nil {
		log.Printf("Runtime error caught: %v", r)
	}
}
