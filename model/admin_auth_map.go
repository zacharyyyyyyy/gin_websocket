package model

type AdminAuthMap struct {
	Role       int
	Auth       int
	CreateTime int
}

type AdminAuthMapDetail struct {
	Role         int
	Auth         int
	RoleName     string
	RoleDescribe string
	AuthName     string
}
