package main

import (
	"github.com/bocchi-the-cache/indeep/internal/apps"
	"github.com/bocchi-the-cache/indeep/internal/servers/metaservers"
)

func main() { apps.MainServer(metaservers.Metaserver()) }
