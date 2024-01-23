package global

import (
	"mxshop-api/user-web/config"
	"strings"

	ut "github.com/go-playground/universal-translator"
)

var Translator ut.Translator

var ServerConfig *config.ServerConfig = &config.ServerConfig{}

type JWTInfo struct {
	SigningKey string
}

func RemoveTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}
