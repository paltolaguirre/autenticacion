package main

import (
	"github.com/jinzhu/gorm"
	"github.com/xubiosueldos/conexionBD/Autenticacion/structAutenticacion"
	"log"
	"net/http"

	"github.com/xubiosueldos/framework/configuracion"
)

func main() {

	configuracion := configuracion.GetInstance()

	router := newRouter()

	server := http.ListenAndServe(":"+configuracion.Puertomicroservicioautenticacion, router)

	log.Fatal(server)

}

func cleanConnections(db *gorm.DB)  {
	db.Model(&structAutenticacion.Security{}).Update("necesitaupdate", true)
}