package main

import (
	"github.com/ezavalishin/partygames/internal/orm"
	"github.com/ezavalishin/partygames/pkg/server"

	log "github.com/ezavalishin/partygames/internal/logger"
)

func main() {

	ormer, err := orm.Factory()

	if err != nil {
		log.Panic(err)
	}

	server.Run(ormer)
}
