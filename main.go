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
	wallet, fileIo := japicore.InitWalletSession()
	fileIoQueue := japicore.NewFileIoQueue()

	scrapeQueue := japicore.NewScrapeQueue(fileIoQueue)

	router := bunrouter.New(
		bunrouter.WithMethodNotAllowedHandler(japicore.MethodNotAllowedHandler()),
		bunrouter.WithNotFoundHandler(japicore.RouteNotFoundHandler()),
	)
	group := router.NewGroup("")

	handler := http.Handler(router)
	handler = cors.Default().Handler(handler)

	group.WithGroup("", func(group *bunrouter.Group) {
		group.GET("/version", japicore.VersionHandler())
	})

	group.WithGroup("/fid", func(group *bunrouter.Group) {
		group.GET("/download/:id", japicore.DownloadByFidHandler(fileIo))
		group.GET("/d/:id", japicore.DownloadByFidHandler(fileIo))
		group.GET("/ipfs/:id", japicore.IpfsHandler(fileIo, fileIoQueue))

		group.POST("/upload", japicore.UploadByPathHandler(fileIo, fileIoQueue))
		group.POST("/u", japicore.UploadByPathHandler(fileIo, fileIoQueue))
		group.DELETE("/del/:id", japicore.DeleteByFidHandler(fileIo, fileIoQueue))
	})

	group.WithGroup("/p", func(group *bunrouter.Group) {
		group.GET("/download/:id", japicore.DownloadByPathHandler(fileIo))
		group.GET("/d/:id", japicore.DownloadByPathHandler(fileIo))

		group.POST("/import", japicore.ImportHandler(fileIo, scrapeQueue))
		group.POST("/upload", japicore.UploadByPathHandler(fileIo, fileIoQueue))
		group.POST("/u", japicore.UploadByPathHandler(fileIo, fileIoQueue))
		group.DELETE("/del/:id", japicore.DeleteByPathHandler(fileIo, fileIoQueue))
	})

	port := jutils.LoadEnvVarOrFallback("JAPI_PORT", "3535")

	portNum, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		panic(err)
	}

	fmt.Println("<<<<< * >>>>>")
	fmt.Printf("🌍 Started JAPI: http://0.0.0.0:%d\n", portNum)
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", portNum), handler)
	if err != nil {
		panic(err)
	}

	fmt.Printf("🌍 JAPI Wallet: %s\n", wallet.GetAddress())
	fmt.Printf("🌍 JAPI Network: %s\n", wallet.GetChainID())
	fmt.Println("<<<<< * >>>>>")

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server Closed\n")
		return
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
