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
	"strings"

	"github.com/markkurossi/privx-secrets/api"
	"github.com/markkurossi/privx-secrets/api/vault"
)

func cmdGet(client *api.Client) {
	separator := flag.String("separator", ".", "Data element separator")
	flag.Parse()

	v, err := vault.NewClient(client)
	if err != nil {
		log.Fatalf("failed to create Vault client: %s", err)
	}

	for _, key := range flag.Args() {
		parts := strings.Split(key, *separator)
		bag, err := v.Get(parts[0])
		if err != nil {
			log.Fatalf("failed to get secret '%s': %s", key, err)
		}
		data, ok := bag["data"]
		if !ok {
			log.Fatalf("no 'data' in vault response")
		}
		for i := 1; i < len(parts); i++ {
			switch element := data.(type) {
			case map[string]interface{}:
				el, ok := element[parts[i]]
				if !ok {
					log.Fatalf("element '%s' not found",
						strings.Join(parts[:i+1], *separator))
				}
				data = el

			default:
				log.Fatalf("can't index %T", element)
			}
		}
		fmt.Printf("%s\n", flatten(0, data))
	}
}

func flatten(level int, data interface{}) string {
	switch element := data.(type) {
	case map[string]interface{}:
		var result string
		for k, v := range element {
			if len(result) > 0 {
				result += ","
			}
			result += fmt.Sprintf("%s=%s", k, flatten(level+1, v))
		}
		return fmt.Sprintf("{%s}", result)

	case string:
		if level == 0 {
			return element
		} else {
			return fmt.Sprintf("%q", element)
		}

	default:
		log.Fatalf("can't flatten %T", element)
		return ""
	}
}
