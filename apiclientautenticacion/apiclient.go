package apiclientautenticacion

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	s "strings"

	"github.com/xubiosueldos/framework/configuracion"

	"github.com/xubiosueldos/autenticacion/publico"
	"github.com/xubiosueldos/framework"
)

func CheckTokenValidoConMicroservicioAutenticacion(r *http.Request) (*publico.Security, *publico.Error) {

	config := configuracion.GetInstance()

	var tokenAutenticacion *publico.Security
	var tokenError *publico.Error
	puerto := config.Puertomicroservicio
	if puerto == "" {
		puerto = config.Puertomicroserivicioautenticacion
	}
	url := configuracion.GetUrlMicroservicio(puerto) + "/auth/check-token"

	req, _ := http.NewRequest("GET", url, nil)

	header := r.Header.Get("Authorization")

	req.Header.Add("Authorization", header)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != http.StatusUnauthorized {

		// tokenAutenticacion = &(TokenAutenticacion{})
		tokenAutenticacion = new(publico.Security)
		json.Unmarshal([]byte(string(body)), tokenAutenticacion)

	} else {
		tokenError = new(publico.Error)
		errorrespuesta := s.Split(res.Status, " ")
		tokenError.ErrorNombre = errorrespuesta[1]
		tokenError.ErrorCodigo = res.StatusCode
		json.Unmarshal([]byte(string(body)), tokenError)

	}

	return tokenAutenticacion, tokenError
}

func ErrorToken(w http.ResponseWriter, tokenError *publico.Error) {
	errorToken := *tokenError
	framework.RespondError(w, errorToken.ErrorCodigo, errorToken.ErrorNombre)

}

func CheckTokenValido(w http.ResponseWriter, r *http.Request) (bool, *publico.Security) {
	var tokenValido bool = true
	tokenAutenticacion, tokenError := CheckTokenValidoConMicroservicioAutenticacion(r)

	if tokenError != nil {
		ErrorToken(w, tokenError)
		tokenValido = false
	}

	return tokenValido, tokenAutenticacion
}
