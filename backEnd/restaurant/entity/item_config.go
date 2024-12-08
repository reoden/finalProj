package entity

type ItemStatus int

const (
	ItemShow ItemStatus = iota
	ItemHide
)

type ItemConfig struct {
	App       ItemStatus `json:"app"`
	Chef      ItemStatus `json:"chef"`
	Guest     ItemStatus `json:"guest"`
	News      ItemStatus `json:"news"`
	Prize     ItemStatus `json:"prize"`
	AppName   string     `json:"app_name"`
	ChefName  string     `json:"chef_name"`
	GuestName string     `json:"guest_name"`
	NewsName  string     `json:"news_name"`
	PrizeName string     `json:"prize_name"`
}
