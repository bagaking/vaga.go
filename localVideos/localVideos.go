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
	RootPath    string `json:"root"`

	MaxDepth        int `json:"max_depth"`
	ShrinkThreshold int `json:"shrink_threshold"`

	Videos []Video             `json:"videos"`
	DirMap map[string][]*Video `json:"dir_map"`
}

//var Videos = make([]Video, 0, 10000)
var VideosMap = map[string]Video{}

func (videoBlob *VideoBlob) Initial() error {
	if err := videoBlob.scanVideos(videoBlob.RootPath, 0); err != nil {
		return err
	}
	videoBlob.initialDirMap()
	return nil
}

func (videoBlob *VideoBlob) initialDirMap() {
	videoBlob.DirMap = make(map[string][]*Video)
	maxDepth := videoBlob.MaxDepth
	for ind, videoMeta := range videoBlob.Videos {

		dirName := videoMeta.DirName
		lstDir := strings.Split(dirName, "/")

		if maxDepth > 0 && len(lstDir) > maxDepth {
			dirName = path.Join(lstDir[0:maxDepth]...)
		}

		if nil == videoBlob.DirMap[dirName] {
			videoBlob.DirMap[dirName] = make([]*Video, 0)
		}
		videoBlob.DirMap[dirName] = append(videoBlob.DirMap[dirName], &videoBlob.Videos[ind])
	}

	for finished := false; !finished; {
		finished = true
		for dirName, videos := range videoBlob.DirMap {
			if len(videos) >= videoBlob.ShrinkThreshold {
				continue
			}
			dirNameParent := path.Dir(dirName)
			nameCurrent := path.Base(dirName)
			if strings.Trim(dirName, " ") == "" ||
				dirNameParent == dirName {
				continue;
			}

			if videoBlob.DirMap[dirNameParent] == nil {
				videoBlob.DirMap[dirNameParent] = make([]*Video, 0, videoBlob.ShrinkThreshold)
			}

			for _, pVideo := range videoBlob.DirMap[dirName] {
				pVideo.FileName = path.Join(nameCurrent, pVideo.FileName)
				videoBlob.DirMap[dirNameParent] = append(videoBlob.DirMap[dirNameParent], pVideo)
			}

			delete(videoBlob.DirMap, dirName)
			finished = false
		}
	}
}

func (videoBlob *VideoBlob) scanVideos(baseName string, depth int) error {
	rd, err := ioutil.ReadDir(baseName)
	count := 0
	for _, fi := range rd {
		fileName := fi.Name()

		pathFile := path.Join(baseName, fileName)
		if fi.IsDir() {
			pathDir := pathFile
			videoBlob.scanVideos(pathDir, depth+1)
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
