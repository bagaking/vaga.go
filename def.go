package main

import "bagaking.com/vaga.go/localVideos"

var (
	VBlobStudyArtDrawing = localVideos.VideoBlob{
		Name:        "Art-Drawing",
		Description: "绘画方面的教学视频",
		RootPath:    "K:\\美术-绘画\\",
	}
	AllAvailableVideoBlobs = []*localVideos.VideoBlob {
		&VBlobStudyArtDrawing,
	}
)