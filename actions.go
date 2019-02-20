package main

import (
	"fmt"
    "net/http"
   // "log"
  //  "github.com/gorilla/mux"
 //   "encoding/json"
  	//"github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
  	"time"
 //   "github.com/xubiosueldos/conexionBD"
  	"github.com/xubiosueldos/autenticacion/publico"
)





//var db *gorm.DB
var err error
var m = make(map[string] publico.TokenAutenticacion)
var empty struct{}

func Login(w http.ResponseWriter, r *http.Request){

	username := r.FormValue("username")
	pass := r.FormValue("pass")
	tenant := r.FormValue("tenant")

	token := time.Now().Format("2006-01-02 15:04:05.000000")
	fecha := time.Now()

	autenticacion := publico.TokenAutenticacion{Username: username, Pass: pass, Tenant: tenant, Token: token, FechaCreacion: fecha}

	//VER COMO PREGUNTAR SI EL TOKEN YA ESTA INGRESADO EN EL HASHMAP
	m[string(token)] = autenticacion
	fmt.Println(m[string(token)])
	fmt.Println(autenticacion)
}


func CheckToken(r*http.Request)(bool){

	header := r.Header

	token := header.Get("Token")

	fmt.Println(m[token])

	_,ok := m[token]
	
	return ok
}


