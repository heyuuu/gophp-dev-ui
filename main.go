package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"gophp-dev-ui/router"
	"log"
	"os"
	"strconv"
)

func main() {
	// flag.Parse
	var workdir string
	var port int
	flag.StringVar(&workdir, "w", "", "workdir")
	flag.IntVar(&port, "p", 8080, "port")
	flag.Parse()
	fmt.Println("\n>>> workdir=" + workdir)
	if workdir != "" {
		err := os.Chdir(workdir)
		if err != nil {
			log.Panicln(err)
		}
	}
	if port < 1 || port > 65535 {
		log.Panicln("port must between 1~65535")
	}

	// run
	r := gin.Default()
	router.InitRouter(r)
	r.Run(":" + strconv.Itoa(port))
}
