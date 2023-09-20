package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/JackalLabs/jackalapi/japicore"
	"github.com/JackalLabs/jackalapi/jutils"
	"github.com/rs/cors"
	"github.com/uptrace/bunrouter"
)

func main() {
	_, fileIo := japicore.InitWalletSession()
	fileIoQueue := japicore.NewFileIoQueue()

	scrapeQueue := japicore.NewScrapeQueue(fileIoQueue)

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
		group.GET("/ipfs/:id", japicore.IpfsHandler(fileIo, fileIoQueue))

		group.POST("/import", japicore.ImportHandler(fileIo, scrapeQueue))

		// upload func gives an example of a path generator function
		uploadFunc := japicore.UploadHandler(fileIo, fileIoQueue, func(req *bunrouter.Request) (string, error) {
			fmt.Println(req.URL) // has access to the request
			return "example path", nil
		})
		group.POST("/upload", uploadFunc)
		group.POST("/u", uploadFunc)

		group.DELETE("/del/:id", japicore.DeleteHandler(fileIo, fileIoQueue))
	})

	port := jutils.LoadEnvVarOrFallback("JAPI_PORT", "3535")

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
