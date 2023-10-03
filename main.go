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
	coreSession := japicore.InitJApiCore()

	router := bunrouter.New(
		bunrouter.WithMethodNotAllowedHandler(coreSession.MethodNotAllowedHandler()),
		bunrouter.WithNotFoundHandler(coreSession.RouteNotFoundHandler()),
	)
	group := router.NewGroup("")

	handler := http.Handler(router)
	handler = cors.Default().Handler(handler)

	group.WithGroup("", func(group *bunrouter.Group) {
		group.GET("/version", coreSession.VersionHandler())
	})

	group.WithGroup("/fid", func(group *bunrouter.Group) {
		group.GET("/download/:id", coreSession.DownloadByFidHandler())
		group.GET("/d/:id", coreSession.DownloadByFidHandler())
		group.GET("/ipfs/:id", coreSession.IpfsHandler())

		group.POST("/upload", coreSession.UploadByPathHandler())
		group.POST("/u", coreSession.UploadByPathHandler())
		group.DELETE("/del/:id", coreSession.DeleteByFidHandler())
	})

	group.WithGroup("/p", func(group *bunrouter.Group) {
		group.GET("/downloadfrombulk/*location", coreSession.BasicDownloadFromBulkByPathHandler())
		group.GET("/download/*location", coreSession.BasicDownloadByPathHandler())
		group.GET("/d/*location", coreSession.BasicDownloadByPathHandler())

		group.POST("/import", coreSession.ImportHandler())
		group.POST("/upload", coreSession.UploadByPathHandler())
		group.POST("/u", coreSession.UploadByPathHandler())
		group.DELETE("/delfrombulk/:filename/*location", coreSession.BasicDeleteFromBulkByPathHandler())
		group.DELETE("/del/:filename/*location", coreSession.BasicDeleteByPathHandler())
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

	fmt.Printf("🌍 JAPI Wallet: %s\n", coreSession.Wallet.GetAddress())
	fmt.Printf("🌍 JAPI Network: %s\n", coreSession.Wallet.GetChainID())
	fmt.Println("<<<<< * >>>>>")

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server Closed\n")
		return
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
