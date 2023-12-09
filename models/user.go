package models

import (
	"github.com/VDiPaola/go-backend-module/module_models"
)

type User struct {
	Id             uint               `gorm:"primary_key" json:"id"`
	Email          string             `gorm:"unique" json:"email"`
	Password       string             `json:"-"`
	AuthorName     string             `json:"author_name"`
	Role           RoleType           `gorm:"default:user" json:"role"`
	Verified       bool               `gorm:"default:false" json:"verified"`
	HasPassword    bool               `gorm:"default:false" json:"has_password"`
	AccountCreated int64              `gorm:"autoCreateTime" json:"account_created"`
	Code           module_models.Code `gorm:"embedded;embeddedPrefix:code_" json:"-"`
	LastLoginUnix  int64              `json:"last_login_unix"`
}

type RoleType string

var Role = struct {
	User   RoleType
	Member RoleType
}{
	User:   "user",
	Member: "member",
}
