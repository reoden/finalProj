package entity

type PrizeBanner struct {
	Picture     string `json:"picture"`
	RedirectUrl string `json:"url"`
}

var prizeAdver *PrizeBanner = &PrizeBanner{
	Picture:     "",
	RedirectUrl: "",
}

func GetPrizeAver() *PrizeBanner {
	return prizeAdver
}

func UpdatePrizeAdver(param *PrizeBanner) {
	prizeAdver.Picture = param.Picture
	prizeAdver.RedirectUrl = param.RedirectUrl
}
