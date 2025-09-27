package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("python3", "-c", "print('Hello from Python!')")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing Python code:", err)
		return
	}
	fmt.Print(string(out))
}
