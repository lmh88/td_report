package bean

type Cfg struct {
	Redis        *redis        `yaml:"redis"`
	Cos          *cos          `yaml:"cos"`
	DumperClient *dumperClient `yaml:"dumper_client"`
	Extra        *extra `yaml:"extra"`
	Schedule *schedule `yaml:"schedule"`
}

type redis struct {
	Addr         string `yaml:"addr"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	PoolSize     int    `yaml:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns"`
	IdleTimeout  int    `yaml:"idle_timeout"`
}

type cos struct {
	Key               string `yaml:"key"`
	Secret            string `yaml:"secret"`
	BucketName        string `yaml:"bucket_name"`
	NaAshburnEndpoint string `yaml:"na-ashburn_endpoint"`
	HkEndpoint        string `yaml:"hk_endpoint"`
}

type dumperClient struct {
	EndpointList         []string `yaml:"endpoint_list"`
	FileNames            []string `yaml:"file_names"`
	FileDir              string   `yaml:"file_dir"`
	ReqTickerIntervalSec int      `yaml:"req_ticker_interval_sec"`
}

type extra struct {
	EventCount int `yaml:"event_count"`
}

type schedule struct {
	Dsp *dsp
	Sp *sp
	Sd *sd
	Sb *sb
}
type dsp struct {
	Audience []int `yaml:"audience"`
	Detail []int `yaml:"detail"`
	Inventory []int `yaml:"inventory"`
	Order []int `yaml:"order"`
}

type sp struct {
	AdGroups []int `yaml:"adGroups"`
	Asins []int `yaml:"asins"`
	Campaigns []int `yaml:"campaigns"`
	KeywordsQuery []int `yaml:"keywordsQuery"`
	Keywords []int `yaml:"keywords"`
	ProductAds []int `yaml:"productAds"`
	TargetQuerys []int `yaml:"targetQuerys"`
	Targets []int `yaml:"targets"`

}

type sd struct {
	AdGroups []int `yaml:"adGroups"`
	Asins []int `yaml:"asins"`
	Campaigns []int `yaml:"campaigns"`
	Keywords []int `yaml:"keywords"`
	ProductAds []int `yaml:"productAds"`
	Targets []int `yaml:"targets"`
}

type sb struct {
	AdGroup []int `yaml:"adGroup"`
	AdGroupVideo []int `yaml:"adGroupVideo"`
	Campaigns []int `yaml:"campaigns"`
	CampaignsVideo []int `yaml:"campaignsVideo"`
	Keywords []int `yaml:"keywords"`
	KeywordsVideo []int `yaml:"keywordsVideo"`
	KeywordsQuery []int `yaml:"keywordsQuery"`
	KeywordsQueryVideo []int `yaml:"keywordsQueryVideo"`
	Targets []int `yaml:"targets"`
	TargetsVideo []int `yaml:"targetsVideo"`
}



