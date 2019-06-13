package publico

import (
	"time"

	"github.com/xubiosueldos/conexionBD/structGormModel"
)

type Security struct {
	structGormModel.GormModel
	Username      string    `json:"username"`
	Pass          string    `json:"pass"`
	Tenant        string    `json:"tenant"`
	Token         string    `json:"token"`
	FechaCreacion time.Time `json:"fechacreacion"`
}
