package main

import (
	"gocv.io/x/gocv"
	"image"
)

func main() {
	webcam, _ := gocv.OpenVideoCapture(0)
	webcam.Set(gocv.VideoCaptureFPS, 30)
	defer webcam.Close()
	window := gocv.NewWindow("Test")
	// kernelTrackbar := window.CreateTrackbar("kernel", 10)
	// dkernelTrackbar := window.CreateTrackbar("dkernel", 100)
	threshTrackbar := window.CreateTrackbar("thresh", 255)
	openingCount := window.CreateTrackbar("openingCount", 10)
	closingCount := window.CreateTrackbar("closingCount", 10)
	defer window.Close()
	output := gocv.NewMat()
	defer output.Close()
	current := gocv.NewMat()
	defer current.Close()
	diff := gocv.NewMat()
	defer diff.Close()

	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
	dkernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(20, 20))

	bg := gocv.NewMat()
	defer bg.Close()
	//capture background
	for i := 0; i < 10; i++ {
		webcam.Read(&bg)
		// gocv.BilateralFilter(output, &output, 9, 100, 100)
		window2 := gocv.NewWindow("bg")
		window2.IMShow(bg)
		window2.WaitKey(1)
	}

	window3 := gocv.NewWindow("current")
	defer window.Close()

	for {
		// kernel = gocv.GetStructuringElement(gocv.MorphCross, image.Pt(kernelTrackbar.GetPos(), kernelTrackbar.GetPos()))
		// dkernel = gocv.GetStructuringElement(gocv.MorphCross, image.Pt(dkernelTrackbar.GetPos(), dkernelTrackbar.GetPos()))
		webcam.Read(&current)
		gocv.AbsDiff(current, bg, &output)
		gocv.CvtColor(output, &output, gocv.ColorBGRToGray)
		gocv.Threshold(output, &output, float32(threshTrackbar.GetPos()), 255, gocv.ThresholdBinary)
		// gocv.Erode(output, &output, kernel)
		// gocv.Dilate(output, &output, kernel)
		for i := 0; i < openingCount.GetPos(); i++ {
			gocv.MorphologyEx(output, &output, gocv.MorphOpen, kernel)
		}
		for i := 0; i < closingCount.GetPos(); i++ {
			gocv.MorphologyEx(output, &output, gocv.MorphClose, dkernel)
		}
		// gocv.BitwiseAnd(current, output, &output)
		// gocv.BitwiseOr(bg, current, &output)
		// gocv.Subtract(current, bg, &output)
		window.IMShow(output)
		window3.IMShow(current)
		window.WaitKey(1)
		window3.WaitKey(1)
	}
}
