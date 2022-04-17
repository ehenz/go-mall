package model

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityID uint // 1-有管理员权限
	jwt.StandardClaims
}
