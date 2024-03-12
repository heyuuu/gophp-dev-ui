package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"gophp-dev-ui/router"
	"log"
	"os"
)

func main() {
	var workdir string
	flag.StringVar(&workdir, "w", "", "workdir")
	flag.Parse()
	fmt.Println("\n>>> workdir=" + workdir)
	if workdir != "" {
		err := os.Chdir(workdir)
		if err != nil {
			log.Panicln(err)
		}
	}

	r := gin.Default()
	router.InitRouter(r)
	r.Run(":8080")
}
