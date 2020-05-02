//
// Copyright (c) 2020 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/markkurossi/privx-secrets/api"
	"github.com/markkurossi/privx-secrets/oauth"
)

const (
	etcDir         = "/opt/etc/privx-secrets"
	configFileName = "privx-secrets.toml"
)

var (
	config Config
)

type Config struct {
	Verbose bool
	API     APIConfig
	Auth    oauth.Config
}

type APIConfig struct {
	Endpoint    string
	Certificate *Certificate
}

type Certificate struct {
	X509 *x509.Certificate
}

func (cert *Certificate) UnmarshalText(text []byte) error {
	block, _ := pem.Decode(text)
	if block == nil {
		return fmt.Errorf("could not decode certificate PEM data")
	}
	c, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}
	cert.X509 = c
	return nil
}

var commands = map[string]func(client *api.Client){
	"login": cmdLogin,
}

func main() {
	var defaultConfig string

	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("failed to get user's home directory: %s", err)
		defaultConfig = path.Join(etcDir, configFileName)
		log.Printf("fallback to '%s'", defaultConfig)
	} else {
		defaultConfig = path.Join(home, fmt.Sprintf(".%s", configFileName))
	}

	configFile := flag.String("config", defaultConfig, "configuration file")
	flag.Parse()

	if err := readConfig(*configFile); err != nil {
		log.Fatalf("failed to read config file '%s': %s", *configFile, err)
	}

	// Construct API client.
	auth, err := oauth.NewClient(config.Auth, config.API.Endpoint,
		config.API.Certificate.X509, config.Verbose)
	if err != nil {
		log.Fatal(err)
	}
	client, err := api.NewClient(auth)
	if err != nil {
		log.Fatal(err)
	}

	if len(flag.Args()) == 0 {
		fmt.Printf("No command specified. Possible commands are\n")
		for key := range commands {
			fmt.Printf(" - %s\n", key)
		}
		return
	}
	os.Args = flag.Args()
	fn, ok := commands[flag.Arg(0)]
	if !ok {
		fmt.Printf("Unknown command: %s\n", flag.Arg(0))
		os.Exit(1)
	}
	flag.CommandLine = flag.NewFlagSet(
		fmt.Sprintf("privx-secrets %s", os.Args[0]),
		flag.ExitOnError)
	fn(client)
}

func readConfig(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	return nil
}
