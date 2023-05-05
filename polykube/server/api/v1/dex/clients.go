package dex

import (
	"fmt"
	"net/http"
	"strconv"

	// external
	"github.com/gin-gonic/gin"
	// project
	"github.com/ghilbut/polykube/pkg/dex"
)

// ListClients godoc
// @Summary      List Dex clients
// @Description  get dex's OIDC client list
// @Tags         dex
// @Produce      json
// @Param        after  query  string  false  "Last cursor position"
// @Param        size   query  int     false  "Max returned item size"
// @Success      200  {object}  ClientList
// @Failure      400  {object}  string
// @Failure      500  {object}  string
// @Router       /v1/dex/clients [get]
// @Security     ApiKey
func ListClients(ctx *gin.Context) {
	v := ctx.MustGet("DEX")
	client := v.(*dex.Conn)

	after := ctx.Query("after")
	ssize := ctx.Query("size")
	if ssize == "" {
		ssize = "20"
	}
	size, err := strconv.ParseInt(ssize, 10, 16)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Size(%s) is invalid value - %v", ssize, err)
		return
	}

	items, total := client.List(after, uint16(size))
	res := ClientList{
		Clients: make([]Client, len(items)),
		Total:   total,
	}
	for index, item := range items {
		c := &res.Clients[index]
		c.ID = item.ID
		c.Secret = item.Secret
		c.RedirectURIs = item.RedirectURIs
	}

	ctx.JSON(http.StatusOK, res)
}

// CreateClient godoc
// @Summary      Create Dex client
// @Description  create dex's new OIDC client
// @Tags         dex
// @Accept       json
// @Produce      json
// @Param        data  body  Client  true  "New Client"
// @Success      200  {object}  string
// @Failure      500  {object}  string
// @Router       /v1/dex/clients [post]
func CreateClient(ctx *gin.Context) {
	v := ctx.MustGet("DEX")
	d := v.(*dex.Conn)
	c := &dex.Client{}
	res, err := d.Create(ctx, c)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ret := &Client{
		ID:           res.ID,
		Secret:       res.Secret,
		RedirectURIs: res.RedirectURIs,
	}
	ctx.JSON(http.StatusOK, ret)
}

// GetClient godoc
// @Summary      Get Dex client
// @Description  get dex's specific OIDC client
// @Tags         dex
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Client ID"
// @Success      200  {object}  string
// @Failure      500  {object}  string
// @Router       /v1/dex/clients/{id} [get]
func GetClient(ctx *gin.Context) {
	id := ctx.Param("id")
	fmt.Println("[ID", id)
	ctx.Status(http.StatusNotImplemented)
}

// UpdateClient godoc
// @Summary      Update Dex client
// @Description  update dex's specific OIDC client
// @Tags         dex
// @Accept       json
// @Produce      json
// @Param        id    path  string  true  "Client ID"
// @Param        data  body  Client  true  "Updated Client Data"
// @Success      200  {object}  string
// @Failure      500  {object}  string
// @Router       /v1/dex/clients/{id} [put]
func UpdateClient(ctx *gin.Context) {
	ctx.Status(http.StatusNotImplemented)
}

// DeleteClient godoc
// @Summary      Delete Dex client
// @Description  delete dex's specific OIDC client
// @Tags         dex
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Client ID"
// @Success      200  {object}  string
// @Failure      500  {object}  string
// @Router       /v1/dex/clients/{id} [delete]
func DeleteClient(ctx *gin.Context) {
	ctx.Status(http.StatusNotImplemented)
}
