package diff

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func harmonizeSlicesSize(sliceA, sliceB []string) ([]string, []string) {
	if len(sliceA) > len(sliceB) {
		for len(sliceB) != len(sliceA) {
			sliceB = append(sliceB, "")
		}
	} else {
		for len(sliceA) != len(sliceB) {
			sliceA = append(sliceA, "")
		}
	}

	return sliceA, sliceB
}

func readLines(path string) (lines map[string]string, lineChecksums []string, fileChecksum string) {
	lines = make(map[string]string)

	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	fileChecksum = generateFileChecksum(bytes.NewReader(fileData))

	scanner := bufio.NewScanner(bytes.NewReader(fileData))
	for scanner.Scan() {
		checksum := generateStringChecksum(scanner.Text())
		lines[checksum] = scanner.Text()
		lineChecksums = append(lineChecksums, checksum)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}

	return
}

func generateStringChecksum(input string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(input)))
}

func generateFileChecksum(reader io.Reader) string {
	hash := md5.New()
	if _, err := io.Copy(hash, reader); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}
