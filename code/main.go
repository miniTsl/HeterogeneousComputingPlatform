package main

import "HCPlatform/code/cmd"

func main() {
	err := cmd.Execute()
	if err != nil {
		return
	}

}
