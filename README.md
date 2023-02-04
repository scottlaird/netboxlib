# netboxlib

netboxlib is an overly-simple wrapper around
[`go-netbox`]("github.com/netbox-community/go-netbox), which is an
overly-complicated auto-generated library for talking to
[Netbox](http://netbox.dev)'s API.

This is generally focused on reading all records of a given type from
Netbox and creating structs with most of the useful fields populated.
It is in no way, shape, or form intended as a full replacement for
`go-netbox` or any other library.  Instead, it exists to give some of
my other tools a lighter-weight method of accessing Netbox.
