package main

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/xubiosueldos/autenticacion/apiclientautenticacion"
	"github.com/xubiosueldos/autenticacion/publico"
	"github.com/xubiosueldos/conexionBD/apiclientconexionbd"
	"github.com/xubiosueldos/framework/configuracion"
)

func main() {
	configuracion := configuracion.GetInstance()
	var tokenAutenticacion publico.Security
	tokenAutenticacion.Tenant = "security"
	versionMicroservicio := obtenerVersionSecurity()
	tenant := apiclientautenticacion.ObtenerTenant(&tokenAutenticacion)
	apiclientconexionbd.ObtenerDB(tenant, "autenticacion", versionMicroservicio, AutomigrateTablasPrivadas)

	router := newRouter()

	server := http.ListenAndServe(":"+configuracion.Puertomicroservicioautenticacion, router)

	log.Fatal(server)

}

func AutomigrateTablasPrivadas(db *gorm.DB) {
	db.AutoMigrate(&publico.Security{})
}

func obtenerVersionSecurity() int {
	configuracion := configuracion.GetInstance()

	return configuracion.Versionsecurity
}
