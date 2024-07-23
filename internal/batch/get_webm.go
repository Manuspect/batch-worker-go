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

func GetWebmCsv(pathGzipStream string) ([]FileCsv, error) {

	gzipStream, err := os.Open(pathGzipStream)
	if err != nil {
		fmt.Println("error")
	}
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return nil, err
	}

	tarReader := tar.NewReader(uncompressedStream)

	HeaderWebm := ""
	PathCsv := ""
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
			if err := os.Mkdir(header.Name, 0755); err != nil {
				return nil, err
			}
		case tar.TypeReg:
			outFile, err := os.Create(header.Name)
			if err != nil {
				return nil, err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return nil, err
			}
			outFile.Close()

		default:
			log.Printf("ExtractTarGz: uknown type: %v in %s", header.Typeflag, header.Name)
		}

		arr := strings.Split(fmt.Sprintf("%v", header.Name), ".")

		if len(arr) > 1 && arr[1] == "webm" {
			HeaderWebm = header.Name
		}

		if len(arr) > 1 && arr[1] == "csv" {
			PathCsv = header.Name
		}
	}

	frames := "frames"
	err = os.MkdirAll(frames, 0755)
	if err != nil {
		return nil, err
	}

	ar := strings.Split(fmt.Sprintf("%v", HeaderWebm), "/")
	saveWebmToFrames(frames, HeaderWebm)

	events, err := GetEvetsCsv(PathCsv)
	if err != nil {
		return nil, err
	}

	err = os.RemoveAll(ar[0])
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
