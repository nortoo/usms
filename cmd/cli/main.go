package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nortoo/usms/internal/app/cli"
)

func main() {
	host := flag.String("host", "127.0.0.1", "host")
	port := flag.Int("port", 8080, "port")
	servername := flag.String("servername", "", "server name")
	certPath := flag.String("cert", "", "certificate path")
	keyPath := flag.String("key", "", "private key path")
	caPath := flag.String("ca", "", "ca path")

	username := flag.String("u", "", "username")
	password := flag.String("p", "", "password")
	email := flag.String("e", "", "email")
	tenant := flag.String("t", "", "tenant")

	createAll := flag.Bool("create-all", false, "create all nonexisting entities")

	flag.Parse()

	if *servername == "" {
		fmt.Println("servername is required")
		flag.Usage()
		os.Exit(1)
	}
	if *certPath == "" {
		fmt.Println("cert is required")
		flag.Usage()
		os.Exit(1)
	}
	if *keyPath == "" {
		fmt.Println("key is required")
		flag.Usage()
		os.Exit(1)
	}
	if *caPath == "" {
		fmt.Println("ca is required")
		flag.Usage()
		os.Exit(1)
	}
	if *username == "" {
		fmt.Println("username is required")
		flag.Usage()
		os.Exit(1)
	}
	if *password == "" {
		fmt.Println("password is required")
		flag.Usage()
		os.Exit(1)
	}
	if *email == "" {
		fmt.Println("email is required")
		flag.Usage()
		os.Exit(1)
	}
	if *tenant == "" {
		fmt.Println("tenant is required")
		flag.Usage()
		os.Exit(1)
	}

	client, err := cli.NewClient(*certPath, *keyPath, *caPath, *servername, *host, *port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	user, err := client.CreateAdministration(*username, *password, *tenant, *email, *createAll)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Create Administration:")
	fmt.Printf("Username: %s\n", user.Username)
	fmt.Printf("Password: %s\n", *password)
	fmt.Printf("Email: %s\n", *email)
	fmt.Printf("Tenant: %s\n", *tenant)
}
