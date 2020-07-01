package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/ardanlabs/service/business/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

func main() {
	// if err := keyGen(); err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }

	if err := gentoken(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func gentoken() error {
	privatePEM, err := ioutil.ReadFile("private.pem")
	if err != nil {
		return errors.Wrap(err, "reading auth private key")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return errors.Wrap(err, "parsing auth private key")
	}

	keyLookupFunc := func(kid string) (*rsa.PublicKey, error) {
		switch kid {
		case "1":
			return privateKey.Public().(*rsa.PublicKey), nil
		}
		return nil, fmt.Errorf("no public key found for the specified kid: %s", kid)
	}
	a, err := auth.New(privateKey, "1", "RS256", keyLookupFunc)
	if err != nil {
		return errors.Wrap(err, "constructing auth")
	}

	claims := auth.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "service project",
			Subject:   "12345",
			ExpiresAt: time.Now().Add(8760 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Roles: []string{"ADMIN"},
	}

	str, err := a.GenerateToken(claims)
	if err != nil {
		return err
	}

	fmt.Println("-------  BEGIN TOKEN -------")
	fmt.Println(str)
	fmt.Println("-------  END TOKEN -------")

	return nil
}

func keyGen() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	privateFile, err := os.Create("private.pem")
	if err != nil {
		return errors.Wrap(err, "creating private file")
	}
	defer privateFile.Close()

	privateBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	if err := pem.Encode(privateFile, &privateBlock); err != nil {
		return errors.Wrap(err, "encoding to private file")
	}

	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return errors.Wrap(err, "marshaling public key")
	}

	publicFile, err := os.Create("public.pem")
	if err != nil {
		return errors.Wrap(err, "creating public file")
	}
	defer privateFile.Close()

	publicBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	if err := pem.Encode(publicFile, &publicBlock); err != nil {
		return errors.Wrap(err, "encoding to public file")
	}

	return nil
}
