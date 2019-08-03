package util

import (
	"net"
	"net/url"
)

func ParseIPAddresses(ipsS []string) []net.IP {
	var ipAddresses []net.IP

	for _, ipName := range ipsS {
		ip := net.ParseIP(ipName)
		if ip != nil {
			ipAddresses = append(ipAddresses, ip)
		}
	}

	return ipAddresses
}

func ParseURIs(urisS []string) ([]*url.URL, error) {
	var uris []*url.URL

	for _, uriS := range urisS {
		uri, err := url.Parse(uriS)
		if err != nil {
			return nil, err
		}

		uris = append(uris, uri)
	}

	return uris, nil
}
