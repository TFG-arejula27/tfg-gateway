package receptor

import (
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type PowerInfo struct {
	Message string        `json:"message"`
	Averge  PowerInfoData `json:"averge"`
	Max     PowerInfoData `json:"max"`
	Min     PowerInfoData `json:"min"`
	C1      CStateData    `json:"c1"`
	C2      CStateData    `json:"c2"`
	Poll    CStateData    `json:"poll"`
	C0      CStateData    `json:"c0"`
	frames  []PowerInfoData
	Time    time.Duration
}
type PowerInfoData struct {
	Power     string
	Frecuency string
}
type CStateData struct {
	Resident string
	Count    string
	Latency  string
}

//Run executes powerstate
func runPowerstat() (PowerInfo, error) {
	cmd := exec.Command("powerstat", "-R", "-n", "1", "60")
	pwrInf := PowerInfo{frames: make([]PowerInfoData, 0)}
	data, err := cmd.Output()
	if err != nil {
		return pwrInf, err

	}
	output := string(data)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.Join(strings.Fields(line), " ")

		if strings.HasPrefix(line, "Average") {
			frameInfoData, err := getFrameInfoData(line)
			if err != nil {
				log.Println("cant get average info")
				continue
			}
			pwrInf.Averge = frameInfoData

		} else {
			continue
		}

	}
	//parsedLines = append(parsedLines, strings.Split(line, " "))

	return pwrInf, nil

}

func getFrameInfoData(line string) (PowerInfoData, error) {
	parsedLine := strings.Split(line, " ")

	return PowerInfoData{
		Power:     parsedLine[12],
		Frecuency: parsedLine[13],
	}, nil
}

func getCstateData(line string) (CStateData, error) {
	var cstatedata = CStateData{}
	line = strings.Replace(line, "%", "", 1)
	parsedLine := strings.Split(line, " ")
	resident := parsedLine[1]
	cstatedata.Resident = resident
	if parsedLine[0] == "C0" {
		return cstatedata, nil
	}

	cstatedata.Count = parsedLine[2]
	cstatedata.Latency = parsedLine[3]

	return cstatedata, nil

}

func isFrame(line string) bool {
	res := regexp.MustCompile(`^[0-9\.].*$`).MatchString(line)
	return res

}
func framesEnded(line string) bool {
	res := regexp.MustCompile(`^[ ]*Average`).MatchString(line)
	return res

}
func isEmptyLine(line string) bool {
	return line == ""
}
