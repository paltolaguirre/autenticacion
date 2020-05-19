module github.com/xubiosueldos/autenticacion

go 1.12

replace github.com/xubiosueldos/conexionBD => /home/wschmidt/go/src/github.com/xubiosueldos/conexionBD

require (
	github.com/gorilla/mux v1.7.2
	github.com/jinzhu/gorm v1.9.8
	github.com/xubiosueldos/conexionBD v1.1.11
	github.com/xubiosueldos/framework v1.1.2
	github.com/xubiosueldos/monoliticComunication v1.1.0
)
