package main

import (
	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	// plug in Caddy modules here
	_ "github.com/caddyserver/caddy/v2/modules/standard"
	_ "github.com/fore-stun/spanx"
	_ "github.com/greenpau/caddy-security"
	_ "github.com/lindenlab/caddy-s3-proxy"
)

func main() {
	caddycmd.Main()
}
