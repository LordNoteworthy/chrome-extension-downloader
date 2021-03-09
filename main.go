package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

const (
	chromeWebStoreURL = "https://clients2.google.com/service/update2/crx?response=redirect"
	reExtendionID     = `[a-z]{32}`
)

var (
	errHTTPStatusCodeNotOK = errors.New("http download faile, status code != 200")
)

func readInput() ([]string, error) {

	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		lines = append(lines, line)
	}

	return lines, nil
}

// WriteBytesFile write Bytes to a File.
func WriteBytesFile(filename string, r io.Reader) (int, error) {

	// Open a new file for writing only
	file, err := os.OpenFile(
		filename,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Read for the reader bytes to file
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, err
	}

	// Write bytes to disk
	bytesWritten, err := file.Write(b)
	if err != nil {
		return 0, err
	}

	return bytesWritten, nil
}

func downloadFile(url string) ([]byte, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errHTTPStatusCodeNotOK
	}

	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func crx2zip(buf []byte) ([]byte, error) {

	// A CRX₃ file is a binary file of the following format:
	// [4 octets]: "Cr24", a magic number.
	// [4 octets]: The version of the *.crx file format used (currently 3).
	// [4 octets]: N, little-endian, the length of the header section.
	// [N octets]: The header (the binary encoding of a CrxFileHeader).
	// [M octets]: The ZIP archive.

	// 43 72 32 34 (Cr24)
	if buf[0] != 67 || buf[1] != 114 || buf[2] != 50 || buf[3] != 52 {
		return nil, errors.New("invalid header does not start with Cr24")
	}

	// 03 00 00 00
	if buf[4] != 3 || buf[5] != 0 || buf[6] != 0 || buf[7] != 0 {
		return nil, errors.New("unexpected crx format version number")
	}

	// header size
	header := binary.LittleEndian.Uint32(buf[8:])

	// Magic number (4), CRX format version (4), lengths (4)
	zipStartOffset := 4 + 4 + 4 + header
	return buf[zipStartOffset:], nil
}

func createDownloadURL(extensionList []string) {
	r := regexp.MustCompile(reExtendionID)
	for _, ext := range extensionList {

		log.Printf("Processing %s\n", ext)

		match := r.FindStringSubmatch(ext)
		if len(match) == 0 {
			continue
		}
		extensionID := match[0]

		// Omitting this value is allowed, but add it just in case.
		// Source: http://cs.chromium.org/file:omaha_query_params.
		// cc%20GetProdIdString
		productID := "chromiumcrx"

		// Channel is "unknown" on Chromium on ArchLinux, so using "unknown"
		// will probably be fine for everyone.
		productChannel := "unknown"

		// As of July, the Chrome Web Store sends 204 responses to user agents
		// when their Chrome/Chromium version is older than version 31.0.1609.0
		// so forcing a version in the future.
		productVersion := "9999.0.9999.0"

		// Different OS options available:  win, linux, mac, android, openbsd,
		// cros.
		os := "win"

		// Different architecture availables: x86-32, x86-64, arm.
		// Same goes for NaCl architecture,
		arch := "x86-64"

		url := chromeWebStoreURL
		url += "&os=" + os
		url += "&arch=" + arch
		url += "&nacl_arch=" + arch
		url += "&prod=" + productID
		url += "&prodchannel=" + productChannel
		url += "&prodversion=" + productVersion
		url += "&acceptformat=" + "crx2,crx3"
		url += "&x=id%3D" + extensionID
		url += "%26uc"

		b, err := downloadFile(url)
		if err != nil {
			continue
		}
		WriteBytesFile(extensionID+".crx", bytes.NewReader(b))

		// Chrome Extensions (CRX) are ZIP-files with an added header in the
		// form of magic number + version number + public key length +
		// signature length + public key + signature
		zipBytes, err := crx2zip(b)
		if err != nil {
			continue
		}

		WriteBytesFile(extensionID+".zip", bytes.NewReader(zipBytes))
	}

}

func main() {

	extensionList, err := readInput()
	if err != nil {
		log.Fatalln(err)
	}

	createDownloadURL(extensionList)
}
