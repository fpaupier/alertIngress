package main

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"image/jpeg"
	"os"
	"testing"
	"time"
)

const (
	fPath  = "testData/bike.jpeg"
	layout = "2006-01-02 15:04:05.0700"
)

func TestCreateAlert(t *testing.T) {
	f, err := os.Open(fPath)
	defer f.Close()
	if err != nil {
		t.Fatalf("failed to load test image: %v\n", err)
	}
	img, err := jpeg.Decode(f)
	if err != nil {
		t.Fatalf("failed to decode image: %v\n", err)
	}
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, nil)
	imgData := buf.Bytes()

	eventTime := time.Now()
	enrolledOn := time.Now().AddDate(0, 0, -10)
	alert := &Alert{
		EventTime: eventTime.Format(layout),
		CreatedBy: &Alert_Device{
			Type:       "Testing device",
			Guid:       "xyz",
			EnrolledOn: enrolledOn.Format(layout),
		},
		Location: &Alert_Location{
			Longitude: 42.354,
			Latitude:  -72.546,
		},
		FaceDetectionModel: &Alert_Model{
			Name:      "test face detection model",
			Guid:      "abc",
			Threshold: 0.32,
		},
		MaskClassifierModel: &Alert_Model{
			Name:      "test mask classifier model",
			Guid:      "def",
			Threshold: 0.47,
		},
		Probability: 0.975,
		Image: &Alert_Image{
			Format: ".jpeg",
			Size: &Alert_Image_Size{
				Width:  224,
				Height: 224,
			},
			Data: imgData,
		},
	}

	marshaled, err := proto.Marshal(alert)
	if err != nil {
		t.Error("marshaling failed")
	}
	newAlert := createAlert(marshaled)
	if !equal(newAlert, alert) {
		t.Errorf("newAlert different from initial alert")
	}
}

// equal return true if two alerts are equal, false otherwise.
func equal(a1, a2 *Alert) bool {
	return a1.CreatedBy.Guid == a2.CreatedBy.Guid &&
		a1.CreatedBy.EnrolledOn == a2.CreatedBy.EnrolledOn &&
		a1.CreatedBy.Type == a2.CreatedBy.Type &&
		a1.Image.Format == a2.Image.Format &&
		a1.Image.Size.Width == a2.Image.Size.Width &&
		a1.Image.Size.Height == a2.Image.Size.Height &&
		bytes.Compare(a1.Image.Data, a2.Image.Data) == 0 &&
		a1.MaskClassifierModel.Guid == a2.MaskClassifierModel.Guid &&
		a1.MaskClassifierModel.Name == a2.MaskClassifierModel.Name &&
		a1.MaskClassifierModel.Threshold == a2.MaskClassifierModel.Threshold &&
		a1.FaceDetectionModel.Threshold == a2.FaceDetectionModel.Threshold &&
		a1.FaceDetectionModel.Guid == a2.FaceDetectionModel.Guid &&
		a1.FaceDetectionModel.Name == a2.FaceDetectionModel.Name &&
		a1.Location.Latitude == a2.Location.Latitude &&
		a1.Location.Longitude == a2.Location.Longitude &&
		a1.EventTime == a2.EventTime

}
