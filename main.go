package main

import "fmt"
import "github.com/rsalmond/kissyface/cmd"

func main() {
	err := cmd.Analyze()

	if err != nil {
		fmt.Println(err)
	}
}
