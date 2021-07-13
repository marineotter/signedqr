package main

import "time"

func main() {
	context := time.Now().Format("2006-0102-150405")
	generatekey(context)
	sign(context)
}
