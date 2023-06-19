package unknownip

import (
	"net/netip"
	"strings"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

// init registers this plugin.
func init() { plugin.Register("unknownip", setup) }

// store prefixes per zone
type prefixList map[string][]netip.Prefix

var prefixLists = make(prefixList)

// setup is the function that gets called when the config parser see the token "unknownip". Setup is responsible
// for parsing any extra options the unknownip plugin may have. The first token this function sees is "unknownip".
func setup(c *caddy.Controller) error {

	// extract zone from Key (dns://foo.tld.:53 => foo.tld.)
	zone := strings.Split(strings.Split(c.Key, "/")[2], ":")[0]

	c.Next() // Ignore "unknownip" and give us the next token.
	args := c.RemainingArgs()

	for i := 0; i < len(args); i++ {
		prefix, pErr := netip.ParsePrefix(args[i])
		if pErr != nil {
			return plugin.Error("unknownip", pErr)
		}
		prefixLists[zone] = append(prefixLists[zone], prefix)
	}

	// Add the Plugin to CoreDNS, so Servers can use it in their plugin chain.
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Unknownip{prefixLists: prefixLists, Next: next}
	})

	// All OK, return a nil error.
	return nil
}
