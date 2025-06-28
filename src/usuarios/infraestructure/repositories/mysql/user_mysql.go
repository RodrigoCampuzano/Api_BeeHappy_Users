package mysql

import (
	"apiusuarios/src/core/auth"
	core "apiusuarios/src/core/db"
	"apiusuarios/src/usuarios/domain/entities"
	"fmt"
	"log"
	"os"

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

	query := `INSERT INTO usuarios (usuario, contrasena, nombres, apellidos, correo_electronico, rol, estado, fecha_registro, two_factor_enabled) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = mysql.conn.DB.Query(query, data.Usuario, data.Contrasena, data.Nombres, data.Apellidos, data.Correo_electronico, data.Rol, data.Estado, data.Fecha_registro, false)
	if err != nil {
		log.Println("Error insertando los valores:", err)
		return err
	}
	fmt.Println("Nuevo usuario creado exitosamente")
	return nil
}

func (mysql *MySql) GetUserByUsuario(usuario string) (*entities.User, error) {
	query := `SELECT usuario, contrasena, nombres, apellidos, correo_electronico, rol, estado, fecha_registro, 
			  COALESCE(two_factor_secret, '') as two_factor_secret, 
			  COALESCE(two_factor_enabled, false) as two_factor_enabled 
			  FROM usuarios WHERE usuario = ?`
	row := mysql.conn.DB.QueryRow(query, usuario)

	var user entities.User
	err := row.Scan(&user.Usuario, &user.Contrasena, &user.Nombres, &user.Apellidos, 
		&user.Correo_electronico, &user.Rol, &user.Estado, &user.Fecha_registro,
		&user.TwoFactorSecret, &user.TwoFactorEnabled)
	if err != nil {
		log.Println("Error obteniendo el usuario:", err)
		return nil, err
	}
	return &user, nil
}

func (mysql *MySql) UpdateTwoFactorSecret(usuario, secret string) error {
	query := `UPDATE usuarios SET two_factor_secret = ? WHERE usuario = ?`
	_, err := mysql.conn.ExecutePreparedQuery(query, secret, usuario)
	if err != nil {
		log.Println("Error actualizando secreto 2FA:", err)
		return err
	}
	return nil
}

func (mysql *MySql) EnableTwoFactor(usuario string) error {
	query := `UPDATE usuarios SET two_factor_enabled = true WHERE usuario = ?`
	_, err := mysql.conn.ExecutePreparedQuery(query, usuario)
	if err != nil {
		log.Println("Error habilitando 2FA:", err)
		return err
	}
	return nil
}

func (mysql *MySql) DisableTwoFactor(usuario string) error {
	query := `UPDATE usuarios SET two_factor_enabled = false, two_factor_secret = NULL WHERE usuario = ?`
	_, err := mysql.conn.ExecutePreparedQuery(query, usuario)
	if err != nil {
		log.Println("Error deshabilitando 2FA:", err)
		return err
	}
	return nil
}