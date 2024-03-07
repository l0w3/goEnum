package webdiscovery

import (
	"net/http"
	"strconv"
)

type WebPages struct {
	Target string
	Port   int
}

func Resolve(target string, port int, protocol string) (WebPages, int) {
	url := protocol + "://" + target + ":" + strconv.Itoa(port)
	_, err := http.Head(url)
	if err == nil {
		return WebPages{Port: port, Target: target}, 0
	}
	return WebPages{}, 1
}
