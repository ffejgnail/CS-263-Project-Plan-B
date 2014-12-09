package main

import (
	"fmt"
	"os"
)

func main() {
	env := NewEnvironment()
	for i := 0; i < Iteration; i++ {
		env.Run(i)
		if ObserveGene {
			for j := 0; j < 6; j++ {
				fmt.Print(env.Aggressiveness[j])
				fmt.Print("\t")
			}
			fmt.Println(env.Aggressiveness[6])
		} else {
			fmt.Printf("%d\n", env.Aggressiveness[0])
		}
	}

	if RecordGIF {
		f, err := os.Create("new.gif")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		env.WriteTo(f)
	}
}
