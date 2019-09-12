package main

import (
	"log"
	"net/http"

	"github.com/xubiosueldos/conexionBD"
	"github.com/xubiosueldos/conexionBD/apiclientconexionbd"
	"github.com/xubiosueldos/framework/configuracion"
)

func main() {
	configuracion := configuracion.GetInstance()
	db := apiclientconexionbd.ObtenerDB("public")

	db := conexionBD.ConnectBD("security")
	conexionBD.AutomigrateTablaSecurity(db)
	router := newRouter()

	server := http.ListenAndServe(":"+configuracion.Puertomicroservicioautenticacion, router)

	log.Fatal(server)

}
