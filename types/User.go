package types

// UserAuthBody 认证请求体
type UserAuthBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserUpdateProfile
type UserUpdateProfile struct {
    NickName *string `json:"nickname"`
    Birth    *uint64 `json:"birth"`
    Gender   *string `json:"gender"`
}
