package sdkuseragent

import (
	"net/http"
	"strings"
)

const userAgent = "SaaSus-SDK-go/" + Version

// Set appends the SDK identifier to the request's User-Agent once.
func Set(req *http.Request) {
	current := strings.Join(req.Header.Values("User-Agent"), " ")
	for _, product := range strings.Fields(current) {
		if product == userAgent {
			return
		}
	}
	req.Header.Set("User-Agent", strings.TrimSpace(current+" "+userAgent))
}
