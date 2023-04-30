package dex

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ListClients godoc
// @Summary      List Dex clients
// @Description  get dex's OIDC client list
// @Tags         dex
// @Accept       json
// @Produce      json
// @Success      200  {object}  string
// @Failure      500  {object}  string
// @Router       /v1/dex/clients [get]
// @Security     ApiKey
func ListClients(ctx *gin.Context) {
	ctx.Status(http.StatusNotImplemented)
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
	ctx.Status(http.StatusNotImplemented)
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
