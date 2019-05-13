package main

import (
	"log"
	"net/http"

	"github.com/xubiosueldos/autenticacion/publico"
	"github.com/xubiosueldos/conexionBD"
)

func main() {

	db := conexionBD.ConnectBD("security")
	db.AutoMigrate(&publico.Security{})

	router := newRouter()

	server := http.ListenAndServe(":8081", router)

	log.Fatal(server)

}
