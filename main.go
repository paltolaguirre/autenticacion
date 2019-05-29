package main

import (
	"log"
	"net/http"

	"github.com/xubiosueldos/autenticacion/publico"
	"github.com/xubiosueldos/conexionBD"
	"github.com/xubiosueldos/framework/configuracion"
)

func main() {
	configuracion := configuracion.GetInstance()
	db := conexionBD.ConnectBD("security")
	db.AutoMigrate(&publico.Security{})

	router := newRouter()

	server := http.ListenAndServe(":"+configuracion.Puertomicroserivicioautenticacion, router)

	log.Fatal(server)

}
