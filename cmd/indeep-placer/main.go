package main

import (
	"os"

	"github.com/bocchi-the-cache/indeep/internal/apps"
	"github.com/bocchi-the-cache/indeep/internal/servers"
)

func main() { apps.RunServer(servers.DefaultPlacer(), os.Args[1:]) }
