package main

import (
	"log"
	"net/http"

	"github.com/xubiosueldos/conexionBD/apiclientconexionbd"
	"github.com/xubiosueldos/framework/configuracion"
)

func main() {
	configuracion := configuracion.GetInstance()

	db_public := apiclientconexionbd.ObtenerDB("public")
	apiclientconexionbd.AutomigrateTablasPublicas(db_public)
	apiclientconexionbd.CerrarDB(db_public)

	db := apiclientconexionbd.ObtenerDB("security")
	apiclientconexionbd.AutomigrateTablaSecurity(db)

	router := newRouter()

	server := http.ListenAndServe(":"+configuracion.Puertomicroservicioautenticacion, router)

	log.Fatal(server)

}
