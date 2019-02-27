package publico


import "time"

type TokenAutenticacion struct {
//	gorm.Model
	Username string
	Pass string
	Tenant string
	Token string
	FechaCreacion time.Time
}