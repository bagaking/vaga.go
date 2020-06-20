package localVideos

import (
	"fmt"
	. "github.com/c2h5oh/datasize"
	"io/ioutil"
	"path"
	"strings"
)

type Video struct {
	NameHash string `json:"hash"`
	Location string `json:"location"`

	DirName   string `json:"dir"`
	FileName  string `json:"name"`
	FileType  string `json:"typ"`
	FileSize  int64  `json:"size"`
	FileSizeH string `json:"size_h"`
}

type VideoBlobConf struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
	RootPath    string `json:"root"`

	MaxDepth        int `json:"max_depth"`
	ShrinkThreshold int `json:"shrink_threshold"`
	ShrinkMaxCount  int `json:"shrink_max_count"`
}

type VideoBlob struct {
	VideoBlobConf

	Videos []*Video            `json:"videos"`
	DirMap map[string][]*Video `json:"dir_map"`
}

//var Videos = make([]Video, 0, 10000)
var VideosMap = map[string]*Video{}

func (videoBlob *VideoBlob) Initial() error {
	if err := videoBlob.scanVideos(videoBlob.RootPath, 0); err != nil {
		return err
	}
	videoBlob.initialDirMap()
	fmt.Print("Videos initialed")
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
		videoBlob.DirMap[dirName] = append(videoBlob.DirMap[dirName], videoBlob.Videos[ind])
	}

	shrinkLeftCount := - videoBlob.ShrinkMaxCount
	for finished := false; !finished && shrinkLeftCount != -1; shrinkLeftCount++ {
		finished = true

		shrinkDirNames := map[string]bool{}
		for dirName, videos := range videoBlob.DirMap {
			if len(videos) >= videoBlob.ShrinkThreshold {
				continue
			}
			dirNameParent := path.Dir(dirName)
			if strings.Trim(dirName, " ") == "" ||
				dirNameParent == dirName {
				continue
			}

			shrinkDirNames[dirNameParent] = true
		}

		for dirName := range videoBlob.DirMap {
			dirNameParent := ""
			dirNameAncient := path.Dir(dirName)
			for ; dirNameAncient != "" && dirNameAncient != path.Dir(dirNameAncient); dirNameAncient = path.Dir(dirNameAncient) {
				if shrinkDirNames[dirNameAncient] {
					dirNameParent = dirNameAncient
				}
			}

			if !shrinkDirNames[dirNameParent] {
				continue
			}

			//fmt.Println("dirNameParent", dirNameParent)

			if videoBlob.DirMap[dirNameParent] == nil {
				videoBlob.DirMap[dirNameParent] = []*Video{}
			}

			dir := videoBlob.DirMap[dirName]
			for i, pVideo := range dir {
				dir[i].FileName = path.Join(strings.Replace(pVideo.DirName, dirNameParent, "", 1), pVideo.FileName)
				dir[i].DirName = dirNameParent
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
		ext := strings.ToLower(path.Ext(fileName))
		if fi.IsDir() {
			pathDir := pathFile
			videoBlob.scanVideos(pathDir, depth+1)
		} else {
			if ext == ".mp4" || ext == ".ogg" || ext == ".webm" {
				//fmt.Println("==", len(videoBlob.Videos), pathFile)
				pVideo := &Video{
					NameHash: Sha1([]byte(pathFile)),
					DirName:  baseName,
					Location: pathFile,

					FileName: fileName,
					FileType: strings.ToLower(path.Ext(fileName)[1:]),
					FileSize: fi.Size(),
				}
				pVideo.FileSizeH = ByteSize(pVideo.FileSize).HR()
				videoBlob.Videos = append(videoBlob.Videos, pVideo)
				VideosMap[pVideo.NameHash] = pVideo
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
