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
	cshell := flag.Bool("c", false, "Generate C-shell commands on stdout")
	bourne := flag.Bool("s", false, "Generate Bourne shell commands on stdout")
	flag.Parse()

	v, err := vault.NewClient(client)
	if err != nil {
		log.Fatalf("failed to create Vault client: %s", err)
	}

	for _, key := range flag.Args() {
		parts := strings.Split(key, *separator)

		name := parts[0]
		var kv bool
		var env string
		idx := strings.IndexRune(name, '=')
		if idx > 0 {
			env = name[:idx]
			name = name[idx+1:]
			kv = true
		}

		bag, err := v.Get(name)
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

		if kv {
			value := flatten(true, 1, data)
			if *cshell {
				fmt.Printf("setenv %s %s\n", env, value)
			} else if *bourne {
				fmt.Printf("%s=%s; export %s;\n", env, value, env)
			} else {
				fmt.Printf("%s=%s\n", env, value)
			}
		} else {
			fmt.Printf("%s\n", flatten(false, 0, data))
		}
	}
}

func flatten(all bool, level int, data interface{}) string {
	switch element := data.(type) {
	case map[string]interface{}:
		var result string
		for k, v := range element {
			if len(result) > 0 {
				result += ","
			}
			result += fmt.Sprintf("%s=%s", k, flatten(all, level+1, v))
		}
		if all {
			return fmt.Sprintf("'{%s}'", result)
		} else {
			return fmt.Sprintf("{%s}", result)
		}

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
