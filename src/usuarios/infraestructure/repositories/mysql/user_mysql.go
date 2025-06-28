package mysql

import (
	core "apiusuarios/src/core/db"
	"apiusuarios/src/usuarios/domain/entities"
	"fmt"
	"log"
)

type MySql struct {
	conn *core.Conn_MySQL
}

func NewMySql() *MySql {
	conn := core.GetDBPool()
	if conn.Err != "" {
		log.Fatal("Error connecting to MySQL database: %v", conn.Err)
	}
	return &MySql{conn: conn}
}

func (mysql *MySql) CreateUser(data *entities.User) error {
	query := `INSERT INTO usuarios (usuario, contrasena, nombres, apellidos, correo_electronico, rol, estado, fecha_registro) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := mysql.conn.DB.Query(query, data.Usuario, data.Contrasena, data.Nombres, data.Apellidos, data.Correo_electronico, data.Rol, data.Estado, data.Fecha_registro)
	if err != nil {
		log.Println("Error insertando los valores:", err)
		return err
	}
	fmt.Println("Nuevo dato almacenado con ID:")
	return nil
}
