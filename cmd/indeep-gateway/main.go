package main

import (
	"github.com/bocchi-the-cache/indeep/internal/apps"
	"github.com/bocchi-the-cache/indeep/internal/servers/gateways"
)

func main() { apps.MainServer(gateways.Gateway()) }
