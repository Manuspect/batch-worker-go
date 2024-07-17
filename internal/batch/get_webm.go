package batch

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func GetWebmCsv(pathGzipStream string) ([]image.Image, []FileCsv, error) {

	gzipStream, err := os.Open(pathGzipStream)
	if err != nil {
		fmt.Println("error")
	}
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return nil, nil, err
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
			return nil, nil, err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(header.Name, 0755); err != nil {
				return nil, nil, err
			}
		case tar.TypeReg:
			outFile, err := os.Create(header.Name)
			if err != nil {
				return nil, nil, err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return nil, nil, err
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
		return nil, nil, err
	}

	ar := strings.Split(fmt.Sprintf("%v", HeaderWebm), "/")
	webmToFrame(frames, HeaderWebm)

	files, err := os.ReadDir(frames)
	if err != nil {
		return nil, nil, err
	}
	var images []image.Image
	for j, file := range files {
		if !file.IsDir() {
			image, err := loadFrames(fmt.Sprintf("%s/%d.jpeg", frames, j+1))
			if err != nil {
				return nil, nil, err
			}

			images = append(images, image)
			fmt.Println(j+1, fmt.Sprintf("%s/%d.jpeg", frames, j+1))
		}
	}

	events, err := GetEvetsCsv(PathCsv)
	if err != nil {
		return nil, nil, err
	}

	err = os.RemoveAll(ar[0])
	if err != nil {
		return nil, nil, err
	}

	err = os.RemoveAll(frames)
	if err != nil {
		return nil, nil, err
	}

	return images, events, err
}

func webmToFrame(frames, path string) {
	c := exec.Command(
		"ffmpeg", "-i",
		path,
		"-r", "40", frames+"/%d.jpeg",
	)
	c.Stderr = os.Stderr
	c.Run()
}

func loadFrames(path string) (image.Image, error) {
	existingImageFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer existingImageFile.Close()

	existingImageFile.Seek(0, 0)

	loadedImage, err := jpeg.Decode(existingImageFile)
	if err != nil {
		return nil, err
	}

	return loadedImage, nil
}
