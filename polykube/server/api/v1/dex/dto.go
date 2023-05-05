package dex

type Client struct {
	ID           string   `json:"clientId"`
	Secret       string   `json:"clientSecret"`
	RedirectURIs []string `json:"redirectURIs"`
}

type ClientList struct {
	Clients []Client `json:"clients"`
	Total   uint32   `json:"total"`
}
