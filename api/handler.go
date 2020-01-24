package api

import (
	"bagaking.com/vaga.go/localVideos"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type H interface{}

var AllAvailableVideoBlobs []*localVideos.VideoBlob // deep copy

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
}

/**
 * handler of showing page - index
 */
func HandlerVideoIndex(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	t, _ := template.ParseFiles("./tpl/index.html", "./tpl/basic_style.html")
	_ = t.ExecuteTemplate(writer, "index", AllAvailableVideoBlobs)
}

/**
 * handler of showing video tree of a video blob
 */
func HandlerVideoTree(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	t, _ := template.ParseFiles("./tpl/tree.html", "./tpl/basic_style.html")
	indStr := params.ByName("blob_ind")
	ind, _ := strconv.Atoi(indStr)


	blob := AllAvailableVideoBlobs[ind]
	fmt.Println(blob.Name, &blob)

	resultMap := make(map[string][]localVideos.Video)
	lastDirName := ""
	for _, videoMeta := range blob.Videos {
		if videoMeta.DirName != lastDirName {
			lastDirName = videoMeta.DirName
			resultMap[lastDirName] = make([]localVideos.Video, 0)
		}
		resultMap[lastDirName] = append(resultMap[lastDirName], videoMeta)
	}
	_ = t.ExecuteTemplate(writer, "tree", map[string]interface{}{
		"resultMap": resultMap,
		"blob":      blob,
	})
}

/**
 * handler of a video watching page
 */
func HandlerVideoWatch(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	hash := params.ByName("video_hash")
	t, _ := template.ParseFiles("./tpl/watch.html", "./tpl/basic_style.html")
	_ = t.ExecuteTemplate(writer, "watch", localVideos.VideosMap[hash])
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
