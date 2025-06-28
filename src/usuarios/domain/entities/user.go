package entities

// User representa la estructura de un usuario en el sistema
type User struct {
	Usuario            string `json:"usuario" db:"usuario"`
	Contrasena         string `json:"contrasena" db:"contrasena"`
	Nombres            string `json:"nombres" db:"nombres"`
	Apellidos          string `json:"apellidos" db:"apellidos"`
	Correo_electronico string `json:"correo_electronico" db:"correo_electronico"`
	Rol                string `json:"rol" db:"rol"`
	Estado             string `json:"estado" db:"estado"`
	Fecha_registro     string `json:"fecha_registro" db:"fecha_registro"`
	TwoFactorSecret    string `json:"-" db:"two_factor_secret"` // No se expone en JSON por seguridad
	TwoFactorEnabled   bool   `json:"two_factor_enabled" db:"two_factor_enabled"`
}

// LoginRequest representa la solicitud de login
type LoginRequest struct {
	Usuario     string `json:"usuario" binding:"required" example:"usuario123"`
	Contrasena  string `json:"contrasena" binding:"required" example:"password123"`
	TokenTOTP   string `json:"token_totp,omitempty" example:"123456"`
}

// LoginResponse representa la respuesta del login
type LoginResponse struct {
	Token   string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	Message string `json:"message,omitempty" example:"Login exitoso"`
}

// Setup2FARequest representa la solicitud para configurar 2FA
type Setup2FARequest struct {
	TokenTOTP string `json:"token_totp" binding:"required" example:"123456"`
}

// Setup2FAResponse representa la respuesta para configurar 2FA
type Setup2FAResponse struct {
	QRCode    string `json:"qr_code" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."`
	Secret    string `json:"secret" example:"JBSWY3DPEHPK3PXP"`
	Message   string `json:"message" example:"2FA configurado exitosamente"`
}

// ErrorResponse representa una respuesta de error
type ErrorResponse struct {
	Error string `json:"error" example:"Error message"`
}