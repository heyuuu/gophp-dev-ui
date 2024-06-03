package main

import (
	"embed"
	"flag"
	"github.com/gin-gonic/gin"
	"gophp-dev-ui/router"
	"log"
	"strconv"
)

//go:embed static
var staticFs embed.FS

func main() {
	// flag.Parse
	var port int
	flag.IntVar(&port, "p", 8080, "port")
	flag.Parse()
	if port < 1 || port > 65535 {
		log.Panicln("port must between 1~65535")
	}

	// run
	r := gin.Default()
	router.InitRouter(r, staticFs)
	err := r.Run(":" + strconv.Itoa(port))
	if err != nil {
		panic(err)
	}
}
