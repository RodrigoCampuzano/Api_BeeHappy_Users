package entities

const (
	VerificationTypeLogin          = "login"
	VerificationTypePasswordReset  = "reset"
	VerificationTypePasswordChange = "change"
)

type User struct {
	ID                   int    `json:"id"`
	Usuario              string `json:"usuario"`
	Contrasena           string `json:"contrasena"`
	Nombres              string `json:"nombres"`
	Apellidos            string `json:"apellidos"`
	Correo_electronico   string `json:"correo_electronico"`
	Rol                  string `json:"rol"`
	Estado               string `json:"estado"`
	Fecha_registro       string `json:"fecha_registro"`
	VerificacionDospasos bool   `json:"verificacion_dos_pasos"`
}

type VerificationCode struct {
	ID                int    `json:"id"`
	CorreoElectronico string `json:"correo_electronico"`
	Codigo            string `json:"codigo"`
	FechaExpiracion   string `json:"fecha_expiracion"`
	Tipo              string `json:"tipo"`
}
