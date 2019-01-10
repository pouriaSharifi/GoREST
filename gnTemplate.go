package main

var controllerTemplate = `
package %sController



//Auto Generate


%s

`
var packageTemplate = `
package %s



//Auto Generate


%s

`
var RepositoryGetAllMethod = `



func (Model *%s) Get%ss() *[]%s {
	var %ss []%s

	Session_.C("%ss").Find(nil).All(&%ss)
	return &%ss
}

`
var MainFileTemplate = `



package main

 
func main() {

	var session = Models.SessionMgo()
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	Models.Session_ = session.DB("DatabaseName")
	//Models.SetKey("TableName", []string{"_id"})  set new key to table

	//########################################Router########################################\\
	r := mux.NewRouter()
	r.HandleFunc("/", HomeController.Index).Methods("GET")
	//########################################Router########################################\\

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./assets"))))
	http.ListenAndServe(":8080", r)


}


`
var RepositoryCreateMethod = `



func (Model *%s) Create%s() bool {

	if x := Session_.C("%ss").Insert(&Model); x == nil {
		return true
	} else {
		log.Printf("Method:Create%s => Error:%%+v", x)
		return false
	}

}

`
var RepositoryDestroyByID = `

func (Model *%s) DestroyByID(data string) bool {
	d, err := hex.DecodeString(data)
	if err != nil || len(d) != 12 {
		log.Print(fmt.Sprintf("Method:DestroyByID => Errore:invalid input to ObjectIdHex: %%q", data))
		return false
	}

	if x := Session_.C("%ss").Remove(bson.M{"_id": bson.ObjectIdHex(data)}); x == nil {
		return true
	} else {
		log.Printf("Method:DestroyByID => Errore:%%+v", x)
		return false
	}

}

`
var RepositoryDropCollection = `
func (Model *%s) DropCollection() bool {
	if x := Session_.C("%ss").DropCollection(); x == nil {
		return true
	} else {
		log.Printf("Method:DropCollection => Errore:%%+v", x)
		return false
	}

}


`
var RepositoryUpdate = `
func (Model *%s) Update%sByBson(QueryBson bson.M) bool { 
    //bson.M{"name": "new Name"}
	if x := Session_.C("%ss").Update(bson.M{"_id": Model.Id}, bson.M{"$set":QueryBson}); x == nil {
		return true
	} else {
		log.Printf("Method:Update%s => Error:%%+v", x)
		return false
	}

}


`
var RepositoryUpdateByModel = `

func (Model *%s) Update%sModel() bool { 
	 
	if x := Session_.C("%ss").Update(bson.M{"_id": Model.Id}, Model); x == nil {
		return true
	} else {
		log.Printf("Method:Update%s => Error:%%+v", x)
		return false
	}

}


`
var RepositoryGetById = `
func (Model *%s) Get%s() (error) {

	return Session_.C("%ss").FindId(Model.Id).One(Model)
}


`
var makeModelpackageTemplate = `
package Models

type %s struct {

	Id  bson.ObjectId ` + "`" + `json:"_id" bson:"_id,omitempty" ` + "`" + `

}



`
var controllerFuncTemplate = `

func %s (w http.ResponseWriter, r *http.Request) {

	%s

}

`
var controllerIndexHomeTemplate = `

tmpl := template.Must(template.ParseFiles("Views/Home/index.html"))
	tmpl.Execute(w, struct {
		Name string
	}{Name:"index"})
`
var databaseTemplate = `

package Models

var Session_ *mgo.Database

func SessionMgo() mgo.Session {
	const (
		hosts    = "Ip:port"
		username = ""
		password = ""
	)

	info := &mgo.DialInfo{
		Addrs:    []string{hosts},
		Timeout:  20 * time.Second,
		Username: username,
		Password: password,
	}

	session, err := mgo.DialWithInfo(info)
	if err != nil {
		panic(err)

	}
	return *session

}
func SetKey(table string, keys []string) {
	c := Session_.C(table)
	index := mgo.Index{
		Key:    keys,
		Unique: true,
	}
	c.EnsureIndex(index)
}

`
var middlewareTemplate = `
package Middleware

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}
func Method(m string) Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			if r.Method != m {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}
func Logging() Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			start := time.Now()
			defer func() { log.Println(r.URL.Path, time.Since(start)) }()

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

`
var MiddlewareForMUXTemplate = `
package Middleware

type authenticationMiddleware struct {

}

func (amw *authenticationMiddleware) LoggingMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		defer func() { log.Println(r.URL.Path, time.Since(start)) }()

		next.ServeHTTP(w, r)

	})
}
func (amw *authenticationMiddleware) Jwt(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/","/Api/Token"}
		requestPath := r.URL.Path

		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized

			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")

			return
		}



		token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
			return []byte("#%$#^T%TGRSFASDFR8@EHYTH%Y^E%GR7TARGETED5RAF%E$%WR5F"), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual

			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")

			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server

			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")

			return
		}
		context.Set(r,"user",token.Claims)
		next.ServeHTTP(w, r)

	})
}



var Auth = &authenticationMiddleware{}



`
var getModelFieldsTemplate = `
	package main

	func main(){
		lo := %s{}
		var reply interface{} = lo
        t := reflect.TypeOf(reply)
		for i := 0; i < t.NumField(); i++ {
			fmt.Printf(";%+v=%+v", t.Field(i).Name, t.Field(i).Type)
		}
		fmt.Print(";"+t.Name()+"="+"modelName")	
	}
	`
var setValueForModelTemplate = `
func (%s *%s) Set%s(x string) (bool,string) {

	//L.%s = x	you need convert format string x to %s

	%s	

}

`
var bsonObjectTemplate = `
	if len(x) <= 0 {
		return false,"error for field %s "
	}
	data, err := hex.DecodeString(x)
		if err != nil || len(data) != 12 {
			log.Print(fmt.Sprintf("Method:SetBson => Errore:invalid input to ObjectIdHex: %s", x))
			return false,"error for field %s"
		}
		%s.%s = bson.ObjectIdHex(x)
		return true,""

`
var stringTemplate = `
	if len(x) > 0 {
		%s.%s = x
		return true,""
	}
	return false,"error for field %s "
`
var floatTemplate = `
	if len(x) <= 0 {
		return false,"error for field %s "
	}
	if convert, err := strconv.ParseFloat(x, %s); err == nil {
	
			%s.%s = convert
			return true,""
	}
	return false,"error for field %s "
`
var intTemplate = `
	if len(x) <= 0 {
		return false,"error for field %s "
	}
	if convert, err := strconv.Atoi(x); err == nil {

		%s.%s = convert
		return true,""
	}
	return false,"error for field %s "
`
var validModelBody = `

json.NewEncoder(writer).Encode(struct {
		Status     bool   ` + "`" + `json:"status"` + "`" + `
		StatusCode int   ` + "`" + `json:"statusCode"` + "`" + `
		Msg        string ` + "`" + `json:"msg"` + "`" + `
	}{Msg: msg, Status: status, StatusCode: statusCode})

`
var validationMethod = `

func %s(writer http.ResponseWriter, msg string, status bool, statusCode int) {

	
	%s


}

`
var GetAllDataMethodForControllerTemplate = `

	var x = Models.%s{}
	json.NewEncoder(w).Encode(struct {
		Status     bool   ` + "`" + `json:"status"` + "`" + `
		StatusCode int    ` + "`" + `json:"statusCode"` + "`" + `
		Msg        string ` + "`" + `json:"msg"` + "`" + `
		Data []Models.%s ` + "`" + `json:"data"` + "`" + `
	}{Msg: "all %s data sent to you", Status: true, StatusCode: 200,Data: *x.Get%ss()})


`
var DestroyByIDMethodForController = `

    route := mux.Vars(r)
	var Model = Models.%s{}

	isDestroy := Model.DestroyByID(route["id"])

	if isDestroy == true {
		json.NewEncoder(w).Encode(struct {
			Msg bool  ` + "`" + `json:"msg"` + "`" + `
		}{Msg: true})
	} else {
		json.NewEncoder(w).Encode(struct {
			Msg bool  ` + "`" + `json:"msg"` + "`" + `
		}{Msg: false})
	}



`
var HomeIndexTemplate = `

<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{.Name}}</title>
</head>
<body>
Created By Gn View {{.Name}}
</body>
</html>

`
var assetsIndexTemplate = `

<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Static File</title>
</head>
<body>
Created By Gn Open me as Link http://127.0.0.1:8080/static/
</body>
</html>

`
var GetByIdMethodForController = `

    	var (
		checkValidation bool
		err             string
	)
    Vars := mux.Vars(r)
	var %sModel = &Models.%s{}
	checkValidation, err = %sModel.SetId(Vars["id"])

	if !checkValidation {
		valid%s(w, err, false, 500)
		return
	}
	if  err:=%sModel.Get%s();err!=nil{
		valid%s(w, err.Error(), false, 404)
		return
	}

	json.NewEncoder(w).Encode(struct {
		Status     bool        ` + "`" + `json:"status"` + "`" + `
		StatusCode int         ` + "`" + `json:"statusCode"` + "`" + `
		Msg        string      ` + "`" + `json:"msg"` + "`" + `
		Data       Models.%s  ` + "`" + `json:"data"` + "`" + `
	}{Msg: "%s data sent to you", Status: true, StatusCode: 200, Data: *%sModel })




`
var muxRouteTemplate = `
r.HandleFunc("/Api/%s", Middleware.Chain(%s.%s,Middleware.InRoles("%s"))).Methods("GET")
`
var isDestroyTemplate = `

if isDestroy == true {
	 
		valid%s(w,"data create in database",true,200)
	} else {
		 
		valid%s(w,"data no create in database",false,500)

	}

`
var muxRouteFuncTemplate = `
package Route

//Generated By Gn

func %sControllerRoutes(r *mux.Router) {
	 
 
		%s

}


`
var isDestroyUpdateTemplate = `

if isDestroy == true {
	 
		valid%s(w,"Data Update in database",true,200)
	} else {
		 
		valid%s(w,"Data no Update in database",false,500)

	}

`
var checkValidation = `

if !checkValidation {
		valid%s(w, err, false, 500)
		return
	}

`
var helpTemplate = `
YOUR WELCOME TO GN
new		projectName		#create new project in goPatch
############################################ MAKE HELP ############################################

make:controller		controllerName		#create new controller in controller folder you need create controller folder
make:rest		modelName		#create RestFull api crud use mongodb
make:model		modelName		#create New Model

############################################ package HELP ############################################

install					#install your need package.
install		goPackageName		#install by package name.

`
