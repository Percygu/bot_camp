package proto

type TokenResp struct {
	Code int            `json:"code"`
	Msg  string         `json:"msg"`
	Data *TokenRespInfo `json:"data"`
}
type TokenRespInfo struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	Name             string `json:"name"`
	EnName           string `json:"en_name"`
	AvatarUrl        string `json:"avatar_url"`
	AvatarThumb      string `json:"avatar_thumb"`
	AvatarMiddle     string `json:"avatar_middle"`
	AvatarBig        string `json:"avatar_big"`
	OpenId           string `json:"open_id"`
	UnionId          string `json:"union_id"`
	Email            string `json:"email"`
	EnterpriseEmail  string `json:"enterprise_email"`
	UserId           string `json:"user_id"`
	Mobile           string `json:"mobile"`
	TenantKey        string `json:"tenant_key"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
}
