package restapi

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	//"github.com/labstack/gommon/log"
	"github.com/vishvananda/netlink"
)

type RestAPI struct {
	echo *echo.Echo
}

func NewRestApi(port int) (RestAPI, error) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	api := RestAPI{
		echo: e,
	}

	// Routes
	e.GET("/help", help)
	e.GET("/interfaces", api.getInterfaces)
	e.GET("/interfaces/:name", api.getInterfaces)

	log.Infof("Starting echo REST API on port %d", port)

	go func() {
		//err := e.StartTLS(fmt.Sprintf(":%d", port), cert.Certificate, xxxx)
		err := e.Start(fmt.Sprintf(":%d", port))
		if err != nil {
			log.Errorf("Start echo failed with [%s]", err.Error())
			panic(err.Error())
		}
	}()

	return api, nil
}

func (r RestAPI) Close() {
	if r.echo != nil {
		r.echo.Close()
		r.echo = nil
	}
}

func help(c echo.Context) error {
	hostname, _ := os.Hostname()
	banner := fmt.Sprintf("collector alive on [%s] (%s on %s/%s). Available routes: ",
		hostname, runtime.Version(), runtime.GOOS, runtime.GOARCH)

	routes := c.Echo().Routes()
	for _, route := range routes {
		banner = banner + "<li>" + route.Path + "</li> "
	}

	return c.HTML(http.StatusOK, banner)
}

type InterfaceStats struct {
	Name string
	Rx   uint64
	Tx   uint64
}

func (api *RestAPI) getInterfaces(c echo.Context) error {
	name := c.Param("name")
	var ret interface{}

	if name != "" {
		l, err := netlink.LinkByName(name)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Interface [%s] not found", name))
		}

		stats := l.Attrs().Statistics
		ret = InterfaceStats{
			Name: l.Attrs().Name,
			Rx:   stats.RxBytes,
			Tx:   stats.TxBytes,
		}
	} else {
		links, err := netlink.LinkList()
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Interfaces not found"))
		}

		arraystats := make([]InterfaceStats, 0)
		for _, l := range links {
			stats := l.Attrs().Statistics
			is := InterfaceStats{
				Name: l.Attrs().Name,
				Rx:   stats.RxBytes,
				Tx:   stats.TxBytes,
			}

			arraystats = append(arraystats, is)
		}

		ret = arraystats
	}

	return c.JSON(http.StatusOK, ret)
}
