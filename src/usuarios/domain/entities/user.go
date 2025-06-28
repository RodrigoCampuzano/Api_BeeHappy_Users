package entities

type User struct {
	Usuario  string  `json:"usuario"`
	Contrasena string  `json:"contrasena"`
	Nombres  string  `json:"nombres"`
	Apellidos string  `json:"apellidos"`
	Correo_electronico string  `json:"correo_electronico"`
	Rol	   string  `json:"rol"`
	Estado   string  `json:"estado"`
	Fecha_registro string  `json:"fecha_registro"`
}