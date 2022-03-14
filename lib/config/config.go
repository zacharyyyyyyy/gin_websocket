package config

import "github.com/go-ini/ini"

type ConfMap interface {
	MapTo()
	getPath() string
}

type ConfHandle struct {
	cfg *ini.File
}

func init() {

}

func (ConfHandle *ConfHandle) Load() {

}
