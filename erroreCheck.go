package main

import (
	"fmt"
	"os/exec"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func fixImport(patchName string) {
	//fmt.Println(getPatch() + patchName)
	fmt.Println("fix file: " + getPatch() + patchName)
	exec.Command("goimports", "-w", getPatch()+patchName).Output()
	exec.Command("gofmt", "-w", getPatch()+patchName).Output()
}
