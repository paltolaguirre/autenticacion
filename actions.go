package main

import (
	"encoding/base64"
	"math/rand"
	"net/http"
	"strconv"
	s "strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/xubio-inc/sueldos-lib-conexionBD"
	"github.com/xubio-inc/sueldos-lib-conexionBD/Autenticacion/structAutenticacion"
	"github.com/xubio-inc/sueldos-lib-framework"
	"github.com/xubio-inc/sueldos-lib-framework/configuracion"
	"github.com/xubio-inc/sueldos-lib-monolitico"
)

//var db *gorm.DB
var err error

// Sirve para controlar si el server esta OK
func Healthy(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Healthy."))
}

func Login(w http.ResponseWriter, r *http.Request) {

	var datosCorrectos bool = true

	tokenEncode := obtenerTokenHeader(r)
	configuracion := configuracion.GetInstance()

	if configuracion.Checkmonolitico == true {

		datosCorrectos = monoliticComunication.CheckAuthenticationMonolitico(tokenEncode, r)
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

		err = Actualizar(security)
		if err != nil {
			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		framework.RespondJSON(w, http.StatusOK, security)

	} else {
		framework.RespondError(w, http.StatusNotFound, framework.DatosIncorrectos)
		return
	}

}

func Logout(w http.ResponseWriter, r *http.Request) {

	token := obtenerTokenHeader(r)

	db := conexionBD.ObtenerDB("security")
	//defer db.Close()
	defer conexionBD.CerrarDB(db)

	if err := db.Unscoped().Where("token = ?", token).Delete(structAutenticacion.Security{}).Error; err != nil {

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

func insertarTokenSecurity(tokenDecode []byte, w http.ResponseWriter) *structAutenticacion.Security {

	db := conexionBD.ObtenerDB("security")
	//defer db.Close()
	defer conexionBD.CerrarDB(db)

	infoUser := s.Split(string(tokenDecode), ":")

	username := infoUser[0]
	//pass := infoUser[1]
	tenant := infoUser[2]

	numeroRandom := rand.Int63()
	token := strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.FormatInt(numeroRandom, 10)

	fecha := time.Now()

	security := structAutenticacion.Security{Username: username, Tenant: tenant, Token: token, FechaCreacion: fecha, Necesitaupdate: true}

	if err := db.Create(&security).Error; err != nil {
		framework.RespondError(w, http.StatusInternalServerError, err.Error())
		return nil
	}

	return &security
}
func checkTokenDB(w http.ResponseWriter, token string) (*structAutenticacion.Security, bool, error) {

	var security structAutenticacion.Security
	var err error = nil

	dbSecurity := conexionBD.ObtenerDB("security")
	defer conexionBD.CerrarDB(dbSecurity)

	if err = dbSecurity.Set("gorm:auto_preload", true).First(&security, "token = ?", token).Error; gorm.IsRecordNotFoundError(err) {
		return nil, false, err
	}

	if security.Necesitaupdate {
		err = Actualizar(&security)
		if err != nil {
			return nil, false, err
		}
	}

	return &security, true, err
}

func obtenerTokenHeader(r *http.Request) string {

	header := r.Header.Get("Authorization")

	token := s.Split(header, " ")[1]

	return token

}

func LoginWithNoCheck(w http.ResponseWriter, r *http.Request) {

	tokenEncode := obtenerTokenHeader(r)
	configuracion := configuracion.GetInstance()

	tokenDecode, err := base64.StdEncoding.DecodeString(tokenEncode)

	if err != nil {
		framework.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	splited := s.Split(string(tokenDecode), ":")

	if configuracion.Codesecurity != splited[1] {
		framework.RespondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	//token := time.Now().Format("2006-01-02 15:04:05.000000")

	security := insertarTokenSecurity(tokenDecode, w)

	err = Actualizar(security)
	if err != nil {
		framework.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	framework.RespondJSON(w, http.StatusOK, security)

}
