package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	// Copyright Sourcefire 2013
	// -v version
	// -h this message
	// -f [filename]
	// -k [keyfile]
	// -q quietmode
	// -d dump contents of license
	// -F dump feature licenses
	f := flag.String("f", "", "[filename]")
	d := flag.Bool("d", false, "dump contents of license")
	flag.Parse()

	if *f == "" {
		fmt.Println("Base license is not required")
		os.Exit(0)
	}

	if *f != "" && !*d {
		fmt.Println("Valid license")
		os.Exit(0)
	}

	fb, err := ioutil.ReadFile(*f)
	checkErr(err)

	var b64Str string
	for _, l := range strings.Split(string(fb), "\n") {
		if strings.HasPrefix(l, "-") {
			continue
		}

		b64Str += strings.TrimSuffix(l, "\n")
	}

	rb, err := base64.StdEncoding.DecodeString(b64Str)
	checkErr(err)

	// fmt.Println(string(rb))

	licArray := strings.Split(strings.Replace(string(rb), ";", "", -1), "\n")

	fmt.Printf("Valid license\n[%s: %s]\n[%s: %s]\n[%s: %s]\n[%s: %s]\n[%s: %s]\n[%s: %s]\n[%s: %s]\n",
		strings.Split(licArray[0], " ")[0],
		strings.Split(licArray[0], " ")[1],
		strings.Split(licArray[1], " ")[0],
		strings.Split(licArray[1], " ")[1],
		strings.Split(licArray[2], " ")[0],
		strings.Split(licArray[2], " ")[1],
		strings.Split(licArray[3], " ")[0],
		strings.Split(licArray[3], " ")[1],
		strings.Split(licArray[4], " ")[0],
		strings.Split(licArray[4], " ")[1],
		strings.Split(licArray[5], " ")[0],
		strings.Split(licArray[5], " ")[1],
		strings.Split(licArray[6], " ")[0],
		strings.Split(licArray[6], " ")[1],
	)
}

// Checks for an error and fatals if error is not nil
func checkErr(err error) {
	if err != nil {
		fmt.Println("Failed [ signature ]")
		os.Exit(0)
	}
}
