package models

type Domain struct {
	Domain string `json:"domain"`
}

type DomainList struct {
	Items []string `json:"items"`
}
