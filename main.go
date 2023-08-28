package main

import (
	"fmt"
	"github.com/JackalLabs/jackalapi/japicore"
	"github.com/rs/cors"
	"github.com/uptrace/bunrouter"
	"net/http"
	"os"
	"strconv"
)

func main() {
	_, fileIo := japicore.InitWalletSession()
	queue := japicore.NewQueue()

	router := bunrouter.New(
		bunrouter.WithMethodNotAllowedHandler(japicore.MethodNotAllowedHandler()),
	)
	group := router.NewGroup("")

	handler := http.Handler(router)
	handler = cors.Default().Handler(handler)

	group.WithGroup("", func(group *bunrouter.Group) {
		group.GET("/version", japicore.VersionHandler())
		group.GET("/download/:id", japicore.DownloadHandler(fileIo))
		group.GET("/d/:id", japicore.DownloadHandler(fileIo))
		group.GET("/ipfs/:id", japicore.IpfsHandler(fileIo))

		group.POST("/import", japicore.ImportHandler(fileIo))
		group.POST("/upload", japicore.UploadHandler(fileIo, queue))
		group.POST("/u", japicore.UploadHandler(fileIo, queue))
		group.DELETE("/del/:id", japicore.DeleteHandler(fileIo))
	})

	os.Setenv("JHTTP_PORT", "1234")

	port := os.Getenv("JHTTP_PORT")
	if len(port) == 0 {
		port = "3535"
	}

	portNum, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		panic(err)
	}

	fmt.Printf("🌍 Started JHN: http://0.0.0.0:%d\n", portNum)
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", portNum), handler)
	if err != nil {
		panic(err)
	}

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server Closed\n")
		return
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
