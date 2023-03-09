package proto

type LarkEventCallbackRequest struct {
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
	Type      string `json:"type"`
}

type LarkEventCallbackResponse struct {
	Challenge string `json:"challenge"`
}

// JoinGroupEvent im.chat.member.user.added_v1
type JoinGroupEvent struct {
	ChatId     string `json:"chat_id"`
	OperatorId struct {
		UnionId string `json:"union_id"`
		UserId  string `json:"user_id"`
		OpenId  string `json:"open_id"`
	} `json:"operator_id"`
	External          bool   `json:"external"`
	OperatorTenantKey string `json:"operator_tenant_key"`
	Users             []struct {
		Name      string `json:"name"`
		TenantKey string `json:"tenant_key"`
		UserId    struct {
			UnionId string `json:"union_id"`
			UserId  string `json:"user_id"`
			OpenId  string `json:"open_id"`
		} `json:"user_id"`
	} `json:"users"`
	Name      string `json:"name"`
	I18NNames struct {
		ZhCn string `json:"zh_cn"`
		EnUs string `json:"en_us"`
		JaJp string `json:"ja_jp"`
	} `json:"i18n_names"`
}
