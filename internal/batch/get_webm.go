package batch

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func GetWebmJson(pathStream string) ([]FileJson, error) {

	gzipStream, err := os.Open(pathStream)
	if err != nil {
		fmt.Println("error")
	}
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return nil, err
	}

	tarReader := tar.NewReader(uncompressedStream)

	arrCat := strings.Split(pathStream, "/")

	HeaderWebm := ""
	PathJson := ""
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(arrCat[1]+"/"+header.Name, 0755); err != nil {
				return nil, err
			}

		case tar.TypeReg:
			outFile, err := os.Create(arrCat[1] + "/" + header.Name)
			if err != nil {
				return nil, err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return nil, err
			}
			outFile.Close()

		default:
			log.Printf("GetWebmJson: uknown type: %s in %s", string(header.Typeflag), arrCat[1]+"/"+header.Name)
		}

		arr := strings.Split(header.Name, ".")

		if len(arr) > 1 && arr[1] == "webm" {
			HeaderWebm = arrCat[1] + "/" + header.Name
		}

		if len(arr) > 1 && arr[1] == "json" {
			PathJson = arrCat[1] + "/" + header.Name
		}
	}

	frames := arrCat[1] + "/" + "frames"
	err = os.MkdirAll(frames, 0755)
	if err != nil {
		return nil, err
	}

	saveWebmToFrames(frames, HeaderWebm)

	events, err := GetEventsJson(PathJson)
	if err != nil {
		return nil, err
	}

	return events, err
}

func saveWebmToFrames(frames, path string) {
	c := exec.Command(
		"ffmpeg", "-i",
		path,
		"-r", "40", frames+"/%d.jpeg",
	)
	c.Stderr = os.Stderr
	c.Run()
}
