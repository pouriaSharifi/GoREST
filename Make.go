package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func makeREST(name string) {
	x := getModelFields(name)
	if x == "" {
		log.Print("Failed make REST API Model Not Find")
		return
	}
	mainModel := strings.Split(x, ";")[1:]
	myModel := mainModel[0 : len(mainModel)-1]
	myModelName := mainModel[len(mainModel)-1]
	ModelName := strings.Split(myModelName, "=")[0]
	makeSeterForModel(ModelName, myModel)
	var controllerBodyTemplate string

	controllerBodyTemplate += fmt.Sprintf(controllerFuncTemplate, "Create"+ModelName, makeCreateMethodForController(ModelName, myModel))
	controllerBodyTemplate += fmt.Sprintf(controllerFuncTemplate, "Update"+ModelName, makeUpdateMethodForController(ModelName, myModel))
	controllerBodyTemplate += fmt.Sprintf(controllerFuncTemplate, "Get"+ModelName+"s", makeGetAllDataMethodForController(ModelName))
	controllerBodyTemplate += fmt.Sprintf(controllerFuncTemplate, "DestroyByID"+ModelName, makeDestroyByIDMethodForController(ModelName))
	controllerBodyTemplate += fmt.Sprintf(controllerFuncTemplate, "GetById"+ModelName, makeGetByIdMethodForController(ModelName))
	controllerBodyTemplate += fmt.Sprintf(validationMethod, "valid"+ModelName, validModelBody)
	makeController(ModelName, "", controllerBodyTemplate)
	var ModelRepository string
	ModelRepository += fmt.Sprintf(RepositoryGetAllMethod, ModelName, ModelName, ModelName, ModelName, ModelName, ModelName, ModelName, ModelName)
	ModelRepository += fmt.Sprintf(RepositoryCreateMethod, ModelName, ModelName, ModelName, ModelName)
	ModelRepository += fmt.Sprintf(RepositoryDestroyByID, ModelName, ModelName)
	ModelRepository += fmt.Sprintf(RepositoryDropCollection, ModelName, ModelName)
	ModelRepository += fmt.Sprintf(RepositoryUpdate, ModelName, ModelName, ModelName, ModelName)
	ModelRepository += fmt.Sprintf(RepositoryUpdateByModel, ModelName, ModelName, ModelName, ModelName)
	ModelRepository += fmt.Sprintf(RepositoryGetById, ModelName, ModelName, ModelName)

	makeCreateRepositoryForModel(ModelName, ModelRepository)
	makeRoute(strings.Title(ModelName))
}
func makeRoute(name string) {
	patch := fmt.Sprintf("%s\\Route\\%sControllerRoutes.go", getPatch(), strings.Title(name))
	if _, err := os.Stat(patch); os.IsExist(err) {
		log.Print("Failed make  File Set" + name + ".go")
		return
	}

	file, _ := os.Create(patch)

	var TemplateRoute string
	TemplateRoute += fmt.Sprintf(muxRouteTemplate, "Create"+name, name+"Controller", "Create"+name, "Create"+name)
	TemplateRoute += fmt.Sprintf(muxRouteTemplate, "DestroyByID"+name+"/{id}", name+"Controller", "DestroyByID"+name, "DestroyByID"+name)
	TemplateRoute += fmt.Sprintf(muxRouteTemplate, "Get"+name+"s", name+"Controller", "Get"+name+"s", "Get"+name+"s")
	TemplateRoute += fmt.Sprintf(muxRouteTemplate, "Update"+name, name+"Controller", "Update"+name, "Update"+name)
	TemplateRoute += fmt.Sprintf(muxRouteTemplate, "GetById"+name+"/{id}", name+"Controller", "GetById"+name, "GetById"+name)
	file.WriteString(fmt.Sprintf(muxRouteFuncTemplate, strings.Title(name), TemplateRoute))

	file.Close()
	fixImport(fmt.Sprintf("\\Route\\%sControllerRoutes.go", strings.Title(name)))
}
func makeModel(name string) {
	patch := fmt.Sprintf("%s\\Models\\%s.go", getPatch(), strings.Title(name))
	if _, err := os.Stat(patch); os.IsExist(err) {
		log.Print("Failed make  File Set" + name + ".go")
		return
	}

	file, _ := os.Create(patch)

	file.WriteString(fmt.Sprintf(makeModelpackageTemplate, strings.Title(name)))

	file.Close()
	fixImport(fmt.Sprintf("\\Models\\%s.go", strings.Title(name)))

}
func makeCreateRepositoryForModel(ModelName string, FileBode string) {
	patch := fmt.Sprintf("%s\\Models\\%sRepository.go", getPatch(), ModelName)
	if _, err := os.Stat(patch); os.IsExist(err) {
		log.Print("Failed make  File Set" + ModelName + ".go")
		return
	}
	//file, _ := os.OpenFile(patch, os.O_APPEND|os.O_WRONLY, 0600)~
	file, _ := os.Create(patch)

	file.WriteString(fmt.Sprintf(packageTemplate, "Models", FileBode))
	file.WriteString("//=>" + fmt.Sprintf("%s\\Models\\%sRepository.go", getPatch(), ModelName))

	file.Close()
	fixImport(fmt.Sprintf("\\Models\\%sRepository.go", ModelName))

}
func makeGetAllDataMethodForController(ModelName string) string {
	return fmt.Sprintf(GetAllDataMethodForControllerTemplate, ModelName, ModelName, ModelName, ModelName)
}
func makeDestroyByIDMethodForController(ModelName string) string {
	return fmt.Sprintf(DestroyByIDMethodForController, ModelName)
}
func makeGetByIdMethodForController(ModelName string) string {
	return fmt.Sprintf(GetByIdMethodForController, ModelName, ModelName, ModelName, ModelName, ModelName, ModelName, ModelName, ModelName, ModelName, ModelName)
}
func makeCreateMethodForController(ModelName string, myModel []string) string {
	var bodyTemplateForCreate string
	bodyTemplateForCreate += fmt.Sprintf("var %s = Models.%s{}\n", ModelName+"Model", ModelName)
	bodyTemplateForCreate += `
	var   (
		checkValidation bool
		err string
	) 

	`
	for _, item := range myModel {
		value := strings.Split(item, "=")[0]
		bodyTemplateForCreate += fmt.Sprintf(`checkValidation, err = %s.Set%s(r.FormValue("%s"))`+"\n", ModelName+"Model", value, value)
		bodyTemplateForCreate += fmt.Sprintf(checkValidation, ModelName)
	}
	bodyTemplateForCreate += fmt.Sprintf(`isDestroy := %s.Create%s()`+"\n", ModelName+"Model", ModelName)
	bodyTemplateForCreate += fmt.Sprintf(isDestroyTemplate, ModelName, ModelName)
	return bodyTemplateForCreate
}
func makeUpdateMethodForController(ModelName string, myModel []string) string {
	var bodyTemplateForCreate string
	bodyTemplateForCreate += fmt.Sprintf("var %s = Models.%s{}\n", ModelName+"Model", ModelName)
	bodyTemplateForCreate += `
	var   (
		checkValidation bool
		err string
	) 

	`
	for _, item := range myModel {
		value := strings.Split(item, "=")[0]
		bodyTemplateForCreate += fmt.Sprintf(`checkValidation, err = %s.Set%s(r.FormValue("%s"))`+"\n", ModelName+"Model", value, value)
		bodyTemplateForCreate += fmt.Sprintf(checkValidation, ModelName)
	}
	bodyTemplateForCreate += fmt.Sprintf(`isDestroy := %s.Update%sModel()`+"\n", ModelName+"Model", ModelName)
	bodyTemplateForCreate += fmt.Sprintf(isDestroyUpdateTemplate, ModelName, ModelName)
	return bodyTemplateForCreate
}
func makeSeterForModel(ModelName string, myModel []string) {
	patch := fmt.Sprintf("%s\\Models\\Set%s.go", getPatch(), ModelName)
	if _, err := os.Stat(patch); os.IsExist(err) {
		log.Print("Faild make  File Set" + ModelName + ".go")
		return
	}
	//file, _ := os.OpenFile(patch, os.O_APPEND|os.O_WRONLY, 0600)~
	file, _ := os.Create(patch)

	file.WriteString("package Models\n\n")
	for _, item := range myModel {
		//file.WriteString("your value: " + strings.Split(item, "=")[0] + "\nyour format: " + strings.Split(item, "=")[1] + "\n")
		fieldName := strings.Split(item, "=")[0]
		fieldFormat := strings.Split(item, "=")[1]
		template := getTemplateByFormat(fieldFormat, fieldName, ModelName)
		file.WriteString(fmt.Sprintf(setValueForModelTemplate, ModelName+"Model", ModelName, fieldName, fieldName, fieldFormat, template))

	}
	file.Close()
	fixImport("\\Models\\Set" + ModelName + ".go")
}
func makeController(name string, patchName string, body string) {
	name = strings.Title(name)
	//check dir exist
	if _, err := os.Stat(getPatch() + patchName + "\\Controller\\" + name + "Controller"); os.IsNotExist(err) {
		os.Mkdir(getPatch()+patchName+"\\Controller\\"+name+"Controller", os.ModePerm)
		//fmt.Println("Create Controller " + name)
	}
	//check file exist
	if _, err := os.Stat(getPatch() + patchName + "\\Controller\\" + name + "Controller\\" + name + ".go"); err == nil {
		//panic(err)
		fmt.Println("Controller Exist!!")
		fmt.Println("Create Controller stop!!")
		return
	}
	f, err := os.Create(getPatch() + patchName + "\\Controller\\" + name + "Controller\\" + name + ".go")
	check(err)
	_, err = f.WriteString(fmt.Sprintf(controllerTemplate, name, body))
	f.Close()
	fixImport(patchName + "\\Controller\\" + name + "Controller\\" + name + ".go")

}
func createMain(patch string) {

	f, err := os.Create(getPatch() + "\\" + patch + "\\main.go")
	check(err)
	_, err = f.WriteString(MainFileTemplate)
	f.Close()
	fixImport("\\" + patch + "\\main.go")

}
func makeProject(name string) {
	if _, err := os.Stat(getPatch() + "\\" + name); os.IsNotExist(err) {
		root := getPatch() + "\\" + name
		os.Mkdir(root, os.ModePerm)
		os.Mkdir(root+"\\Controller", os.ModePerm)
		os.Mkdir(root+"\\Controller\\HomeController", os.ModePerm)
		os.Mkdir(root+"\\Models", os.ModePerm)
		os.Mkdir(root+"\\Views", os.ModePerm)
		os.Mkdir(root+"\\assets", os.ModePerm)
		os.Mkdir(root+"\\Route", os.ModePerm)
		os.Mkdir(root+"\\Middleware", os.ModePerm)
		os.Mkdir(root+"\\Views\\Home", os.ModePerm)

		makeController("Home", "\\"+name, fmt.Sprintf(controllerFuncTemplate, "Index", "//Auto by gn /n"+controllerIndexHomeTemplate))

		HtmlCreatorText, _ := os.Create(root + "\\Views\\Home\\index.html")
		HtmlCreatorText.WriteString(HomeIndexTemplate)
		HtmlCreatorText.Close()

		HtmlCreatorTextAsset, _ := os.Create(root + "\\assets\\index.html")
		HtmlCreatorTextAsset.WriteString(assetsIndexTemplate)
		HtmlCreatorTextAsset.Close()

		ModelCreatorText, _ := os.Create(root + "\\Models\\Database.go")
		ModelCreatorText.WriteString(databaseTemplate)
		ModelCreatorText.Close()

		MiddlewareCreatorText, _ := os.Create(root + "\\Middleware\\Middleware.go")
		MiddlewareCreatorText.WriteString(middlewareTemplate)
		MiddlewareCreatorText.Close()

		MiddlewareForMUXCreatorText, _ := os.Create(root + "\\Middleware\\MiddlewareForMUX.go")
		MiddlewareForMUXCreatorText.WriteString(MiddlewareForMUXTemplate)
		MiddlewareForMUXCreatorText.Close()

		createMain(name)

		fixImport("\\" + name + "\\Models\\Database.go")
		fixImport("\\" + name + "\\Middleware\\Middleware.go")
		fixImport("\\" + name + "\\Middleware\\MiddlewareForMUX.go")

	} else {
		panic(err)
	}

}
