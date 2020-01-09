package models

//Information struct
type Information struct {
	ServersChanged   bool     `json:"server_changed"`
	SslGrade         string   `json:"ssl_grade"`
	PreviousSslGrade string   `json:"previous_ssl_grade"`
	Logo             string   `json:"logo"`
	Title            string   `json:"title"`
	IsDown           bool     `json:"is_down"`
	Servers          []Server `json:"servers"`
	// ID      int64  `json:"id"`

}

//Server struct
type Server struct {
	Address  string `json:"address"`
	SSLGrade string `json:"ssl_grade"`
	Country  string `json:"country"`
	Owner    string `json:"owner"`
}
