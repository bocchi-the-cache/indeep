package main

import (
	"github.com/bocchi-the-cache/indeep/internal/apps"
	"github.com/bocchi-the-cache/indeep/internal/servers/placers"
)

func main() { apps.MainServer(placers.Placer()) }
