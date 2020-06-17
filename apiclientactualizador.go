package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xubiosueldos/conexionBD/Autenticacion/structAutenticacion"
	"github.com/xubiosueldos/framework/configuracion"
	"io/ioutil"
	"net/http"
)

func Actualizar(security *structAutenticacion.Security) error {
	config := configuracion.GetInstance()
	url := configuracion.GetUrlMicroservicio(config.Puertomicroservicioactualizacion) + "updater/update"

	payload, err := json.Marshal(*security)

	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))

	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	fmt.Println("URL:", url)

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	if res.StatusCode != http.StatusOK {
		var respuestaDeError = make(map[string] string)
		json.Unmarshal(body, &respuestaDeError)
		return errors.New(respuestaDeError["error"])
	}

	return nil
}
