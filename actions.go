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
	"github.com/xubiosueldos/autenticacion/publico"
	"github.com/xubiosueldos/conexionBD"
	"github.com/xubiosueldos/framework"
	"github.com/xubiosueldos/framework/configuracion"
)

//var db *gorm.DB
var err error
var m = make(map[string]publico.Security)
var errors publico.Error

func Login(w http.ResponseWriter, r *http.Request) {

	header := r.Header.Get("Authorization")
	tokenEncode := s.Split(header, " ")[1]

	//Chequear con el monolitico que los datos ingresados sean correctos
	if chequeoAuthenticationMonolitico(tokenEncode, r) {

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

func CheckToken(w http.ResponseWriter, r *http.Request) {

	header := r.Header.Get("Authorization")

	token := s.Split(header, " ")[1]

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

	configuracion := configuracion.GetInstance()
	infoUserValida := false

	url := configuracion.Url + "SecurityAuthenticationGo"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", tokenEncode)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	if res.StatusCode == http.StatusAccepted {
		infoUserValida = true
	}

	return infoUserValida
}

func insertarTokenSecurity(tokenDecode []byte, w http.ResponseWriter) publico.Security {

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
	}

	return security
}

func checkTokenDB(w http.ResponseWriter, token string) (publico.Security, bool, error) {

	var existeToken bool = true
	var security publico.Security
	var err error = nil
	db := conexionBD.ConnectBD("security")

	if err = db.Set("gorm:auto_preload", true).First(&security, "token = ?", token).Error; gorm.IsRecordNotFoundError(err) {
		framework.RespondError(w, http.StatusNotFound, err.Error())
		existeToken = false
	}

	return security, existeToken, err
}
