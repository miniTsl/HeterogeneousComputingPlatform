package main

import "fmt"
import "github.com/spf13/cobra"

import(
	"HCPlatform/cmd"
)

func main(){
	fmt.Println("Hello World!")
	cmd.execute();
}