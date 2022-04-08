package model

type AdminAuth struct {
	Id         int
	Name       string
	Method     string
	Path       string
	Enable     int
	ParentId   int
	CreateTime int
}
