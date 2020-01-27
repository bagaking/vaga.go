package api

import (
	"github.com/bagaking/vaga.go/localVideos"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"strings"
)

type H interface{}

var AllAvailableVideoBlobs []*localVideos.VideoBlob // deep copy
var VideoDirMap map[string][]*localVideos.Video

/**
 * init all video blobs
 */
func initial(allAvailableVideoBlobs []*localVideos.VideoBlob) {
	for _, videoBlob := range allAvailableVideoBlobs {
		if err := videoBlob.Initial(); err != nil {
			fmt.Printf("Initial Blob %s failed", videoBlob.Name)
		}
	}
	AllAvailableVideoBlobs = allAvailableVideoBlobs

	VideoDirMap = make(map[string][]*localVideos.Video)
	for _, blob := range AllAvailableVideoBlobs {
		for dir, pVideos := range blob.DirMap {
			fmt.Println("=", dir, pVideos)
			if nil == VideoDirMap[dir] {
				VideoDirMap[dir] = make([]*localVideos.Video, 0)
			}
			for _, pVideo := range pVideos {
				VideoDirMap[dir] = append(VideoDirMap[dir], pVideo)
			}
		}
	}

}

func executeTemplate(wr io.Writer, name string, data interface{}) error {
	itoaVal := 0
	tpl, _ := template.New(name + ".html").Funcs(template.FuncMap{
		"replace": func(input, from, to string) string {
			return strings.Replace(input, from, to, -1)
		},
		"itoa": func() int {
			itoaVal = itoaVal + 1
			return itoaVal
		},
	}).ParseFiles(
		"./tpl/"+name+".html",
		"./tpl/basic_style.html",
	)
	return tpl.ExecuteTemplate(wr, name, data)
}

/**
 * handler of showing page - index
 */
func HandlerVideoIndex(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	executeTemplate(writer, "index",
		AllAvailableVideoBlobs,
	)
}

/**
 * handler of showing video tree of a video blob
 */
func HandlerVideoTree(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	indStr := params.ByName("blob_ind")
	ind, _ := strconv.Atoi(indStr)

	blob := AllAvailableVideoBlobs[ind]
	fmt.Println(blob.Name, &blob)

	executeTemplate(writer, "tree", map[string]interface{}{
		"VideoDirMap": blob.DirMap,
		"blob":        blob,
	})
}

/**
 * handler of a video watching page
 */
func HandlerVideoWatch(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	hash := params.ByName("video_hash")
	video := localVideos.VideosMap[hash]

	executeTemplate(writer, "watch",
		map[string]interface{}{
			"dir": VideoDirMap[video.DirName],
			"video":  video,
		})
}

/**
 * handler of video stream
 */
func HandlerVideoStream(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	hash := params.ByName("video_hash")
	pathVideo := localVideos.VideosMap[hash]

	video, err := os.Open(pathVideo.Location)
	if err != nil {
		io.WriteString(writer, "Open video failed")
		return
	}
	defer video.Close()

	writer.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(writer, request, video.Name(), time.Now(), video)
}
