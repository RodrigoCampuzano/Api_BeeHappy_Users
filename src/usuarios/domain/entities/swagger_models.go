package entities

// LoginRequest modelo para la solicitud de inicio de sesión
type SwaggerLoginRequest struct {
	// Usuario o correo electrónico
	Usuario string `json:"usuario" example:"rodrigo.martinez" description:"Usuario o correo electrónico"`
	// Contraseña del usuario
	Contrasena string `json:"contrasena" example:"Abc123*" description:"Contraseña del usuario"`
}

// LoginResponse modelo para la respuesta de inicio de sesión
type SwaggerLoginResponse struct {
	// Token JWT si la autenticación es exitosa
	Token string `json:"token,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTE0Mzc2NDEsInJvbCI6ImFkbWluaXN0cmFkb3IiLCJ1c3VhcmlvIjoiUm9kcmlnbyJ9.-kQtpw4rfS5xFd9ZDvpGIHJw2ppvCgqIedWjCK-JZUs" description:"Token JWT si la autenticación es exitosa"`
	// Indica si se requiere verificación en dos pasos
	RequireTwoFactor bool `json:"require_two_factor" example:"true" description:"Indica si se requiere verificación en dos pasos"`
	// Correo electrónico del usuario (solo si se requiere verificación en dos pasos)
	Email string `json:"email,omitempty" example:"rodrigo.martinez@empresa.com" description:"Correo electrónico del usuario (solo si se requiere verificación en dos pasos)"`
	// Mensaje informativo
	Message string `json:"message" example:"Se ha enviado un código de verificación a su correo electrónico" description:"Mensaje informativo"`
}

// VerifyCodeRequest modelo para la solicitud de verificación de código
type SwaggerVerifyCodeRequest struct {
	// Correo electrónico del usuario
	Email string `json:"email" example:"rodrigo.martinez@empresa.com" description:"Correo electrónico del usuario"`
	// Código de verificación de 6 dígitos
	Code string `json:"code" example:"847291" description:"Código de verificación de 6 dígitos"`
}

// VerifyCodeResponse modelo para la respuesta de verificación de código
type SwaggerVerifyCodeResponse struct {
	// Token JWT si la verificación es exitosa
	Token string `json:"token,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTE0Mzc2NDEsInJvbCI6ImFkbWluaXN0cmFkb3IiLCJ1c3VhcmlvIjoiUm9kcmlnbyJ9.-kQtpw4rfS5xFd9ZDvpGIHJw2ppvCgqIedWjCK-JZUs" description:"Token JWT si la verificación es exitosa"`
	// Mensaje informativo
	Message string `json:"message" example:"Código verificado correctamente. Sesión iniciada." description:"Mensaje informativo"`
}

// SwaggerPasswordResetRequest Modelo para solicitar restablecimiento de contraseña
type SwaggerPasswordResetRequest struct {
	Email string `json:"email" example:"rodrigo.martinez@empresa.com" description:"Correo electrónico del usuario"`
}

// SwaggerPasswordResetResponse Modelo de respuesta para solicitud de restablecimiento
type SwaggerPasswordResetResponse struct {
	Message string `json:"message" example:"Se ha enviado un código de verificación a su correo electrónico" description:"Mensaje informativo"`
}

// SwaggerResetPasswordRequest Modelo para restablecer contraseña
type SwaggerResetPasswordRequest struct {
	Email       string `json:"email" example:"rodrigo.martinez@empresa.com" description:"Correo electrónico del usuario"`
	Code        string `json:"code" example:"847291" description:"Código de verificación"`
	NewPassword string `json:"new_password" example:"NuevaContraseña123*" description:"Nueva contraseña"`
}

// SwaggerResetPasswordResponse Modelo de respuesta para el restablecimiento de contraseña
type SwaggerResetPasswordResponse struct {
	Message string `json:"message" example:"Contraseña actualizada exitosamente" description:"Mensaje informativo"`
}

// ToggleTwoFactorRequest modelo para activar/desactivar la verificación en dos pasos
type SwaggerToggleTwoFactorRequest struct {
	// Estado de la verificación en dos pasos
	Estado bool `json:"estado" example:"true" description:"Estado de la verificación en dos pasos"`
}

// ToggleTwoFactorResponse modelo para la respuesta de activar/desactivar 2FA
type SwaggerToggleTwoFactorResponse struct {
	// Mensaje informativo
	Message string `json:"message" example:"La verificación en dos pasos ha sido activada exitosamente" description:"Mensaje informativo"`
}

// ErrorResponse modelo para respuestas de error
type SwaggerErrorResponse struct {
	// Mensaje de error
	Error string `json:"error" example:"Credenciales inválidas. Por favor, verifique su usuario y contraseña" description:"Mensaje de error"`
}

// CreateUserRequest modelo para la creación de usuario
type SwaggerCreateUserRequest struct {
	// Nombre de usuario (sin espacios, solo letras, números y puntos)
	Usuario string `json:"usuario" example:"rodrigo.martinez" description:"Nombre de usuario (sin espacios, solo letras, números y puntos)"`
	// Contraseña (mínimo 8 caracteres, debe incluir mayúsculas, minúsculas y números)
	Contrasena string `json:"contrasena" example:"Abc123*" description:"Contraseña (mínimo 8 caracteres, debe incluir mayúsculas, minúsculas y números)"`
	// Nombres del usuario
	Nombres string `json:"nombres" example:"Rodrigo Antonio" description:"Nombres del usuario"`
	// Apellidos del usuario
	Apellidos string `json:"apellidos" example:"Martínez García" description:"Apellidos del usuario"`
	// Correo electrónico
	Correo_electronico string `json:"correo_electronico" example:"rodrigo.martinez@empresa.com" description:"Correo electrónico"`
	// Rol del usuario (administrador, apicultor, tecnico)
	Rol string `json:"rol" example:"apicultor" description:"Rol del usuario (administrador, apicultor, tecnico)" enums:"administrador,apicultor,tecnico"`
}

// CreateUserResponse modelo para la respuesta de creación de usuario
type SwaggerCreateUserResponse struct {
	// Mensaje informativo
	Message string `json:"message" example:"Usuario creado exitosamente. Se ha enviado un correo de bienvenida." description:"Mensaje informativo"`
}

// SwaggerChangePasswordRequest Modelo para solicitar cambio de contraseña
type SwaggerChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" example:"ContraseñaActual123" description:"Contraseña actual del usuario"`
}

// SwaggerChangePasswordResponse Modelo de respuesta para cambio de contraseña
type SwaggerChangePasswordResponse struct {
	Message string `json:"message" example:"Código de verificación enviado" description:"Mensaje informativo"`
}
