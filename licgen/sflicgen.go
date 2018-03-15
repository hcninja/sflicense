package main

import (
	"crypto"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	CAPABILITIES_BASE = "PROTECT+CONTROL+VPN+SSL"
	// CAPABILITIES = "PROTECT+CONTROL+URLFilter"
	// CAPABILITIES = "PROTECT+CONTROL+URLFilter+VPN"
	// CAPABILITIES = "PROTECT+CONTROL+URLFilter+VPN+SSL"
)

var licFSM = `model 0x42;
expires forever;
node %s;
serial_number %s;
feature_id 0xC;
model_info 66E:50000:HOST,66E:50000:USER;
66E VirtualDC64bit;
`

var licSensor = `model 0x42;
expires forever;
node %s;
serial_number %s;
feature_id 0xA;
series_3_model_info %s:%s:%s;
%s %s;
`

var licURL = `model 0x42;
expires forever;
node %s;
serial_number %s;
feature_id 0xB;
model_info %s:%s:URLFilter;
%s %s;
license_type SUBSCRIPTION;
`

var licAMP = `model 0x42;
expires forever;
node %s;
serial_number %s;
feature_id 0xA;
series_3_model_info 63L:1:MALWARE;
63L 3D7125;
`

var licHeader = "--- BEGIN SourceFire Product License :\n"
var licFooter = "\n\n--- END SourceFire Product License ---"

func main() {
	licenseKeyFlag := flag.String("l", "", "DC license key [66:00:11:22:33:44:55]")
	privateRsaKeyFlag := flag.String("k", "", "Private RSA key to sign the license")
	modelIdFlag := flag.String("mid", "63G", "Sets the sensor model ID code")
	modelNameFlag := flag.String("mod", "3D7110", "Sets the sensor model number")
	numberOfLicsFlag := flag.String("n", "1", "Number of sensors to apply the license")
	isFsmLicenseFlag := flag.Bool("fsm", false, "Generates a FSM license [default sensor license]")
	isUrlLicenseFlag := flag.Bool("url", false, "Generates a URL subscription license")
	flag.Parse()

	kr := regexp.MustCompile("^([0-9A-Fa-f]{2}:){6}([0-9A-Fa-f]{2})$")

	if !kr.Match([]byte(*licenseKeyFlag)) {
		fmt.Println("No valid license key")
		flag.Usage()
		os.Exit(1)
	}

	if *isFsmLicenseFlag && *isUrlLicenseFlag {
		fmt.Println("fsm and url licenses ar not compatible.")
		os.Exit(1)
	}

	sn := GenSN()
	larr := strings.SplitN(strings.ToUpper(*licenseKeyFlag), ":", 2)

	var licText string
	if *isFsmLicenseFlag {
		licText = fmt.Sprintf(licFSM, larr[1], sn)
	} else if *isUrlLicenseFlag {
		licText = fmt.Sprintf(
			licURL,
			larr[1],
			sn,
			*modelIdFlag,
			*numberOfLicsFlag,
			*modelIdFlag,
			*modelNameFlag,
		)
	} else {
		licText = fmt.Sprintf(
			licSensor,
			larr[1],
			sn,
			*modelIdFlag,
			*numberOfLicsFlag,
			CAPABILITIES_BASE,
			*modelIdFlag,
			*modelNameFlag,
		)
	}

	var lic []byte
	lic = append(lic, []byte(licText)...) // License data

	var signature []byte
	if *privateRsaKeyFlag != "" {
		signature = SignLicense(lic, *privateRsaKeyFlag)
	}

	lic = append(lic, []byte("---")...) // Signature start delimiter
	lic = append(lic, signature...)     // License signature rsa(sha1(license))

	b64lic := base64.StdEncoding.EncodeToString(lic)

	for i, c := range b64lic {
		if i%64 == 0 {
			licHeader += "\n"
		}

		licHeader += string(c)
	}

	licHeader += licFooter

	fmt.Println(licHeader)
}

func GenSN() string {
	rand.Seed(time.Now().Unix())

	var sn string
	for i := 0; i < 9; i++ {
		sn += strconv.Itoa(rand.Intn(9))
	}

	return sn
}

func SignLicense(license []byte, keyFile string) []byte {
	keyBuff, err := ioutil.ReadFile(keyFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rng := crand.Reader

	block, _ := pem.Decode(keyBuff)
	ki, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	hash := sha1.Sum(license)
	sign, err := rsa.SignPKCS1v15(rng, ki, crypto.SHA1, hash[:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return sign
}
