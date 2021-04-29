package main

import (
	"gocv.io/x/gocv"
	"image"
	"image/color"
)

const MinimumArea = 3000

func main() {
	webcam, _ := gocv.OpenVideoCapture(0)
	webcam.Set(gocv.VideoCaptureFPS, 30)
	defer webcam.Close()
	window := gocv.NewWindow("Test")
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

	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(5, 5))
	dkernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(20, 20))

	bg := gocv.NewMat()
	defer bg.Close()
	//capture background
	for i := 0; i < 10; i++ {
		webcam.Read(&bg)
		gocv.MedianBlur(bg, &bg, 3)
		window2 := gocv.NewWindow("bg")
		window2.IMShow(bg)
		window2.WaitKey(1)
	}
	exposure := webcam.Get(gocv.VideoCaptureExposure)
	webcam.Set(gocv.VideoCaptureAutoExposure, 0)
	webcam.Set(gocv.VideoCaptureExposure, exposure)

	window3 := gocv.NewWindow("current")
	defer window.Close()

	for {
		webcam.Read(&current)
		gocv.MedianBlur(current, &current, 3)
		gocv.AbsDiff(current, bg, &output)
		gocv.CvtColor(output, &output, gocv.ColorBGRToGray)
		gocv.Threshold(output, &output, float32(threshTrackbar.GetPos()), 255, gocv.ThresholdBinary)
		for i := 0; i < openingCount.GetPos(); i++ {
			gocv.MorphologyEx(output, &output, gocv.MorphOpen, kernel)
		}
		for i := 0; i < closingCount.GetPos(); i++ {
			gocv.MorphologyEx(output, &output, gocv.MorphClose, dkernel)
		}

		contours := gocv.FindContours(output, gocv.RetrievalExternal, gocv.ChainApproxNone)

		for i := 0; i < contours.Size(); i++ {
			area := gocv.ContourArea(contours.At(i))
			if area < MinimumArea {
				continue
			}

			statusColor := color.RGBA{255, 0, 0, 0}
			gocv.DrawContours(&current, contours, i, statusColor, 2)
		}
		window.IMShow(output)
		window3.IMShow(current)
		window.WaitKey(1)
		window3.WaitKey(1)
	}
}
