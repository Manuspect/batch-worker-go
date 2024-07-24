package batch

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type FileCsv struct {
	Timestamp     string `json:"timestamp"`
	Process_path  string `json:"process_path"`
	Title         string `json:"title"`
	Class_name    string `json:"class_name"`
	Window_left   string `json:"window_left"`
	Window_top    string `json:"window_top"`
	Window_right  string `json:"window_right"`
	Window_bottom string `json:"window_bottom"`
	Event         string `json:"event"`
	Mouse_x_pos   string `json:"mouse_x_pos"`
	Mouse_y_pos   string `json:"mouse_y_pos"`
	Modifiers     string `json:"modifiers"`
}

type Info struct {
	Timestamp     string `json:"timestamp"`
	Process_path  string `json:"process_path"`
	Title         string `json:"title"`
	Class_name    string `json:"class_name"`
	Window_left   string `json:"window_left"`
	Window_top    string `json:"window_top"`
	Window_right  string `json:"window_right"`
	Window_bottom string `json:"window_bottom"`
	Event         string `json:"event"`
	Mouse_x_pos   string `json:"mouse_x_pos"`
	Mouse_y_pos   string `json:"mouse_y_pos"`
	Modifiers     string `json:"modifiers"`
	FilePath      string `json:"filePath"`
}

func GetEvetsCsv(pathCsv string) ([]FileCsv, error) {
	file, err := os.Open(pathCsv)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var events []FileCsv

	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()

		lineArr := strings.Split(line, ",")

		if len(lineArr) != 12 {
			return nil, fmt.Errorf("сan't parse string: %s", line)
		}

		events = append(events, FileCsv{
			Timestamp:     lineArr[0],
			Process_path:  lineArr[1],
			Title:         lineArr[2],
			Class_name:    lineArr[3],
			Window_left:   lineArr[4],
			Window_top:    lineArr[5],
			Window_right:  lineArr[6],
			Window_bottom: lineArr[7],
			Event:         lineArr[8],
			Mouse_x_pos:   lineArr[9],
			Mouse_y_pos:   lineArr[10],
			Modifiers:     lineArr[11],
		})

	}
	return events, err
}
