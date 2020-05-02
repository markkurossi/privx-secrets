//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/markkurossi/privx-secrets/api"
)

func cmdLogin(client *api.Client) {
	flag.Parse()

	_, err := client.Auth.Token()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ok")
}
