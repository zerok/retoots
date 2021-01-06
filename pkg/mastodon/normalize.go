package mastodon

import (
	"fmt"
	"net/url"
	"strings"
)

func NormalizeAcct(server string, acct string) string {
	if strings.Contains(acct, "@") {
		return acct
	}
	u, err := url.Parse(server)
	if err != nil {
		return "<err>"
	}
	return fmt.Sprintf("%s@%s", acct, u.Host)
}
