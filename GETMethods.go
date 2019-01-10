package main

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func getPatch() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func getModelFields(name string) string {
	randName := fmt.Sprintf("%s", time.Now().Format("20060102150405"))
	tempDirName := randName + "temp"
	rempDirPatch := fmt.Sprintf("%s\\%stemp\\temp.go", getPatch(), randName)
	os.Mkdir(tempDirName, os.ModeDir)
	f, err := os.Create(rempDirPatch)
	check(err)
	_, err = f.WriteString(fmt.Sprintf(getModelFieldsTemplate, name, "%+v", "%+v"))
	f.Sync()
	f.Close()
	fixImport("\\" + randName + "temp\\temp.go")
	//exec.Command("goimports", "-w", getPatch()+).Output()
	//exec.Command("gofmt", "-w", getPatch()+"\\temp.go").Output()
	out, _ := exec.Command("go", "run", rempDirPatch).Output()
	os.Remove(rempDirPatch)
	os.Remove(getPatch() + "\\" + tempDirName)
	return string(out)

}

func getGoPatch() string {

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return gopath
}
func getCheckPatch() bool {
	x := strings.Index(getPatch(), getGoPatch()+"\\src")

	if x != -1 {

		return true
	}
	fmt.Println("Error\nyour go patch is => " + getGoPatch() + "\\src\name" + "no Run gn in " + getPatch())
	return false
}
func getTemplateByFormat(format string, fieldName string, ModelName string) string {
	if format == "bson.ObjectId" {
		return fmt.Sprintf(bsonObjectTemplate, fieldName, "%q", fieldName, ModelName+"Model", fieldName)
	} else if format == "string" {
		return fmt.Sprintf(stringTemplate, ModelName+"Model", fieldName, fieldName)
	} else if format == "float64" {
		return fmt.Sprintf(floatTemplate, fieldName, "64", ModelName+"Model", fieldName, fieldName)
	} else if format == "float32" {
		return fmt.Sprintf(floatTemplate, fieldName, "32", ModelName+"Model", fieldName, fieldName)
	} else if format == "int" {
		return fmt.Sprintf(intTemplate, fieldName, ModelName+"Model", fieldName, fieldName)
	} else {
		return `return false,"error Message: you need change method ` + "set" + fieldName + `() validation in file model/set` + ModelName + `.go"`
	}
}
