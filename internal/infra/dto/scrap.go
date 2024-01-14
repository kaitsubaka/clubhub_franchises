package dto

import whoisparser "github.com/likexian/whois-parser"

type ScrapDTO struct {
	HTMLData HTMLDataDTO
	ProtocolAndJumpsDTO
	WhoisData whoisparser.WhoisInfo
}

type ProtocolAndJumpsDTO struct {
	Protocol string
	Jumps    int
}

type HTMLDataDTO struct {
	Title       string
	Description string
	Image       string
	SiteName    string
}

type SSLLabsResponseDTO struct {
	Host      string        `json:"host"`
	Port      int           `json:"port"`
	Protocol  string        `json:"protocol"`
	Endpoints []EndpointDTO `json:"endpoints"`
}

type EndpointDTO struct {
	IpAddress string `json:"ipAddress"`
	Grade     string `json:"grade"`
	SeverName string `json:"serverName"`
}
