package main

import (
	"github.com/bitwormhole/starter"
	starterobjectbucket "github.com/bitwormhole/starter-object-bucket"
)

func main() {
	i := starter.InitApp()
	i.Use(starterobjectbucket.ModuleForDemo())
	i.Run()
}
