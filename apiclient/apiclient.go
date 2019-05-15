package apiclient

import (
	"net/http"

	s "strings"

	"github.com/jinzhu/gorm"
	"github.com/xubiosueldos/autenticacion/publico"
	"github.com/xubiosueldos/conexionBD"
	"github.com/xubiosueldos/framework"
)

func checkTokenDB(w http.ResponseWriter, token string) (*publico.Security, bool, error) {

	var existeToken bool = true
	var security publico.Security
	var err error = nil
	db := conexionBD.ConnectBD("security")

	if err = db.Set("gorm:auto_preload", true).First(&security, "token = ?", token).Error; gorm.IsRecordNotFoundError(err) {
		framework.RespondError(w, http.StatusNotFound, err.Error())
		existeToken = false
	}

	return &security, existeToken, err
}

func obtenerTokenHeader(r *http.Request) string {

	header := r.Header.Get("Authorization")

	token := s.Split(header, " ")[1]

	return token

}
