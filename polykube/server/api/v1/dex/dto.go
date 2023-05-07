package dex

type Client struct {
	ID           string   `json:"id"`
	Secret       string   `json:"secret,omitempty"`
	Name         string   `json:"name"`
	RedirectURIs []string `json:"redirect_uris"`
}

type ClientList struct {
	Clients []Client `json:"clients"`
	Total   uint32   `json:"total"`
}
