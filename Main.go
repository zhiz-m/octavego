package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/zhiz-m/octavego/audio"
)

var (
	Token  string
	Prefix = "a."
)

func init() {
	args := os.Args[1:]
	if len(args) == 0 || len(args) > 2 {
		fmt.Println("Usage: octavego <token> <prefix: default a.>")
		os.Exit(1)
	}
	Token = args[0]
	if len(args) == 2 {
		Prefix = args[1]
	}
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		panic("Error creating Discord session")
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildVoiceStates

	audio.AddAudio(dg, Prefix)

	err = dg.Open()
	if err != nil {
		panic("Error opening Discord session")
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	audio.Cleanup()

	dg.Close()
}

/*
func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		panic("Error creating Discord session")
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildVoiceStates

	err = dg.Open()
	if err != nil {
		panic("Error opening Discord session")
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, Prefix) {
		cmd := m.Content[len(Prefix):]
		switch {
		case strings.HasPrefix(cmd, "hi"):
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("hi %s", m.Author.ID))
		case strings.HasPrefix(cmd, "test"):
			c, err := s.State.Channel(m.ChannelID)

			if err != nil {
				fmt.Println(err)
				return
			}

			g, err := s.State.Guild(c.GuildID)
			if err != nil {
				fmt.Println("error: failed to find guild")
				return
			}

			for _, vs := range g.VoiceStates {
				if vs.UserID == m.Author.ID {
					playSong(s, g.ID, vs.ChannelID)
				}
			}
		default:
			s.ChannelMessageSend(m.ChannelID, HelpMsg)
		}
	}
}

const (
	FrameRate = 48000
	//FrameSize = 3840
	FrameSize = 960
	Channels  = 2
	URL       = "https://r2---sn-8vq54vox50-hxal.googlevideo.com/videoplayback?expire=1622661047&ei=VoO3YJP0O_WAjuMPkcS7mAM&ip=2405%3A6e00%3A3174%3A5200%3Afc8b%3Ae46b%3Aa1e1%3Ab2a8&id=o-AIq7WML3k15m1WcL92_bYxEHambiIgtOeLlUPoaQciiD&itag=140&source=youtube&requiressl=yes&mh=u9&mm=31%2C29&mn=sn-8vq54vox50-hxal%2Csn-hxa7zn7s&ms=au%2Crdu&mv=m&mvi=2&pl=44&initcwndbps=1151250&vprv=1&mime=audio%2Fmp4&ns=-94OGEDpwFjTfmxHf6OK2jEF&gir=yes&clen=5931025&dur=366.434&lmt=1575219701971364&mt=1622639092&fvip=4&keepalive=yes&fexp=24001373%2C24007246&c=WEB&txp=5531432&n=Qt0MtVCtBIbMgZD0NY&sparams=expire%2Cei%2Cip%2Cid%2Citag%2Csource%2Crequiressl%2Cvprv%2Cmime%2Cns%2Cgir%2Cclen%2Cdur%2Clmt&lsparams=mh%2Cmm%2Cmn%2Cms%2Cmv%2Cmvi%2Cpl%2Cinitcwndbps&lsig=AG3C_xAwRAIgE2mPRdK2WSqNqBXIyn3p2bmJhF1KQefnm9_EKlJABXkCIChA_e-ZfAX_9jusqu9JVkXMIv0ILpIlFQ0QT7zL7k4S&sig=AOq0QJ8wRQIgWLCNC7Wd3KZipRpNOvmAANGe9rMzqqgNdoWUItzy6HACIQDc7o1FotI_GVDkLzr5llx4q6Jp0IL2mipx_U7WmMl0OA=="
)

func playSong(s *discordgo.Session, guildID string, channelID string) {
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		println("error joining voice channel")
		return
	}
	vc.Speaking(true)
	defer vc.Speaking(false)

	cmd := exec.Command("ffmpeg", "-i", URL, "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		println("error creating pipe")
		return
	}
	buffer := bufio.NewReaderSize(pipe, 16384)

	inChan := make(chan []int16, 2)
	go encode(inChan, vc)

	err = cmd.Start()
	if err != nil {
		println("error starting ffmpeg proc")
		return
	}
	//fmt.Printf("hi1\n")
	for {
		buf := make([]int16, FrameSize*Channels)
		err = binary.Read(buffer, binary.LittleEndian, &buf)
		//fmt.Printf("hi2\n")
		if err != nil {
			fmt.Println("error", err)
			return
		}
		//fmt.Printf("hi3\n")
		inChan <- buf
	}

}

func encode(in chan []int16, out *discordgo.VoiceConnection) {
	encoder, err := gopus.NewEncoder(FrameRate, Channels, gopus.Audio)
	if err != nil {
		print("error creating pcm encoder")
		return
	}
	for {
		raw := <-in
		encoded, err := encoder.Encode(raw, FrameSize, FrameSize)
		//fmt.Printf("hi4\n")
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
		}
		if err != nil {
			fmt.Println("error", err)
			return
		}
		//fmt.Printf("sending to out.OpusSend\n")
		out.OpusSend <- encoded
	}
}
*/
