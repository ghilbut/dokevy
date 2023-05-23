package dex

import (
	"encoding/hex"
	"encoding/json"
	"github.com/dexidp/dex/api/v2"
	"net/http"
	"strconv"

	// external
	"github.com/gin-gonic/gin"
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
	conn := getDexConnectionFromContext(ctx)

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
	if size < 1 {
		size = 20
	}

	var result struct {
		Rows  string
		Total uint32
	}

	const sql = `SELECT
        (SELECT COUNT(*) FROM client) AS total,
        (SELECT json_agg(client.*) FROM
          (SELECT id, name, redirect_uris::text FROM client WHERE id > ? ORDER BY id LIMIT ?) AS client
        ) AS rows;`
	tx := conn.DB.Raw(sql, after, size).Scan(&result)
	if tx.Error != nil {
		panic(tx.Error)
	}

	var rows []struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		RedirectURIs string `json:"redirect_uris"`
	}
	if err := json.Unmarshal([]byte(result.Rows), &rows); err != nil {
		panic(err)
	}

	res := ClientList{
		Clients: make([]Client, len(rows)),
		Total:   result.Total,
	}
	for index, row := range rows {
		c := &res.Clients[index]
		c.ID = row.ID
		c.Name = row.Name

		src := []byte(row.RedirectURIs[2:])
		dst := make([]byte, len(src))
		len, err := hex.Decode(dst, src)
		if err != nil {
			panic(err)
		}

		uris := dst[:len]
		if err := json.Unmarshal(uris, &c.RedirectURIs); err != nil {
			panic(err)
		}
	}

	ctx.JSON(http.StatusOK, res)
}

// CreateClient godoc
// @Summary      Create Dex client
// @Description  create dex's new OIDC client
// @Tags         dex
// @Accept       json
// @Produce      json
// @Param        data  body  createReq  true  "New Client"
// @Success      200  {object}  string
// @Failure      500  {object}  string
// @Router       /v1/dex/clients [post]
func CreateClient(ctx *gin.Context) {
	conn := getDexConnectionFromContext(ctx)

	var body createReq
	if err := json.NewDecoder(ctx.Request.Body).Decode(&body); err != nil {
		ctx.String(http.StatusBadRequest, "Failed json decoding - %v", err)
		return
	}

	req := &api.CreateClientReq{
		Client: &api.Client{
			Id:           body.ID,
			Name:         body.Name,
			RedirectUris: body.RedirectURIs,
		},
	}
	res, err := conn.API.CreateClient(ctx, req)
	if err != nil {
		panic(err)
	}

	c := res.GetClient()
	if res.AlreadyExists {
		const f = "Client ID (%s) is already exists"
		ctx.String(http.StatusBadRequest, f, body.ID)
	}

	ctx.JSON(http.StatusOK, &Client{
		ID:           c.Id,
		Secret:       c.Secret,
		Name:         c.Name,
		RedirectURIs: c.RedirectUris,
	})
}

type createReq struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	RedirectURIs []string `json:"redirectURIs"`
}

// GetClient godoc
// @Summary      Get Dex client
// @Description  get dex's specific OIDC client
// @Tags         dex
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Client ID"
// @Success      200  {object}  Client
// @Failure      500  {object}  string
// @Router       /v1/dex/clients/{id} [get]
func GetClient(ctx *gin.Context) {
	conn := getDexConnectionFromContext(ctx)

	id := ctx.Param("id")
	entity := &ClientEntity{}
	const sql = `SELECT id, secret, name, redirect_uris FROM client WHERE id = ?;`
	tx := conn.DB.Raw(sql, id).Scan(entity)
	if tx.Error != nil {
		panic(tx.Error)
	}
	if tx.RowsAffected == 0 {
		ctx.Status(http.StatusNotFound)
		return
	}

	var uris []string
	if err := json.Unmarshal([]byte(entity.RedirectURIs), &uris); err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, &Client{
		ID:           entity.ID,
		Name:         entity.Name,
		RedirectURIs: uris,
	})
}

// UpdateClient godoc
// @Summary      Update Dex client
// @Description  update dex's specific OIDC client
// @Tags         dex
// @Accept       json
// @Produce      json
// @Param        id    path  string     true  "Client ID"
// @Param        data  body  updateReq  true  "Client Data to update"
// @Success      200  {object}  string
// @Failure      400  {object}  string
// @Failure      404  {object}  string
// @Failure      500  {object}  string
// @Router       /v1/dex/clients/{id} [put]
func UpdateClient(ctx *gin.Context) {
	conn := getDexConnectionFromContext(ctx)

	id := ctx.Param("id")
	client := &Client{
		ID: id,
	}
	if err := json.NewDecoder(ctx.Request.Body).Decode(&client); err != nil {
		ctx.String(http.StatusBadRequest, "Failed json decoding - %v", err)
		return
	}

	req := &api.UpdateClientReq{
		Name:         client.Name,
		Id:           client.ID,
		RedirectUris: client.RedirectURIs,
	}
	res, err := conn.API.UpdateClient(ctx, req)
	if err != nil {
		panic(err)
	}
	if res.NotFound {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.Status(http.StatusOK)
}

type updateReq struct {
	Name         string   `json:"name"`
	RedirectURIs []string `json:"redirectURIs"`
}

// DeleteClient godoc
// @Summary      Delete Dex client
// @Description  delete dex's specific OIDC client
// @Tags         dex
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Client ID"
// @Success      200  {object}  string
// @Failure      404  {object}  string
// @Failure      500  {object}  string
// @Router       /v1/dex/clients/{id} [delete]
func DeleteClient(ctx *gin.Context) {
	conn := getDexConnectionFromContext(ctx)

	id := ctx.Param("id")
	req := &api.DeleteClientReq{
		Id: id,
	}
	res, err := conn.API.DeleteClient(ctx, req)
	if err != nil {
		panic(err)
	}
	if res.NotFound {
		ctx.Status(http.StatusNotFound)
		return
	}
	ctx.Status(http.StatusOK)
}
