package bean

type ProfileToken struct {
	ProfileId    string `json:"profile_id" orm:"profileId"`
	Region       string `json:"region" orm:"region"`
	RefreshToken string `json:"refresh_token" orm:"refreshToken"`
	ClientId     string `json:"clientId" orm:"clientId" `
	ClientSecret string `json:"clientSecret" orm:"clientSecret"`
	Tag          int    `json:"tag" orm:"tag"`
	ProfileType  string `json:"profileType" orm:"profileType"`
}

// ProfileTokenClient 带client的客户数据
type ProfileTokenClient struct {
	ProfileId    int64  `json:"profile_id"`
	Region       string `json:"region"`
	RefreshToken string `json:"refresh_token"`
	ClientTag    string `json:"client_tag"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type ProfileAll struct {
	ProfileId int64  `json:"profile_id"`
	Region    string `json:"region"`
	NickName  string `json:"nickname"`
}
