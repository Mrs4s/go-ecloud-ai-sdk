package ecloud

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/guonaihong/gout"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
)

type (
	FaceDetectResponse struct {
		ImageId   int64               `json:"imageId"`
		Cost      int64               `json:"cost"`
		FaceCount int32               `json:"faceNum"`
		Details   []*FaceDetectDetail `json:"faceDetectDetailList"`
	}

	FaceDetectDetail struct {
		FaceId             string              `json:"faceId"`
		FaceRectangleArea  *RectangleAreaInfo  `json:"faceDectectRectangleArea"`
		LandmarkAreas      []*LandmarkAreaInfo `json:"faceDetectLandmarkAreaList"`
		FaceScore          float64             `json:"faceScore"`
		Roll               float64             `json:"roll"`
		Pitch              float64             `json:"pitch"`
		Raw                float64             `json:"raw"`
		FaceLandMarkNumber int32               `json:"faceLandmarkNumber"`
		// attribute
	}

	RectangleAreaInfo struct {
		UpperLeftX  float64 `json:"upperLeftX"`
		UpperLeftY  float64 `json:"upperLeftY"`
		LowerRightX float64 `json:"lowerRightX"`
		LowerRightY float64 `json:"lowerRightY"`
	}

	LandmarkAreaInfo struct {
		X float64 `json:"pointX"`
		Y float64 `json:"pointY"`
	}
)

func FileFromUrl(url string) io.Reader {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil
	}
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer func() { _ = rsp.Body.Close() }()
	data, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil
	}
	return bytes.NewReader(data)
}

func (as *AiService) FaceDetect(file io.Reader) (*FaceDetectResponse, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "read stream error")
	}
	var rsp string
	if err = as.g.POST("https://smartlib-api-changsha-1.cmecloud.cn:8444/ecloud/ai/v1/face/v1/detect").
		RequestUse(as.auth()).SetJSON(gout.H{"imageFile": base64.StdEncoding.EncodeToString(data)}).
		BindBody(&rsp).Do(); err != nil {
		return nil, errors.Wrap(err, "request api error")
	}
	if state := gjson.Get(rsp, "state").Str; state != "OK" {
		return nil, errors.Errorf("response state error: %v", state)
	}
	var body FaceDetectResponse
	if err = json.Unmarshal([]byte(gjson.Get(rsp, "body").String()), &body); err != nil {
		return nil, errors.Wrap(err, "unmarshal json error")
	}
	return &body, nil
}

func (as *AiService) CreateFaceSet() {

}
