package l8f

type VideoHeader struct {
	Meta               map[string]string
	LaptopUniqueFrames [][]int
	LaptopFrames       map[int]int
	MobileUniqueFrames [][]int
	MobileFrames       map[int]int
	AudioSize          int
	LaptopVideoSize    int
	MobileVideoSize    int
}

type UniqueFrameDetails struct {
	Hash             string
	FirstFrameNumber int
	Size             int
}

type MakeVideoLumpTemp struct {
	UniqueFrames                []UniqueFrameDetails
	FramesPointerToUniqueFrames map[int]int
}
