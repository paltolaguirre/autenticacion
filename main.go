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

	db_public := conexionBD.ObtenerDB("public")
	err = apiclientconexionbd.AutomigrateTablasPublicas(db_public)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	conexionBD.CerrarDB(db_public)

	db := conexionBD.ObtenerDB("security")
	err = apiclientconexionbd.AutomigrateTablaSecurity(db)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	router := newRouter()

	server := http.ListenAndServe(":"+configuracion.Puertomicroservicioautenticacion, router)

	log.Fatal(server)

}
