package main

import (
	"github.com/bocchi-the-cache/indeep/internal/apps"
	"github.com/bocchi-the-cache/indeep/internal/servers"
)

func main() { apps.MainServer(servers.Gateway()) }
