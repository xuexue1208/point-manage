package utils

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/wonderivan/logger"
)

var JWTToken jwtToken

type jwtToken struct{}

//加解密因子
const (
	SECRET = "xxxxxxxxx"
	//在vue前端的Login.vue中也需要设置一样的秘钥
)

type JWTClaims struct { // token里面添加用户信息，验证token后可能会用到用户信息
	jwt.StandardClaims
	Data string
}

//解析token  兼容xadmin
func (*jwtToken) ParseToken2(token string) (claims string, err error) {
	vToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})
	if err != nil {
		logger.Error("parse token failed ", err)
		//处理token解析后的各种错误
		if ve, ok := err.(*jwt.ValidationError); ok {
			//格式错误
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return "nil", errors.New("TokenMalformed")
				//过期
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return "nil", errors.New("TokenExpired")
				//无效
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return "nil", errors.New("TokenNotValidYet")
			} else {
				return "nil", errors.New("TokenInvalid")
			}
		}
	}
	//讲token对象中的claims断言成*CustomClaims类型并返回
	if claims, ok := vToken.Claims.(*JWTClaims); ok && vToken.Valid {
		//logger.Info("parse token success ", claims)
		return claims.Data, nil
	}
	return "nil", errors.New("解析Token失败")
}

//生产token
func (*jwtToken) GenToken2(mobile string) (string, error) {
	data := map[string]interface{}{
		"mobile": mobile,
	}
	//data转json
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	c := JWTClaims{
		StandardClaims: jwt.StandardClaims{
			Audience:  "",
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			Id:        "",
			IssuedAt:  0,
			Issuer:    "point",
			NotBefore: 0,
			Subject:   "",
		},
		Data: string(dataJson),
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	//必须传入[]byte类型的数据
	return token.SignedString([]byte(SECRET))
}
