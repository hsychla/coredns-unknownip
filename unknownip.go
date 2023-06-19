// Package unknownip is a CoreDNS plugin that logs if the returned IP for a query is not in
// a predefined list.
package unknownip

import (
	"context"
	"fmt"
	"net/netip"
	"strings"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/plugin/pkg/nonwriter"

	"github.com/miekg/dns"
)

// Define log to be a logger with the plugin name in it. This way we can just use log.Info and
// friends to log.
var log = clog.NewWithPlugin("unknownip")

type Unknownip struct {
	Next        plugin.Handler
	prefixLists prefixList
}

// We only care about dns.TypeA and dns.TypeAAAA requests
var allowedQTypes = []uint16{1, 28}

func contains(elems []uint16, v uint16) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

// ServeDNS implements the plugin.Handler interface. This method gets called when unknownip is used
// in a Server.
func (u Unknownip) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	log.Debug(fmt.Sprintf("Enter func unknownIP.ServeDNS() for Question %v", r.Question))

	// Get zone from Question
	zone := r.Question[0].Name
	log.Debug(zone)

	// Check if the request is really for one of our zones
	// Just call next plugin if not
	_, zoneExists := prefixLists[zone]
	if !zoneExists {
		return plugin.NextOrFailure(u.Name(), u.Next, ctx, w, r)
	}

	// Use a nonwriter to capture the response.
	nw := nonwriter.New(w)

	// Call next plugin (if any).
	rcode, err := plugin.NextOrFailure(u.Name(), u.Next, ctx, nw, r)
	if err != nil {
		// Simply return if there was an error.
		return rcode, err
	}

	// We now know we have a valid response so we can
	// check for unknown IPs
	if len(nw.Msg.Answer) > 0 && (contains(allowedQTypes, nw.Msg.Question[0].Qtype)) {

		// Export metric with the server label set to the current server handling the request.
		requestCount.WithLabelValues(metrics.WithServer(ctx)).Inc()

		// Get ip from Answer
		ip, ipErr := netip.ParseAddr(strings.Fields(nw.Msg.Answer[0].String())[4])
		if ipErr != nil {
			// We don't have an error code so just use 0?
			return 0, plugin.Error("unknownip", ipErr)
		}
		log.Debug(fmt.Sprintf("Query returned IP %s", ip.String()))

		// Check if returned IP is in list of known IPs
		ipvalid := false
		for _, prefix := range prefixLists[zone] {
			if prefix.Contains(ip) {
				ipvalid = true
				log.Debug(fmt.Sprintf("IP %s is valid for zone %s", ip.String(), zone))
				break
			}
		}

		if !ipvalid {
			// Export metric with the server label set to the current server handling the request.
			unknownIpCount.WithLabelValues(metrics.WithServer(ctx)).Inc()
			log.Info(fmt.Sprintf("IP %s is NOT valid for zone %s", ip.String(), zone))
		}
	}

	// Return answer to client
	w.WriteMsg(nw.Msg)

	return rcode, err
}

// Name implements the Handler interface.
func (e Unknownip) Name() string { return "unknownip" }
