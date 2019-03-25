package main

import (
	"fmt"
    "net/http"
   // "log"
  //  "github.com/gorilla/mux"
    "encoding/json"
  	//"github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
  	"time"
 //   "github.com/xubiosueldos/conexionBD"
  	"github.com/xubiosueldos/autenticacion/publico"
)





//var db *gorm.DB
var err error
var m = make(map[string] publico.TokenAutenticacion)
var errors publico.Error

func Login(w http.ResponseWriter, r *http.Request){

	username := r.FormValue("username")
	pass := r.FormValue("pass")
	tenant := r.FormValue("tenant")

	token := time.Now().Format("2006-01-02 15:04:05.000000")
	fecha := time.Now()

	autenticacion := publico.TokenAutenticacion{Username: username, Pass: pass, Tenant: tenant, Token: token, FechaCreacion: fecha}

	//VER COMO PREGUNTAR SI EL TOKEN YA ESTA INGRESADO EN EL HASHMAP
	m[string(token)] = autenticacion

	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(200)

	json.NewEncoder(w).Encode(token)
	fmt.Println(m[string(token)])
	fmt.Println(autenticacion)
}


func CheckToken(w http.ResponseWriter, r *http.Request){

	header := r.Header

	token := header.Get("Token")

	fmt.Println(m[token])

	tokenAutenticacion,ok := m[token]

	if(ok){
		w.Header().Set("Content-Type", "application-json")
		w.WriteHeader(200)

		json.NewEncoder(w).Encode(tokenAutenticacion)
	} else 
	{
		errors = publico.Error{ErrorNombre: "Hubo error", ErrorCodigo: 400}

		w.Header().Set("Content-Type", "application-json")
		w.WriteHeader(errors.ErrorCodigo)

		json.NewEncoder(w).Encode(errors)
	}
}


