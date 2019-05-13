package publico

import "time"

type Security struct {
	//	gorm.Model
	Username      string    `json:"username"`
	Pass          string    `json:"pass"`
	Tenant        string    `json:"tenant"`
	Token         string    `json:"token"`
	FechaCreacion time.Time `json:"fechacreacion"`
}
