package dex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	// external packages
	"github.com/dexidp/dex/api/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Conn struct {
	DB  *gorm.DB
	Dex api.DexClient
}

func NewConn() (*Conn, error) {
	host := "localhost"
	port := 5432
	user := "postgres"
	password := "postgrespw"
	database := "dex"

	// postgres
	const dsnf = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Seoul"
	dsn := fmt.Sprintf(dsnf, host, port, user, password, database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// grpc
	creds := insecure.NewCredentials()
	conn, err := grpc.Dial("localhost:5557", grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("dial: %v", err)
	}
	return &Conn{
		DB:  db,
		Dex: api.NewDexClient(conn),
	}, nil
}

func (conn *Conn) List(after string, size uint16) ([]Client, uint32) {

	var result struct {
		Rows  string
		Total uint32
	}

	const sql = `SELECT
        (SELECT COUNT(*) FROM client) AS total,
        (SELECT json_agg(client.*) FROM
          (SELECT id, secret, redirect_uris FROM client WHERE id > ? ORDER BY id LIMIT ?) AS client
        ) AS rows;`
	tx := conn.DB.Raw(sql, after, size).Scan(&result)
	if tx.Error != nil {
		panic(tx.Error)
	}
	if tx.RowsAffected == 0 {
		return nil, result.Total
	}

	var items []Client
	if err := json.Unmarshal([]byte(result.Rows), &items); err != nil {
		panic(err)
	}

	return items, result.Total
}

func (conn *Conn) Create(ctx context.Context, client *Client) (*Client, error) {
	req := &api.CreateClientReq{
		Client: &api.Client{
			Id:           client.ID,
			Secret:       client.Secret,
			RedirectUris: client.RedirectURIs,
		},
	}
	res, err := conn.Dex.CreateClient(ctx, req)
	if err != nil {
		panic(err)
	}

	c := res.GetClient()
	if res.AlreadyExists {
		const f = "Client ID (%s) is already exists"
		t := fmt.Sprintf(f, client.ID)
		return nil, errors.New(t)
	}

	return &Client{
		ID:           c.Id,
		Secret:       c.Secret,
		RedirectURIs: c.RedirectUris,
	}, nil
}

func (conn *Conn) Get() {
}

func (conn *Conn) Update() {
}

func (conn *Conn) Delete() {

}
