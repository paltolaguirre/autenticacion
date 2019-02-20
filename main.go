package main

import ("net/http"
 		"log"
 		//"github.com/xubiosueldos/conexionBD"
)

func main() {

	//db := conexionBD.ConnectBD()

	router := newRouter()
	
	server := http.ListenAndServe(":8081", router)

	log.Fatal(server)

}
