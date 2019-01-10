package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func INIT() {
	prop := os.Args[1:]
	if len(prop) == 0 {
		fmt.Print(helpTemplate)
		return
	}

	if !getCheckPatch() {
		return
	}

	if strings.ToLower(prop[0]) == "help" {
		fmt.Print(helpTemplate)
		return

	}
	if strings.ToLower(prop[0]) == "new" {
		if len(prop) >= 2 {
			makeProject(prop[1])
		}
		return

	}
	if strings.ToLower(prop[0]) == "install" {
		if len(prop) >= 2 {
			out, _ := exec.Command("go", "get", prop[0]).Output()
			fmt.Println(string(out))
		} else {
			exec.Command("go", "get", "github.com/gorilla/mux").Output()
			exec.Command("go", "get", "gopkg.in/mgo.v2").Output()

		}

		return

	}
	if strings.ToLower(prop[0]) == "make:model" {
		if len(prop) >= 2 {
			makeModel(prop[1])
		} else {
			fmt.Println("No Find Model Name")

		}

		return

	}
	if strings.ToLower(prop[0]) == "patch" {

		fmt.Println("your patch : " + getPatch())
		fmt.Println("GoPatch : " + getGoPatch())
		return

	}
	if strings.ToLower(prop[0]) == "make:controller" {
		makeController(prop[1], "", fmt.Sprintf(controllerFuncTemplate, "Index", "//Auto by gn"))
		return

	}
	if strings.ToLower(prop[0]) == "make:rest" {
		if len(prop) >= 1 {
			makeREST(prop[1])
		} else {
			log.Print("errore you no send two props")
		}
		return

	}

	fmt.Print(helpTemplate)
}
