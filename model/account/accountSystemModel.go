package account

type Account struct {
	Account    string `json:"account"`
	Password   string `json:"password"`
	Name       string `json:"name"`
	Type       int    `json:"type"`
	CreateTime string `json:"create_time"`
}

type BindingDevice struct {
	Account    string `json:"account"`
	DeviceId   string `json:"device_id"`
	CreateTime string `json:"create_time"`
}
