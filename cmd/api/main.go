package main

import "github.com/GianOrtiz/bean/cmd/api/server"

func main() {
	server, err := server.New()
	if err != nil {
		panic(err)
	}
	if err := server.Run(); err != nil {
		panic(err)
	}
}
