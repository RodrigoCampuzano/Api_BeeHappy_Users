package mysql

import (
	"apiusuarios/src/core/auth"
	core "apiusuarios/src/core/db"
	"apiusuarios/src/usuarios/domain/entities"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type MySql struct {
	conn *core.Conn_MySQL
}

func NewMySql() *MySql {
	conn := core.GetDBPool()
	if conn.Err != "" {
		log.Fatalf("Error connecting to MySQL database: %v", conn.Err)
	}
	return &MySql{conn: conn}
}

func (mysql *MySql) CreateUser(data *entities.User) error {
    godotenv.Load()
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        log.Println("JWT_SECRET no definido")
        return fmt.Errorf("JWT_SECRET no definido")
    }

	 passwordHasheada := auth.HashPasswordWithSecret(data.Contrasena, secret)
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordHasheada), bcrypt.DefaultCost)
    if err != nil {
        log.Println("Error encriptando la contraseña:", err)
        return err
    }
    data.Contrasena = string(hashedPassword)

	// Establecer valores por defecto
	if data.Estado == "" {
		data.Estado = "activo"
	}
	if data.Fecha_registro == "" {
		data.Fecha_registro = time.Now().Format("2006-01-02 15:04:05")
	}

	query := `INSERT INTO usuarios (usuario, contrasena, nombres, apellidos, correo_electronico, rol, estado, fecha_registro, verificacion_dos_pasos) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = mysql.conn.DB.Query(query,
		data.Usuario,
		data.Contrasena,
		data.Nombres,
		data.Apellidos,
		data.Correo_electronico,
		data.Rol,
		data.Estado,
		data.Fecha_registro,
		false) // verificacion_dos_pasos por defecto es false
    if err != nil {
        log.Println("Error insertando los valores:", err)
        return err
    }
	fmt.Println("Nuevo usuario creado exitosamente")
    return nil
}

func (mysql *MySql) GetUserByUsuario(usuario string) (*entities.User, error) {
	query := `SELECT id, usuario, contrasena, nombres, apellidos, correo_electronico, rol, estado, fecha_registro, verificacion_dos_pasos 
			  FROM usuarios WHERE usuario = ?`
	row := mysql.conn.DB.QueryRow(query, usuario)

	var user entities.User
	err := row.Scan(
		&user.ID,
		&user.Usuario,
		&user.Contrasena,
		&user.Nombres,
		&user.Apellidos,
		&user.Correo_electronico,
		&user.Rol,
		&user.Estado,
		&user.Fecha_registro,
		&user.VerificacionDospasos,
	)
	if err != nil {
		log.Println("Error obteniendo el usuario:", err)
		return nil, err
	}
	return &user, nil
}

func (mysql *MySql) GetUserByEmail(email string) (*entities.User, error) {
	query := `SELECT id, usuario, contrasena, nombres, apellidos, correo_electronico, rol, estado, fecha_registro, verificacion_dos_pasos 
			  FROM usuarios WHERE correo_electronico = ?`
	row := mysql.conn.DB.QueryRow(query, email)

	var user entities.User
	err := row.Scan(
		&user.ID,
		&user.Usuario,
		&user.Contrasena,
		&user.Nombres,
		&user.Apellidos,
		&user.Correo_electronico,
		&user.Rol,
		&user.Estado,
		&user.Fecha_registro,
		&user.VerificacionDospasos,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (mysql *MySql) UpdatePassword(email string, newPassword string) error {
	godotenv.Load()
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return fmt.Errorf("JWT_SECRET no definido")
	}

	passwordHasheada := auth.HashPasswordWithSecret(newPassword, secret)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordHasheada), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `UPDATE usuarios SET contrasena = ? WHERE correo_electronico = ?`
	_, err = mysql.conn.DB.Exec(query, string(hashedPassword), email)
	return err
}

func (mysql *MySql) UpdateVerificacionDospasos(userID int, estado bool) error {
	query := `UPDATE usuarios SET verificacion_dos_pasos = ? WHERE id = ?`
	_, err := mysql.conn.DB.Exec(query, estado, userID)
	return err
}

func (mysql *MySql) SaveVerificationCode(code *entities.VerificationCode) error {
	// Primero eliminamos cualquier código anterior del mismo tipo para el mismo correo
	deleteQuery := `DELETE FROM verification_codes WHERE correo_electronico = ? AND tipo = ?`
	_, err := mysql.conn.DB.Exec(deleteQuery, code.CorreoElectronico, code.Tipo)
	if err != nil {
		return err
	}

	// Insertamos el nuevo código
	query := `INSERT INTO verification_codes (correo_electronico, codigo, fecha_expiracion, tipo) 
			  VALUES (?, ?, ?, ?)`
	_, err = mysql.conn.DB.Exec(query, code.CorreoElectronico, code.Codigo, code.FechaExpiracion, code.Tipo)
	return err
}

func (mysql *MySql) GetVerificationCode(email string, tipo string) (*entities.VerificationCode, error) {
	query := `SELECT id, correo_electronico, codigo, fecha_expiracion, tipo 
			  FROM verification_codes 
			  WHERE correo_electronico = ? AND tipo = ?`
	row := mysql.conn.DB.QueryRow(query, email, tipo)

	var code entities.VerificationCode
	err := row.Scan(&code.ID, &code.CorreoElectronico, &code.Codigo, &code.FechaExpiracion, &code.Tipo)
	if err != nil {
		return nil, err
	}

	// Verificar si el código ha expirado
	expTime, err := time.Parse(time.RFC3339, code.FechaExpiracion)
	if err != nil {
		return nil, err
	}

	if time.Now().After(expTime) {
		// Si el código ha expirado, lo eliminamos y retornamos error
		mysql.DeleteVerificationCode(code.ID)
		return nil, fmt.Errorf("el código ha expirado")
	}

	return &code, nil
}

func (mysql *MySql) DeleteVerificationCode(id int) error {
	query := `DELETE FROM verification_codes WHERE id = ?`
	_, err := mysql.conn.DB.Exec(query, id)
	return err
}
