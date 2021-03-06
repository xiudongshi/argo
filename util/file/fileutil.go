package file

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// IsFileOrDirExistInGZip return true if file or directory exists in GZip file
func IsFileOrDirExistInGZip(sourcePath string, gzipFilePath string) bool {

	fi, err := os.Open(gzipFilePath)

	if os.IsNotExist(err) {
		return false
	}
	defer close(fi)

	fz, err := gzip.NewReader(fi)
	if err != nil {
		return false
	}
	tr := tar.NewReader(fz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {

			return false
		}
		if hdr.FileInfo().IsDir() && strings.Contains(strings.Trim(hdr.Name, "/"), strings.Trim(sourcePath, "/")) {
			return true
		}
		if strings.Contains(sourcePath, hdr.Name) && hdr.Size > 0 {
			return true
		}
	}
	return false
}

//Close the file
func close(f io.Closer) {
	err := f.Close()
	if err != nil {
		log.Warnf("Failed to close the file/writer/reader. %v", err)
	}
}

// CompressEncodeString will return the compressed string with base64 encoded
func CompressEncodeString(content string) string {
	return base64.StdEncoding.EncodeToString(CompressContent([]byte(content)))
}

// DecodeDecompressString will return  decode and decompress the
func DecodeDecompressString(content string) (string, error) {

	buf, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", err
	}
	dBuf, err := DecompressContent(buf)
	if err != nil {
		return "", err
	}
	return string(dBuf), nil
}

// CompressContent will compress the byte array using zip writer
func CompressContent(content []byte) []byte {
	var buf bytes.Buffer
	zipWriter := gzip.NewWriter(&buf)

	_, err := zipWriter.Write(content)
	if err != nil {
		log.Warnf("Error in compressing: %v", err)
	}
	close(zipWriter)
	return buf.Bytes()
}

// DecompressContent will return the uncompressed content
func DecompressContent(content []byte) ([]byte, error) {

	buf := bytes.NewReader(content)
	gZipReader, _ := gzip.NewReader(buf)
	defer close(gZipReader)
	return ioutil.ReadAll(gZipReader)
}
