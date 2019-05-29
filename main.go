package main

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/xubiosueldos/autenticacion/publico"
	"github.com/xubiosueldos/conexionBD/apiclientconexionbd"
	"github.com/xubiosueldos/framework/configuracion"
)

func main() {
	var tokenAutenticacion publico.Security
	tokenAutenticacion.Tenant = "security"
	versionMicroservicio := obtenerVersionSecurity()
	apiclientconexionbd.ObtenerDB(&tokenAutenticacion, "autenticacion", versionMicroservicio, AutomigrateTablasPrivadas)

	router := newRouter()

	server := http.ListenAndServe(":"+configuracion.Puertomicroserivicioautenticacion, router)

	log.Fatal(server)

}

func AutomigrateTablasPrivadas(db *gorm.DB) {
	db.AutoMigrate(&publico.Security{})
}

func obtenerVersionSecurity() int {
	configuracion := configuracion.GetInstance()

	return configuracion.Versionsecurity
}
