package apiclientautenticacion

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	s "strings"

	"github.com/xubio-inc/sueldos-lib-framework/configuracion"

	"github.com/xubiosueldos/conexionBD/Autenticacion/structAutenticacion"
	"github.com/xubio-inc/sueldos-lib-framework"
)

func CheckTokenValidoConMicroservicioAutenticacion(r *http.Request) (*structAutenticacion.Security, *structAutenticacion.Error) {

	var tokenAutenticacion *structAutenticacion.Security
	var tokenError *structAutenticacion.Error
	config := configuracion.GetInstance()
	url := configuracion.GetUrlMicroservicio(config.Puertomicroservicioautenticacion) + "auth/check-token"

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("Error: ", err)
	}
	header := r.Header.Get("Authorization")

	req.Header.Add("Authorization", header)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("URL:", url)

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	if res.StatusCode != http.StatusUnauthorized {

		// tokenAutenticacion = &(TokenAutenticacion{})
		tokenAutenticacion = new(structAutenticacion.Security)
		json.Unmarshal([]byte(string(body)), tokenAutenticacion)

	} else {
		tokenError = new(structAutenticacion.Error)
		errorrespuesta := s.Split(res.Status, " ")
		tokenError.ErrorNombre = errorrespuesta[1]
		tokenError.ErrorCodigo = res.StatusCode
		json.Unmarshal([]byte(string(body)), tokenError)

	}

	return tokenAutenticacion, tokenError
}

func ErrorToken(w http.ResponseWriter, tokenError *structAutenticacion.Error) {
	errorToken := *tokenError
	framework.RespondError(w, errorToken.ErrorCodigo, errorToken.ErrorNombre)

}

func CheckTokenValido(w http.ResponseWriter, r *http.Request) (bool, *structAutenticacion.Security) {
	var tokenValido bool = true
	tokenAutenticacion, tokenError := CheckTokenValidoConMicroservicioAutenticacion(r)

	if tokenError != nil {
		ErrorToken(w, tokenError)
		tokenValido = false
	}

	return tokenValido, tokenAutenticacion
}

func ObtenerTenant(tokenAutenticacion *structAutenticacion.Security) string {

	token := *tokenAutenticacion
	tenant := token.Tenant
	return tenant
}
