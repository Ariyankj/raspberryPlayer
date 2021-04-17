package main

import (
	"encoding/json"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var stream beep.StreamSeeker
var Filename string
var lock bool

func StartHandler(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintf(w,"play");
	fmt.Println("play")
	fmt.Println(r.FormValue("song"))
	f, _ := os.Open("./downloaded/"+r.FormValue("song")+".mp3")
	stream, _, _ = mp3.Decode(f)
	_, format, _ := mp3.Decode(f)
	speaker.Clear()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if lock {
		speaker.Unlock()
	}
	speaker.Play(stream)
}
func pauseHandler(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("pause")
	fmt.Fprintf(w,"pause")
	lock=true
	speaker.Lock()
}
func seekHandler(w http.ResponseWriter, r *http.Request)  {
	a:=r.FormValue("time")
	lock=true
	speaker.Lock()
	b,_ :=strconv.Atoi(a)
	stream.Seek(b)
	fmt.Println(stream.Position())
	lock=false
	speaker.Unlock()

}
func play(w http.ResponseWriter, r *http.Request)  {
	lock=false
	speaker.Unlock()

}
func online(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintf(w,"online");
	fmt.Println("online")
	f, _ := os.Open("./1.mp3")
	stream, _, _ = mp3.Decode(f)
	_, format, _ := mp3.Decode(f)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if lock {
		speaker.Unlock()
	}
	speaker.Play(stream)

}
func upload(w http.ResponseWriter, r *http.Request)  {
	print("here")
	file, handler, err := r.FormFile("uploaded_file")
	defer file.Close()
	if (err!=nil) {
		fmt.Println(err)
	}
	filename :=handler.Filename
	if(strings.Contains(filename,".MP3")) {
		filename=strings.Split(filename,".MP3")[0]+".mp3"

	}else if(strings.Contains(filename,".Mp3")){
		filename=strings.Split(filename,".Mp3")[0]+".mp3"
	}else if(strings.Contains(filename,".mP3")){
		filename=strings.Split(filename,".mP3")[0]+".mp3"
	}

	f, err := os.OpenFile("./downloaded/"+filename, os.O_WRONLY|os.O_CREATE, 0666)
	if(err!=nil){
		fmt.Println(err)
	}
	defer f.Close()
	io.Copy(f, file)
	Filename=handler.Filename
	fmt.Fprintf(w,"done")
	fmt.Println("done")
}
func GetSongList(w http.ResponseWriter, r *http.Request)  {
	list := []string{}
	files, err := ioutil.ReadDir("./downloaded")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.Contains(f.Name() ,".mp3"){
			list=append(list,f.Name())
		}
	}
	encjson, _ := json.Marshal(list)
	fmt.Println("sent")
	fmt.Fprint(w,string(encjson))
}
func deletesong(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintf(w,"delete");
	fmt.Println("play")
	fmt.Println(r.FormValue("song"))
	a:=os.Remove("./downloaded/"+r.FormValue("song")+".mp3")
	if a!=nil {
		fmt.Println(a)
	}
}
func main() {
	http.HandleFunc("/seek", seekHandler)
	http.HandleFunc("/pause", pauseHandler)
	http.HandleFunc("/Start", StartHandler)
	http.HandleFunc("/play",play)
	http.HandleFunc("/online",online)
	http.HandleFunc("/upload",upload)
	http.HandleFunc("/delete",deletesong)
	http.HandleFunc("/GetSongList",GetSongList)
	http.ListenAndServe(":3000", nil)
}