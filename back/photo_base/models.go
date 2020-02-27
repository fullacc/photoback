package photo_base

type ConfigFile struct{
	Port string `json:"port"`
	Host string `json:"host"`
	DbHost string `json:"dbhost"`
	DbPort string `json:"dbport"`
	Password string `json:"password"`
	User string `json:"user"`
	Name string `json:"name"`
}

type User struct {
	Id int `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role int `json:"role"`
}