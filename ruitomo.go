package main

import (
	"encoding/gob"
	"fmt"
	"github.com/banthar/Go-SDL/mixer"
	"io/ioutil"
	"krkrparser"
	"os"
	"strings"
	"time"
)

type sel struct {
	text   string
	target string
}

type config struct {
	File      string
	Sceneid   string
	Scene     string
	SceneEnd  string
	MsgSwitch bool
	Line      int

	VocCache string

	Bgm   string
	Text  string
	Voice string
	Se    string

	Bgmvol int
	Sevol  int
	Vovol  int

	sellst []sel
}

func NewConfig() *config {
	return new(config)
}

var (
	con     = NewConfig()
	bgm     *mixer.Music
	wav, se *mixer.Chunk
	p       *krkrparser.Parser
)

func load() {
	if bgm != nil {
		mixer.FadeOutMusic(100)
		bgm = nil
	}
	if se != nil {
		se.Free()
		se = nil
	}
	line := con.Line
	p = loadsc()
		for n, s, err := p.Next();con.Line < line && (err == nil||len(s)>0); n, s, err = p.Next() {
		switch n {
		case krkrparser.COMMAND:
			execute(strings.TrimSpace(s))
		case krkrparser.TEXT:
			fmt.Print("--", s)
		}
		con.Line++
	}

}

func loadsc() *krkrparser.Parser {
	f, err := os.Open(con.Text + "/" + con.File + ".txt")
	if err != nil {
		return nil
	}
	defer f.Close()
	p := krkrparser.NewParser()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil
	}
	p.Init(b)
	con.Line = 0
	con.SceneEnd = ""
	con.Sceneid = ""
	con.Scene = ""
	return p
}

func playVo() {
	var subdir string
	switch strings.ToUpper(strings.Split(con.VocCache, "_")[0]) {
	case "0A":
		subdir = "才野原惠"
	case "AW":
		subdir = "あわせ"
	case "0D":
		subdir = "花城花鶏"
	case "0B":
		subdir = "皆元るい"
	case "0F":
		subdir = "茅場茜子"
	case "0P", "1U":
		subdir = "三宅"
	case "1K":
		subdir = "蝉丸いずる"
	case "0L":
		subdir = "冬篠宮和"
	case "0G":
		subdir = "白鞘伊代"
	case "3A":
		subdir = "浜江"
	case "2Z":
		subdir = "芳川佐知子"
	case "0E":
		subdir = "鳴滝こより"
	case "1N":
		subdir = "鳴滝小夜里"
	case "2N":
		subdir = "和久津真耶"
	case "0C":
		subdir = "和久津智"
	case "0R":
		subdir = "尹央輝"
	default:
		subdir = "モブ"
	}
	if subdir != "" {
		if wav != nil {
			wav.Free()
		}
		wav = mixer.LoadWAV(con.Voice + "/" + subdir + "/" + con.VocCache + ".ogg")
		if wav == nil {
			println("can not load voice:\t",con.VocCache)
			return
		}
		wav.Volume(con.Vovol)
		wav.PlayChannel(-1, 0)
	}
}

func parseSel(s string) sel {
	text, target := "", ""
	for i := 0; i < len(s); {
		w := make([]byte, 0)
		for ; i < len(s) && s[i] != '\n' && s[i] != ' ' && s[i] != '\t' && s[i] != '='; i++ {
			w = append(w, s[i])
		}
		i++
		switch strings.ToUpper(strings.TrimSpace(string(w))) {
		case "TEXT":
			w = make([]byte, 0)
			for ; i < len(s) && s[i] != '\n' && s[i] != ' ' && s[i] != '\t'; i++ {
				w = append(w, s[i])
			}
			i++
			text = string(w)
		case "TARGET":
			w = make([]byte, 0)
			for ; i < len(s) && s[i] != '\n' && s[i] != ' ' && s[i] != '\t'; i++ {
				w = append(w, s[i])
			}
			i++
			target = string(w)
		default:

		}
	}

	return sel{text, target}
}

func execute(c string) {
	ocommand := c[1 : len(c)-1]
	c = strings.TrimSpace(c)
	argument := strings.Split(c[1:len(c)-1], " ")
	c = argument[0]
	switch {
	case strings.Contains(c, "'"):
		arr := strings.Split(c, "'")
		if len(arr) == 2 {
			fmt.Printf("%s(%s)", arr[0], arr[1])
		}
	case strings.HasPrefix(strings.ToUpper(c), "TWAIT"):
		t := 0
		fmt.Sscanf(argument[1], "time=%d", &t)
		if t > 0 {
			time.Sleep(time.Duration(t) * time.Millisecond)
		}
	case strings.HasPrefix(strings.ToUpper(c), "CH"):
		rbs := ""
		fmt.Sscanf(argument[1], "text=%s", &rbs)
		fmt.Print(strings.Replace(rbs, "\"", "", -1))
		for j := 2; j < len(argument); j++ {
			fmt.Print(argument[j])
		}
	case strings.HasPrefix(strings.ToUpper(c), "SELADD"):
		if con.sellst == nil {
			con.sellst = []sel{parseSel(ocommand)}
		} else {
			con.sellst = append(con.sellst, parseSel(ocommand))
		}
	case strings.HasPrefix(strings.ToUpper(c), "SELECT"):
		if con.sellst == nil {
			fmt.Println("unexpected select command")
		}
		fmt.Println("select:")
		j := -1
		for i := 0; i < len(con.sellst); i++ {
			fmt.Printf("%d. %s\n", i, con.sellst[i].text)
		}
		for j < 0 || j >= len(con.sellst) {
			fmt.Scanf("%d", &j)
		}
		jumpto(con.sellst[j].target)
		con.sellst = nil
	case strings.HasPrefix(strings.ToUpper(c), "BGM"):
		switch len(argument) {
		case 1:
			if bgm != nil {
				bgm.Free()
			}
			bgm = mixer.LoadMUS(con.Bgm + "/" + c + ".ogg")
			if bgm == nil{
				println("can not load bgm:\t",c)
				return
			}
			bgm.PlayMusic(-1)
		case 2:
			t := 0
			_, err := fmt.Sscanf(argument[1], "STOP=%d", &t)
			if err == nil {
				mixer.FadeOutMusic(int(t))
				bgm = nil
				return
			}
		default:
		}

	case strings.HasPrefix(strings.ToUpper(c), "SE"):
		if se != nil {
			se.Free()
		}
		se = mixer.LoadWAV(con.Se + "/" + c + ".ogg")
		if se == nil {
			return
		}
		se.Volume(con.Sevol)
		se.PlayChannel(0, 0)
	case strings.HasPrefix(strings.ToUpper(c), "MSGOFF"):
		con.MsgSwitch = false
	case strings.HasPrefix(strings.ToUpper(c), "BG"):

	case strings.HasPrefix(strings.ToUpper(c), "NEXT"):
		switch {
		case strings.HasPrefix(strings.ToUpper(argument[1]), "STORAGE"):
			fmt.Sscanf(argument[1], "storage=%s", &con.File)
			con.File = strings.Trim(con.File, "\"")
			switch len(argument) {
			case 2:
				sc := loadsc()
				if sc != nil {
					p = sc
				}
			default:
				fmt.Printf("jump to script %s ? (yn)", con.File)
				k := ""
				for strings.ToUpper(strings.TrimSpace(k)) != "Y" && strings.ToUpper(strings.TrimSpace(k)) != "N" {
					fmt.Scanln(&k)
				}
				if strings.ToUpper(strings.TrimSpace(k)) == "Y" {
					sc := loadsc()
					if sc != nil {
						p = sc
					}
				}
			}

		case strings.HasPrefix(strings.ToUpper(argument[1]), "TARGET"):
			tar := ""
			fmt.Sscanf(argument[1], "target=%s", &tar)
			switch len(argument) {
			case 2:
				jumpto(tar)
			default:
				fmt.Printf("jump to tag %s ? (yn)", tar)
				k := ""
				for strings.ToUpper(strings.TrimSpace(k)) != "Y" && strings.ToUpper(strings.TrimSpace(k)) != "N" {
					fmt.Scanln(&k)
					fmt.Println(k)
				}
				if strings.ToUpper(strings.TrimSpace(k)) == "Y" {
					jumpto(tar)
				}
			}
		}
	case len(strings.Split(c, "_")) == 2 && len(strings.Split(c, "_")[0]) == 2:
		con.VocCache = c
		playVo()
	case strings.HasPrefix(strings.ToUpper(c), "BULKSKIP"):
		fmt.Sscanf(argument[1], "target=%s", &con.SceneEnd)

	default:

	}
}

func parseScene(s string) {
	arr := strings.Split(strings.TrimSpace(s), "|")
	switch len(arr) {
	case 2:
		con.Sceneid = arr[0]
		con.Scene = arr[1]
		con.SceneEnd = ""
	case 1:
		con.Sceneid = arr[0]
		con.Scene = ""
		con.SceneEnd = ""
	default:
	}
}

func jumpto(sig string) {
	sig = strings.ToUpper(strings.TrimSpace(sig))
	if strings.HasPrefix(sig, "*") {
		con.VocCache = ""
		for n, s, err := p.Next();strings.ToUpper(strings.TrimSpace(con.Sceneid)) != sig&& (err == nil||len(s)>0); n, s, err = p.Next() {
			switch n {
			case krkrparser.SCENE:
				parseScene(s)
				if strings.ToUpper(strings.TrimSpace(con.Sceneid)) == sig {
					break
				}
			case krkrparser.COMMAND:
			default:
			}
			con.Line++
		}
	}
}

func main(){
fmt.Println("terminal-loader public domain.")
	fmt.Print("checking audio device...")
	if mixer.OpenAudio(mixer.DEFAULT_FREQUENCY, mixer.DEFAULT_FORMAT,
		mixer.DEFAULT_CHANNELS, 4096) != 0 {
		fmt.Println("err")
		return
	}
	defer mixer.HaltMusic()
	defer mixer.CloseAudio()
	fmt.Println("ok")

	fmt.Print("build environment...")
	con.MsgSwitch = true
	con.Line = 0
	con.Bgm = "./data/bgm/"
	con.Se = "./data/sound/"
	con.Text = "./data/scenario/"
	con.Voice = "./data/voice/"
	con.File = "start.ks"
	con.Bgmvol = 20
	con.Sevol = 30
	con.Vovol =	30
	mixer.VolumeMusic(30)
	fmt.Println("ok")

	fmt.Print("test sound...")
	music := mixer.LoadMUS("./start2.wav")
	if music == nil {
		fmt.Println("err")
		return
	}
	sound := mixer.LoadWAV("./start.wav")
	if sound == nil {
		fmt.Println("err")
		return
	}
	sound.Volume(con.Vovol)
	music.PlayMusic(1)
	sound.PlayChannel(0, 0)
	for mixer.PlayingMusic() == 1 {
		time.Sleep(1000)
	}
	sound.Free()
	music.Free()
	fmt.Println("ok")

	fmt.Print("load start script...")
	p = loadsc()
	if p == nil {
		fmt.Println("err")
		return
	}
	fmt.Println("ok")

	var command string
	var stop bool

	for {
		stop = false
		con.VocCache = ""
		for n, s, err := p.Next(); n != krkrparser.SPACE && (err == nil||len(s)>0); n, s, err = p.Next() {
			switch n {
			case krkrparser.COMMA:
			case krkrparser.COMMAND:
				execute(strings.TrimSpace(s))
			case krkrparser.TEXT:
				if !con.MsgSwitch {
					s = ";-" + s
				}
				fmt.Print(s)
				stop = true
			case krkrparser.SCENE:
				parseScene(s)
			default:
			}
			con.Line++
		}
		con.MsgSwitch = true

		for stop {
			fmt.Scanf("%s", &command)

			switch strings.ToUpper(strings.TrimSpace(command)) {
			case "H", "HELP":
				fmt.Println(`welcome to TerminalLoader!
				This is a project for loading Virtual Novel (krkr scripts)
				command list:(case insensitive)
				q/quit: quit program
				n/next: next COMMAND
				s/save: save state
				l/load: load state
				r/repeat: repeat Voice
				v+,v-: voice volume
				s+,s-: se volume
				b+,b-: bgm volume
				v/view: view state
				h/help: show this text`)
			case "Q", "QUIT":
				return
			case "O", "OPEN":

			case "N", "NEXT", "":
				stop = false
			case "S", "SAVE":
				f, err := os.Create("TerminalSave.dat")
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				gob.NewEncoder(f).Encode(*con)
				f.Close()
			case "L", "LOAD":
				f, err := os.Open("TerminalSave.dat")
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				gob.NewDecoder(f).Decode(con)
				f.Close()
				load()
				stop = false
			case "J", "JUMP":
				fmt.Printf("jumping to %s ...", con.SceneEnd)
				jumpto(con.SceneEnd)
				fmt.Println("ok")
				stop = false
			case "R", "REPEAT":
				playVo()
			case "V+":
				con.Vovol += 10
			case "V-":
				con.Vovol -= 10
			case "S+":
				con.Sevol += 10
			case "S-":
				con.Sevol -= 10
			case "B+":
				con.Bgmvol += 10
			case "B-":
				con.Bgmvol -= 10
			case "V", "VIEW":
				fmt.Println(*con)
			default:
			}
		}

	}
}
