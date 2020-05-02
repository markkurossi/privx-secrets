//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package api

import (
	"github.com/markkurossi/privx-secrets/oauth"
)

type Client struct {
	Auth *oauth.Client
}

func NewClient(auth *oauth.Client) (*Client, error) {
	return &Client{
		Auth: auth,
	}, nil
}
