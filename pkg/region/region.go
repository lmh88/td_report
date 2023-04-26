package region

const (
	RegionNa  = "NA"
	RegionEu   = "EU"
	RegionFe   = "FE"
)

var ApiUrl = map[string]string{
	RegionFe: "https://advertising-api-fe.amazon.com",
	RegionEu: "https://advertising-api-eu.amazon.com",
	RegionNa: "https://advertising-api.amazon.com",
}



