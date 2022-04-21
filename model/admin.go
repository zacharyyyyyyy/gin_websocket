package model

type Admin struct {
	Id         int
	Username   string
	Password   string
	Name       string
	Role       int
	CreateTime int
}

type AdminWithRole struct {
	Id         int
	Username   string
	Name       string
	RoleName   string
	Describe   string
	CreateTime int
}
