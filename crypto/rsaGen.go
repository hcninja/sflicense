/*
 * Genarate rsa keys.
 */

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	reader := rand.Reader
	bitSize := 4096

	ki, err := rsa.GenerateKey(reader, bitSize)
	checkError(err)

	kp := ki.PublicKey

	savePEMKey("private.pem", ki)
	savePublicDERKey("public.der", kp)
}

func savePEMKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	checkError(err)
}

func savePublicDERKey(fileName string, pubkey rsa.PublicKey) {
	asn1Bytes, err := asn1.Marshal(pubkey)
	checkError(err)

	derfile, err := os.Create(fileName)
	checkError(err)
	defer derfile.Close()

	derfile.Write(asn1Bytes)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
