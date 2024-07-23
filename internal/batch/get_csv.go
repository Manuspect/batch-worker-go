package batch

import (
	"bufio"
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
	FrameName     string `json:"frameName"`
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

		lineArr2 := strings.Split(line, ",")
		var lineArr = [12]string{}

		if len(lineArr2) == 12 {
			lineArr[0] = lineArr2[0]
			lineArr[1] = lineArr2[1]
			lineArr[2] = lineArr2[2]
			lineArr[3] = lineArr2[3]
			lineArr[4] = lineArr2[4]
			lineArr[5] = lineArr2[5]
			lineArr[6] = lineArr2[6]
			lineArr[7] = lineArr2[7]
			lineArr[8] = lineArr2[8]
			lineArr[9] = lineArr2[9]
			lineArr[10] = lineArr2[10]
			lineArr[11] = lineArr2[11]
		}

		if len(lineArr2) == 13 {
			lineArr[0] = lineArr2[0]
			lineArr[1] = lineArr2[1]
			lineArr[2] = lineArr2[2]
			lineArr[3] = lineArr2[3]
			lineArr[4] = lineArr2[4]
			lineArr[5] = lineArr2[5]
			lineArr[6] = lineArr2[6]
			lineArr[7] = lineArr2[7]
			lineArr[8] = lineArr2[8]
			lineArr[9] = lineArr2[9]
			lineArr[10] = lineArr2[10]
			lineArr[11] = lineArr2[11] + ", " + lineArr2[12]
		}

		events = append(events, FileCsv{
			Timestamp:     lineArr[0] + "|",
			Process_path:  lineArr[1] + "|",
			Title:         lineArr[2] + "|",
			Class_name:    lineArr[3] + "|",
			Window_left:   lineArr[4] + "|",
			Window_top:    lineArr[5] + "|",
			Window_right:  lineArr[6] + "|",
			Window_bottom: lineArr[7] + "|",
			Event:         lineArr[8] + "|",
			Mouse_x_pos:   lineArr[9] + "|",
			Mouse_y_pos:   lineArr[10] + "|",
			Modifiers:     lineArr[11],
		})

	}
	return events, err
}
