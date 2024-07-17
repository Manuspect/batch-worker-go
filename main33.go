// main4.go
package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	// Specify the destination directory
	dst := "unzip"

	// Open the zip file
	fmt.Println("open zip archive...")
	archive, err := zip.OpenReader("/home/user/Видео/batch-1719407918236.zip")
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	// Extract the files from the zip
	for _, f := range archive.File {

		// Create the destination file path
		filePath := filepath.Join(dst, f.Name)
		// Print the file path
		fmt.Println("extracting file ", filePath)

		// Check if the file is a directory
		if f.FileInfo().IsDir() {
			// Create the directory
			fmt.Println("creating directory...")
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				panic(err)
			}
			continue
		}

		// Create the parent directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		// Create an empty destination file
		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		// Open the file in the zip and copy its contents to the destination file
		srcFile, err := f.Open()
		if err != nil {
			panic(err)
		}
		if _, err := io.Copy(dstFile, srcFile); err != nil {
			panic(err)
		}

		// Close the files
		dstFile.Close()
		srcFile.Close()
	}
	// inPath := "unzip/batch-1719407918236/output.webm"
	// outDir := "frames/"
	// ExtractImages(inPath, outDir)

}

// func ExtractImages(inPath string, outDir string) error {
// 	imgFormat := "png"
// 	opts := ffmpeg.Options{
// 		OutputFormat: &imgFormat,
// 	}

// 	err := os.Mkdir(outDir, os.ModeDir|fs.ModePerm)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = GetTranscoder(inPath).
// 		WithOptions(opts).
// 		Output(path.Join(outDir, "%6d.png")).
// 		Start(opts)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// ffmpeg -i input.mp4 -an -s qvga %06d.png
// ffmpeg -i video.VOB -vsync 0 -ss 01:30 -to 01:40 %06d.bmp
// ffmpeg -i input.mp4 -vf "fps=1" frame%04d.png
// ffmpeg -i input.webm -c:v libx264 -c:a aac output.mp4
