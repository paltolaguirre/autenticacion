package main

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	s "strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/xubiosueldos/autenticacion/publico"
	"github.com/xubiosueldos/conexionBD"
	"github.com/xubiosueldos/framework"
	"github.com/xubiosueldos/framework/configuracion"
)

//var db *gorm.DB
var err error
var errors publico.Error

func Login(w http.ResponseWriter, r *http.Request) {

	var datosCorrectos bool = true

	tokenEncode := obtenerTokenHeader(r)
	configuracion := configuracion.GetInstance()

	if configuracion.Checkmonolitico == true {
		datosCorrectos = chequeoAuthenticationMonolitico(tokenEncode, r)
	}
	//Chequear con el monolitico que los datos ingresados sean correctos
	if datosCorrectos {

		tokenDecode, err := base64.StdEncoding.DecodeString(tokenEncode)

		if err != nil {
			framework.RespondError(w, http.StatusNotFound, err.Error())
			return
		}

		//token := time.Now().Format("2006-01-02 15:04:05.000000")

		security := insertarTokenSecurity(tokenDecode, w)

		framework.RespondJSON(w, http.StatusOK, security)

	} else {
		framework.RespondError(w, http.StatusNotFound, framework.DatosIncorrectos)
		return
	}

}

func Logout(w http.ResponseWriter, r *http.Request) {

	token := obtenerTokenHeader(r)

	db := conexionBD.ConnectBD("security")
	defer db.Close()

	if err := db.Unscoped().Where("token = ?", token).Delete(publico.Security{}).Error; err != nil {

		framework.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//Hice que devuelva el token, no se si es necesario
	framework.RespondJSON(w, http.StatusOK, token)
}

func CheckToken(w http.ResponseWriter, r *http.Request) {

	token := obtenerTokenHeader(r)

	security, ok, err := checkTokenDB(w, token)

	if ok {
		framework.RespondJSON(w, http.StatusOK, security)

	} else {
		//errors = publico.Error{ErrorNombre: "Hubo error", ErrorCodigo: http.StatusUnauthorized}
		framework.RespondError(w, http.StatusUnauthorized, err.Error())
		return
	}
}

func chequeoAuthenticationMonolitico(tokenEncode string, r *http.Request) bool {

	infoUserValida := false
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	var prueba []byte = []byte("xubiosueldosimplementadocongo")
	tokenSecurity := base64.StdEncoding.EncodeToString(prueba)

	url := configuracion.GetUrlMonolitico() + "SecurityAuthenticationGo"

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Add("Authorization", tokenEncode)
	req.Header.Add("SecurityToken", tokenSecurity)

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusAccepted {
		infoUserValida = true
	}

	return infoUserValida
}

func insertarTokenSecurity(tokenDecode []byte, w http.ResponseWriter) *publico.Security {

	db := conexionBD.ConnectBD("security")
	defer db.Close()

	infoUser := s.Split(string(tokenDecode), ":")

	username := infoUser[0]
	pass := infoUser[1]
	tenant := infoUser[2]

	numeroRandom := rand.Int63()
	token := strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.FormatInt(numeroRandom, 10)

	fecha := time.Now()

	security := publico.Security{Username: username, Pass: pass, Tenant: tenant, Token: token, FechaCreacion: fecha}

	if err := db.Create(&security).Error; err != nil {
		framework.RespondError(w, http.StatusInternalServerError, err.Error())
		return nil
	}

	return &security
}
func checkTokenDB(w http.ResponseWriter, token string) (*publico.Security, bool, error) {

	var existeToken bool = true
	var security publico.Security
	var err error = nil
	db := conexionBD.ConnectBD("security")

	if err = db.Set("gorm:auto_preload", true).First(&security, "token = ?", token).Error; gorm.IsRecordNotFoundError(err) {
		framework.RespondError(w, http.StatusUnauthorized, err.Error())
		existeToken = false
	}

	return &security, existeToken, err
}

func obtenerTokenHeader(r *http.Request) string {

	header := r.Header.Get("Authorization")

	token := s.Split(header, " ")[1]

	return token

}
