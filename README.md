# API de Usuarios - Documentación

Esta API proporciona servicios de gestión de usuarios con características avanzadas de seguridad, incluyendo autenticación en dos pasos y recuperación de contraseña.

## Índice
- [Configuración](#configuración)
- [Autenticación](#autenticación)
- [Endpoints](#endpoints)
  - [Gestión de Usuarios](#gestión-de-usuarios)
  - [Autenticación](#endpoints-de-autenticación)
  - [Gestión de Contraseñas](#gestión-de-contraseñas)
  - [Seguridad](#seguridad)

## Configuración

### Variables de Entorno Requeridas
```env
EMAIL_FROM=tu_correo@gmail.com
EMAIL_PASSWORD=tu_app_password
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
JWT_SECRET=tu_secreto_jwt
```

### Ejecución
```bash
go run main.go
```

La API estará disponible en `http://localhost:8080`
Documentación Swagger: `http://localhost:8080/swagger/index.html`

## Autenticación

La API utiliza tokens JWT para la autenticación. Para los endpoints protegidos, incluye el token en el header:
```
Authorization: Bearer <tu_token>
```

### Proceso de Login

El proceso de login varía dependiendo de si el usuario tiene activada la verificación en dos pasos:

1. **Usuario sin 2FA activado**:
   ```http
   POST /api/v1/users/login
   Content-Type: application/json

   {
       "usuario": "rodrigo.martinez",
       "contrasena": "Abc123*"
   }
   ```
   **Respuesta (200 OK)**:
   ```json
   {
       "token": "eyJhbGciOiJIUzI1NiIsInR5cCI...",
       "require_two_factor": false,
       "message": "Inicio de sesión exitoso"
   }
   ```

2. **Usuario con 2FA activado**:
   
   **Paso 1: Iniciar sesión**
   ```http
   POST /api/v1/users/login
   Content-Type: application/json

   {
       "usuario": "rodrigo.martinez",
       "contrasena": "Abc123*"
   }
   ```
   **Respuesta (200 OK)**:
   ```json
   {
       "require_two_factor": true,
       "email": "rodrigo.martinez@empresa.com",
       "message": "Se ha enviado un código de verificación a su correo electrónico"
   }
   ```

   **Paso 2: Verificar código**
   ```http
   POST /api/v1/users/login/verify
   Content-Type: application/json

   {
       "email": "rodrigo.martinez@empresa.com",
       "code": "847291"
   }
   ```
   **Respuesta (200 OK)**:
   ```json
   {
       "token": "eyJhbGciOiJIUzI1NiIsInR5cCI...",
       "require_two_factor": false,
       "message": "Código verificado correctamente. Sesión iniciada."
   }
   ```

## Endpoints

### Gestión de Usuarios

#### Crear Usuario
```http
POST /api/v1/users/
```
Crea un nuevo usuario en el sistema.

**Request:**
```json
{
    "usuario": "rodrigo.martinez",
    "contrasena": "Abc123*",
    "nombres": "Rodrigo Antonio",
    "apellidos": "Martínez García",
    "correo_electronico": "rodrigo.martinez@empresa.com",
    "rol": "apicultor"
}
```

**Respuesta Exitosa (200):**
```json
{
    "message": "Usuario creado exitosamente. Se ha enviado un correo de bienvenida."
}
```

### Gestión de Contraseñas

#### Solicitar Cambio de Contraseña (Usuario Autenticado)
```http
POST /api/v1/users/profile/password/change/request
Authorization: Bearer <tu_token>
```
Envía un código de verificación para cambiar la contraseña.

**Request:**
```json
{
    "current_password": "ContraseñaActual123"
}
```

**Respuesta Exitosa (200):**
```json
{
    "message": "Código de verificación enviado"
}
```

#### Confirmar Cambio de Contraseña (Usuario Autenticado)
```http
POST /api/v1/users/profile/password/change
Authorization: Bearer <tu_token>
```
Cambia la contraseña usando el código de verificación.

**Request:**
```json
{
    "email": "tu.correo@empresa.com",
    "code": "847291",
    "new_password": "NuevaContraseña123*"
}
```

**Respuesta Exitosa (200):**
```json
{
    "message": "Contraseña actualizada exitosamente"
}
```

#### Solicitar Restablecimiento de Contraseña (Usuario No Autenticado)
```http
POST /api/v1/users/password/reset/request
```
Envía un código de verificación para restablecer la contraseña.

**Request:**
```json
{
    "email": "rodrigo.martinez@empresa.com"
}
```

**Respuesta Exitosa (200):**
```json
{
    "message": "Se ha enviado un código de verificación a su correo electrónico"
}
```

#### Restablecer Contraseña (Usuario No Autenticado)
```http
POST /api/v1/users/password/reset
```
Restablece la contraseña usando el código de verificación.

**Request:**
```json
{
    "email": "rodrigo.martinez@empresa.com",
    "code": "847291",
    "new_password": "NuevaContraseña123*"
}
```

**Respuesta Exitosa (200):**
```json
{
    "message": "Contraseña restablecida exitosamente"
}
```

### Seguridad

#### Activar/Desactivar Verificación en Dos Pasos
```http
POST /api/v1/users/profile/2fa/toggle
```
Activa o desactiva la verificación en dos pasos (requiere autenticación).

**Request:**
```json
{
    "estado": true
}
```

**Respuesta Exitosa (200):**
```json
{
    "message": "La verificación en dos pasos ha sido activada exitosamente"
}
```

## Códigos de Error

### 400 Bad Request
```json
{
    "error": "Datos inválidos. Por favor, verifique la información proporcionada"
}
```

### 401 Unauthorized
```json
{
    "error": "Credenciales inválidas. Por favor, verifique su usuario y contraseña"
}
```
O
```json
{
    "error": "Token no proporcionado o inválido"
}
```

### 403 Forbidden
```json
{
    "error": "No tiene permisos para realizar esta acción"
}
```

## Notas de Seguridad

1. **Contraseñas**
   - Mínimo 8 caracteres
   - Debe incluir mayúsculas, minúsculas y números
   - Se recomienda incluir caracteres especiales

2. **Tokens JWT**
   - Expiran después de 72 horas
   - Incluyen información del rol y usuario
   - No almacenan información sensible

3. **Verificación en Dos Pasos**
   - Códigos de 6 dígitos
   - Expiran después de 5 minutos
   - Se envían por correo electrónico
   - Máximo 3 intentos por código

4. **Limitación de Intentos**
   - 5 intentos fallidos de login bloquean la cuenta por 15 minutos
   - 3 intentos fallidos de código de verificación invalidan el código
   - Se requiere solicitar un nuevo código después de 3 intentos fallidos 