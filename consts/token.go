/*
Package consts - NekoBlog backend server constants.
This file is for bearer token related constants.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package consts

// TOKEN_EXPIRE_DURATION 有效期
const TOKEN_EXPIRE_DURATION = 7 * 24 * 60 * 60

// TOKEN_SECRET 令牌密钥
const TOKEN_SECRET = "NEKO_MICRO_BLOG_BACKEND_EXAMPLE_SECRET"

// TOKEN_ISSUER 令牌签发者
const TOKEN_ISSUER = "org.kirisakiii.neko"

// MAX_TOKENS_PER_USER 最大令牌数量
const MAX_TOKENS_PER_USER = 5
