package types

// AuthBody 认证请求体
type AuthBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
