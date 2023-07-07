package domain

type Locality struct {
	ID       int    `json:"id"`
	Name     string `json:"locality_name"`
	Province string `json:"province_name"`
	Country  string `json:"country_name"`
}
