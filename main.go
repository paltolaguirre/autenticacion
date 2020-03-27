package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/xubiosueldos/conexionBD"
	"github.com/xubiosueldos/conexionBD/apiclientconexionbd"
	"github.com/xubiosueldos/framework/configuracion"
)

func main() {
	var err error
	configuracion := configuracion.GetInstance()

	dbPublic := conexionBD.ObtenerDB("public")
	txPublic := dbPublic.Begin()
	err = apiclientconexionbd.AutomigrateTablasPublicas(txPublic)
	if err != nil {
		txPublic.Rollback()
		fmt.Println("Error Public Automigrate: ", err)
		return
	}
	txPublic.Commit()
	conexionBD.CerrarDB(dbPublic)

	dbSecurity := conexionBD.ObtenerDB("security")
	txSecurity := dbSecurity.Begin()
	err = apiclientconexionbd.AutomigrateTablaSecurity(txSecurity)
	if err != nil {
		txSecurity.Rollback()
		fmt.Println("Error Security Automigrate: ", err)
		return
	}
	txSecurity.Commit()
	conexionBD.CerrarDB(dbSecurity)

	router := newRouter()

	server := http.ListenAndServe(":"+configuracion.Puertomicroservicioautenticacion, router)

	log.Fatal(server)

}
