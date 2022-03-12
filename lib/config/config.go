package config

import "github.com/go-ini/ini"

type ConfMap interface {
	MapTo()
}

type ConfHandle struct {
	cfg *ini.File
}
