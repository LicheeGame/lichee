package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTInstance struct {
	SecretKey []byte
}

func InitJwt(SecretKey []byte) JWTInstance {
	return JWTInstance{SecretKey}
}

type CustomClaims struct {
	OpenID string `json:"openid"`
	jwt.RegisteredClaims
}

// GenerateJWT 生成JWT
func (this JWTInstance) GenerateJWT(openid string, count time.Duration) string {
	if count == 0 {
		count = 2
	}
	// 设置一些声明
	claims := CustomClaims{
		openid,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(count * time.Hour)), //有效时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                        //签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                        //生效时间
			Issuer:    "test",                                                //签发人
			Subject:   "somebody",                                            //主题
			ID:        "1",                                                   //JWT ID用于标识该JWT
		},
	}

	//使用指定的加密方式和声明类型创建新令牌
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 设置签名并获取token字符串
	token, err := jwtToken.SignedString(this.SecretKey)
	if err != nil {
		return ""
	}

	return token
}

// ParseJWT 解析JWT
func (this JWTInstance) ParseJWT(tokenString string) *CustomClaims {
	// 解析JWT字符串
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return this.SecretKey, nil
	})

	if err != nil {
		return nil
	}

	// 验证token
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims
	}

	return nil
}
