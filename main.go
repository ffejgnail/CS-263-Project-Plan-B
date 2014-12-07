package main

import (
	"fmt"
	"os"
)

func main() {
	env := NewEnvironment()
	for i := 0; i < Iteration; i++ {
		env.Run(i)
	}
	f, err := os.Create("new.gif")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	env.WriteTo(f)
}