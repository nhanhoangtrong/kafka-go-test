package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("go", "run", "lib/list_topic.go")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		log.Fatal("Run error ", err)
	}
	fmt.Printf("output: \n%s\n", out.String())
}
