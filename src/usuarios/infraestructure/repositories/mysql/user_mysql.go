// src/usuarios/infraestructure/repositories/mysql/user_mysql.go (ACTUALIZADO)
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

    query := `INSERT INTO usuarios (usuario, contrasena, nombres, apellidos, correo_electronico, rol, estado, fecha_registro) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
    _, err = mysql.conn.DB.Query(query, data.Usuario, data.Contrasena, data.Nombres, data.Apellidos, data.Correo_electronico, data.Rol, data.Estado, data.Fecha_registro)
    if err != nil {
        log.Println("Error insertando los valores:", err)
        return err
    }
    fmt.Println("Nuevo dato almacenado con ID:")
    return nil
}

func (mysql *MySql) GetUserByUsuario(usuario string) (*entities.User, error) {
	query := `SELECT usuario, contrasena, nombres, apellidos, correo_electronico, rol, estado, fecha_registro FROM usuarios WHERE usuario = ?`
	row := mysql.conn.DB.QueryRow(query, usuario)

	var user entities.User
	err := row.Scan(&user.Usuario, &user.Contrasena, &user.Nombres, &user.Apellidos, &user.Correo_electronico, &user.Rol, &user.Estado, &user.Fecha_registro)
	if err != nil {
		log.Println("Error obteniendo el usuario:", err)
		return nil, err
	}
	return &user, nil
}

// NUEVO: Obtener usuario por email
func (mysql *MySql) GetUserByEmail(email string) (*entities.User, error) {
	query := `SELECT usuario, contrasena, nombres, apellidos, correo_electronico, rol, estado, fecha_registro FROM usuarios WHERE correo_electronico = ?`
	row := mysql.conn.DB.QueryRow(query, email)

	var user entities.User
	err := row.Scan(&user.Usuario, &user.Contrasena, &user.Nombres, &user.Apellidos, &user.Correo_electronico, &user.Rol, &user.Estado, &user.Fecha_registro)
	if err != nil {
		log.Println("Error obteniendo el usuario por email:", err)
		return nil, err
	}
	return &user, nil
}

// NUEVO: Actualizar contraseña
func (mysql *MySql) UpdatePassword(usuario, newPassword string) error {
	query := `UPDATE usuarios SET contrasena = ? WHERE usuario = ?`
	_, err := mysql.conn.ExecutePreparedQuery(query, newPassword, usuario)
	if err != nil {
		log.Println("Error actualizando la contraseña:", err)
		return err
	}
	return nil
}