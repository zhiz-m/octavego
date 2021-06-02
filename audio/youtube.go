package audio

import (
	"os/exec"
)

func YTDL(song *Song) (string, error) {
	args := []string{
		"-x",
		"--skip-download",
		"--get-url",
		"--audio-quality",
		"128K",
		"ytsearch:" + song.SearchQuery,
	}
	out, err := exec.Command("youtube-dl", args...).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

/*
func YTDL(songs []*Song, YTDLKill <-chan bool, YTDLReady chan<- bool) {
	defer func (){
		YTDLReady <- true
	}()
	args := []string{
		"-x",
		"--skip-download",
		"--get-url",
		"--audio-quality",
		"128K",
	}
	for _, song := range songs {
		args = append(args, "ytsearch:"+song.SearchQuery)
	}
	cmd := exec.Command("youtube-dl", args...)
	out, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println("Error creating pipe for youtube-dl process")
		return
	}

	if err := cmd.Start(); err != nil{
		fmt.Println("Error starting youtube-dl process")
		return
	}

	i := 0
	currentURL := ""
	for {
		select {
		case <-YTDLKill:
			if err := cmd.Process.Kill(); err != nil {
				fmt.Println("Error destroying youtube-dl process")
			}
			return
		default:
		}
		buffer := make([]byte, 1)
		_, err := out.Read(buffer)
		if err == io.EOF {
			return
		}
		current := buffer[0]
		if current == byte('\n') {
			songs[i].SourceURL = currentURL
			songs[i].Load()

			currentURL = ""
		} else {
			currentURL = currentURL + string(current)
		}
	}
}

*/
