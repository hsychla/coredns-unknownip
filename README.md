# unknownip

## Name

*unknownip* - check resolved IP against a list of known IPs

## Description

The unknownip plugin checks the IP returned for a A- or AAAA-query against a list of known prefixes and logs a message it is not contained in those.

## Syntax

Create a server block per checked domain

~~~ txt
foo.com:1053 {
    unknownip 198.18.0.0/24 198.19.42.1/32
    forward . 8.8.8.8
}
bar.net:1053 {
    unknownip 198.19.21.0/24 198.19.22.1/32
    forward . 8.8.8.8
}
~~~

## Metrics

If monitoring is enabled (via the *prometheus* directive) the following metric is exported:

* `coredns_unknownip_request_count_total{server}` - query count to the *unknownip* plugin.
* `coredns_unknownip_unknown_count_total{server}` - count of unknown IPs returned.

The `server` label indicated which server handled the request, see the *metrics* plugin for details.

## Ready

This plugin reports readiness to the ready plugin. It will be immediately ready.

## Examples

In this configuration, we forward all queries to 9.9.9.9 and print check if the returned IPs for foo.com and bar.net are included in the given prefixes.

~~~ corefile
foo.com:1053 {
    unknownip 198.18.0.0/24 198.19.42.1/32
    forward . 9.9.9.9
}
bar.net:1053 {
    unknownip 198.19.21.0/24 198.19.22.1/32
    forward . 9.9.9.9
}
. {
  unknownip 
  forward . 9.9.9.9
  example
}
~~~

## Also See

See the [manual](https://coredns.io/manual).
