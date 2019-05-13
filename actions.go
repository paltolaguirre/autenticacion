package main

import (
	"encoding/base64"
	"math/rand"
	"net/http"
	"strconv"
	s "strings"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/xubiosueldos/autenticacion/publico"
	"github.com/xubiosueldos/framework"
)

//var db *gorm.DB
var err error
var m = make(map[string]publico.TokenAutenticacion)
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
		infoUser := s.Split(string(tokenDecode), ":")

		username := infoUser[0]
		pass := infoUser[1]
		tenant := infoUser[2]

		numeroRandom := rand.Int63()
		token := strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.FormatInt(numeroRandom, 10)

		fecha := time.Now()

		autenticacion := publico.TokenAutenticacion{Username: username, Pass: pass, Tenant: tenant, Token: token, FechaCreacion: fecha}

		m[string(token)] = autenticacion

		framework.RespondJSON(w, http.StatusOK, autenticacion)

	} else {
		framework.RespondError(w, http.StatusNotFound, framework.DatosIncorrectos)
		return
	}

}

func CheckToken(w http.ResponseWriter, r *http.Request) {

	header := r.Header.Get("Authorization")

	token := s.Split(header, " ")[1]

	tokenAutenticacion, ok := m[token]

	if ok {
		framework.RespondJSON(w, http.StatusOK, tokenAutenticacion)

	} else {
		errors = publico.Error{ErrorNombre: "Hubo error", ErrorCodigo: http.StatusUnauthorized}
		framework.RespondError(w, errors.ErrorCodigo, errors.ErrorNombre)
		return
	}
}

func chequeoAuthenticationMonolitico(tokenEncode string, r *http.Request) bool {

	infoUserValida := false

	url := "http://localhost:8080/NXV/SecurityAuthenticationGo"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", tokenEncode)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	if res.StatusCode == http.StatusAccepted {
		infoUserValida = true
	}

	return infoUserValida
}
