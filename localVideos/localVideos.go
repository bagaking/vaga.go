package localVideos

import (
	"fmt"
	. "github.com/c2h5oh/datasize"
	"io/ioutil"
	"strings"
	"path"
)

type Video struct {
	NameHash string `json:"hash"`
	Location string `json:"location"`

	DirName   string `json:"dir"`
	FileName  string `json:"name"`
	FileSize  int64  `json:"size"`
	FileSizeH string `json:"size_h"`
}

type VideoBlob struct {
	Name        string `json:"name"`
	Description string `json:"desc"`

	RootPath string  `json:"root"`
	Videos   []Video `json:"videos"`
}

//var Videos = make([]Video, 0, 10000)
var VideosMap = map[string]Video{}

func (videoBlob *VideoBlob) Initial() error {
	return videoBlob.scanVideos(videoBlob.RootPath)
}

func (videoBlob *VideoBlob) scanVideos(baseName string) error {
	rd, err := ioutil.ReadDir(baseName)
	count := 0
	for _, fi := range rd {
		fileName := fi.Name()

		pathFile := path.Join(baseName, fileName)
		if fi.IsDir() {
			pathDir := pathFile
			videoBlob.scanVideos(pathDir)
		} else {
			if strings.HasSuffix(fileName, ".mp4") {
				//fmt.Println("==", len(videoBlob.Videos), pathFile)
				metaVideo := Video{
					NameHash: Sha1([]byte(pathFile)),
					DirName:  baseName,
					Location: pathFile,

					FileName: fileName,
					FileSize: fi.Size(),
				}
				metaVideo.FileSizeH = ByteSize(metaVideo.FileSize).HR()
				videoBlob.Videos = append(videoBlob.Videos, metaVideo)
				VideosMap[metaVideo.NameHash] = metaVideo
				count++
			} else {
				//fmt.Printf("== x %s\n", pathFile)
			}
		}
	}
	if count > 0 {
		fmt.Printf("[%s, %p] Scan on %s: got %d, total %d\n", videoBlob.Name, &videoBlob, baseName, count, len(videoBlob.Videos))
	}
	return err
}
