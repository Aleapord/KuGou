package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

type Music struct {
	Name      string
	Author    string
	Album     string
	MusicHash string
	MusicId   string
	FileName  string
}

type KuGou struct {
	Path    string
	Musics  []Music
	Keyword string
}

func (h KuGou) Init() {
	_ = os.Mkdir(h.Path, 0755)
}

func getBody(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return body
}

func (h *KuGou) Search() {
	var json_obj = make(map[string]interface{})
	unicode := url.Values{}
	compiler, _ := regexp.Compile("\\(.+\\)")
	unicode.Set("keyword", h.Keyword)
	my_url := "https://songsearch.kugou.com/song_search_v2?callback=jQuery1124042761514747027074_1580194546707&" + unicode.Encode() + "&page=1&pagesize=30&userid=-1&clientver=&platform=WebFilter&tag=em&filter=2&iscorrection=1&privilege_filter=0&_=1580194546709"
	s := compiler.FindString(string(getBody(my_url)))
	s = s[1 : len(s)-1]
	err := json.Unmarshal([]byte(s), &json_obj)
	if err != nil {
		panic(err)
	}

	s2 := json_obj["data"]
	for _, i3 := range s2.(map[string]interface{})["lists"].([]interface{}) {
		i3 := i3.(map[string]interface{})
		var music = new(Music)
		music.Album = i3["AlbumName"].(string)
		music.Author = i3["SingerName"].(string)
		compiler, _ = regexp.Compile("<em>|</em>")
		name := compiler.ReplaceAllString(i3["FileName"].(string), "")
		music.MusicHash = i3["FileHash"].(string)
		music.MusicId = i3["AlbumID"].(string)
		music.FileName = name
		h.Musics = append(h.Musics, *music)
	}
}

func (h *KuGou) Downolad(index int) {
	json_url := "https://wwwapi.kugou.com/yy/index.php?r=play/getdata&hash=" +
		h.Musics[index].MusicHash + "&album_id=" +
		h.Musics[index].MusicId +
		"&dfid=2SSV0x4LWcsx0iylej1F6w7P&mid=44328d3dc4bfce21cf2b95cf9e76b968&platid=4"
	j := getBody(json_url)
	var my_json = make(map[string]interface{})
	_ = json.Unmarshal(j, &my_json)
	for s2, i := range my_json["data"].(map[string]interface{}) {
		if s2 == "play_url" {
			json_url = i.(string)
		}
	}

	bs := getBody(json_url)
	_ = ioutil.WriteFile(h.Path+"/"+h.Musics[index].FileName+".mp3", bs, 0755)
}
func main() {
	var key string
	fmt.Println("输入歌曲名:")
	_, _ = fmt.Scan(&key)
	var h = KuGou{
		Keyword: key,
		Path:    "./music",
		Musics:  make([]Music, 0),
	}
	h.Init()
	h.Search()
	for i, s2 := range h.Musics {
		fmt.Println(i+1, s2.FileName)
	}

	fmt.Println("输入要下载的序列号：")
	var choic int
	_, _ = fmt.Scan(&choic)

	h.Downolad(choic - 1)

}
