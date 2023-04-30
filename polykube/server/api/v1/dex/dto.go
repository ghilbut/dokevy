package dex

type Client struct {
	ClientID     string   `json:"clientId"`
	ClientSecret string   `json:"clientSecret"`
	RedirectURIs []string `json:"redirectURIs"`
}
