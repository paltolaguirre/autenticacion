package apiclientautenticacion

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	s "strings"

	"github.com/xubiosueldos/autenticacion/publico"
	"github.com/xubiosueldos/framework"
)

func CheckTokenValido(r *http.Request) (*publico.Security, *publico.Error) {

	var tokenAutenticacion *publico.Security
	var tokenError *publico.Error

	url := "http://localhost:8081/check-token"

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
