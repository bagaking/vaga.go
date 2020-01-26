package main

import "github.com/bagaking/vaga.go/localVideos"

var (
	VBlobStudyArtDrawing = localVideos.VideoBlob{
		Name:        "Art-Drawing",
		Description: "绘画方面的教学视频",
		RootPath:    "/Users/zhouliqihan/geth",
	}
	AllAvailableVideoBlobs = []*localVideos.VideoBlob {
		&VBlobStudyArtDrawing,
	}
)