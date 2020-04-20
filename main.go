package main

import (
	"bitbucket.org/telemetryapp/go_interview/router"
)

func main() {
	r := router.New()

	r.Run()
}
