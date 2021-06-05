package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gocv.io/x/gocv"
)

type imageScoringResult struct {
	AppVersion                    string  `json:"app_version"`
	OpenNsfwScore                 float32 `json:"open_nsfw_score"`
	AnAlgorithmForNudityDetection bool    `json:"an_algorithm_for_nudity_detection"`
	ImageName                     string  `json:"image_name"`
}

type testResult struct {
	OpenNsfwScore float32 `json:"open_nsfw_score,omitempty"`
	Nudity        bool    `json:"an_algorithm_for_nudity_detection"`
}

type pdfResponse struct {
	AppVersion                  string                `json:"app_version"`
	Result                      map[string]testResult `json:"result"`
	NudityDetectionDisabled     bool                  `json:"nudity_detection_disabled"`
	ImageScoringDisabled        bool                  `json:"image_scoring_disabled"`
	ImageName                   string                `json:"image_name"`
	AlgorithmForNudityDetection bool                  `json:"an_algorithm_for_nudity_detection"`
	OpenNsfwScore               float32               `json:"open_nsfw_score,omitempty"`
}

func proceedPDF(w http.ResponseWriter, r *http.Request) {
	dir, err := ioutil.TempDir(os.TempDir(), "adult-image-detector-*-pdf")
	if err != nil {
		HandleError(w, err)
		return
	}
	parsedForm, err := uploadFileFormHandler(r, dir, "pdf")
	if err != nil {
		HandleError(w, err)
		return
	}

	if parsedForm.FileExt != strings.ToLower("pdf") {
		HandleError(w, fmt.Errorf("bad request. Invalid file type"))
		return
	}

	images, dir, err := getImagesFromPDF(parsedForm.FilePath, dir, parsedForm.password)
	if err != nil {
		HandleError(w, err)
		return
	}

	body := pdfResponse{
		AppVersion:              VERSION,
		ImageScoringDisabled:    parsedForm.disableOpenNsfw,
		NudityDetectionDisabled: parsedForm.disableAnAlgorithm,
		ImageName:               parsedForm.Filename,
	}

	var res = make(map[string]testResult)

	protoPath, _ := filepath.Abs("./models/open_nsfw/nsfw_model/deploy.prototxt")
	modelPath, _ := filepath.Abs("./models/open_nsfw/nsfw_model/resnet_50_1by2_nsfw.caffemodel")

	net := gocv.ReadNetFromCaffe(
		protoPath,
		modelPath,
	)
	if net.Empty() {
		HandleError(w, err)
		return
	}

	defer net.Close()

	for _, v := range images {
		var r testResult
		// r.ImageName = v
		if !parsedForm.disableOpenNsfw {
			openNsfwScore, err := getOpenNsfwScore(v, net)
			if err != nil {
				continue
			}

			if body.OpenNsfwScore < openNsfwScore {
				body.OpenNsfwScore = openNsfwScore
			}

			r.OpenNsfwScore = openNsfwScore
		}

		if !parsedForm.disableAnAlgorithm {
			nudity, err := getAnAlgorithmForNudityDetectionResult(v, parsedForm.debug)
			if err != nil {
				continue
			}

			if nudity {
				body.AlgorithmForNudityDetection = nudity
			}

			r.Nudity = nudity
		}

		res[filepath.Base(v)] = r
	}

	body.Result = res

	// remove uploaded file
	removeFile(parsedForm.FilePath)
	removeFile(dir)

	// serialize answer
	out, err := json.Marshal(body)
	if err != nil {
		HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// save uploaded image and get scoring
func proceedImage(w http.ResponseWriter, r *http.Request) {
	dir, err := ioutil.TempDir(os.TempDir(), "adult-image-detector-*-image")
	if err != nil {
		HandleError(w, err)
		return
	}

	defer removeFile(dir)

	// save image and get options
	parsedForm, err := uploadFileFormHandler(r, dir, "image")
	if err != nil {
		HandleError(w, err)
		return
	}

	log.Printf("Uploaded file %s, saved as %s", parsedForm.Filename, parsedForm.SaveAsFilename)

	res := imageScoringResult{
		ImageName:  parsedForm.Filename,
		AppVersion: VERSION,
	}

	protoPath, _ := filepath.Abs("./models/open_nsfw/nsfw_model/deploy.prototxt")
	modelPath, _ := filepath.Abs("./models/open_nsfw/nsfw_model/resnet_50_1by2_nsfw.caffemodel")

	net := gocv.ReadNetFromCaffe(
		protoPath,
		modelPath,
	)
	if net.Empty() {
		HandleError(w, err)
		return
	}

	defer net.Close()

	if parsedForm.disableOpenNsfw == false {
		// get yahoo open nfsw score
		openNsfwScore, err := getOpenNsfwScore(parsedForm.FilePath, net)
		if err != nil {
			HandleError(w, err)
			return
		}
		res.OpenNsfwScore = openNsfwScore

		log.Printf("For file %s, openNsfwScore=%f", parsedForm.SaveAsFilename, openNsfwScore)
	}

	if parsedForm.disableAnAlgorithm == false {
		// get An Algorithm for Nudity Detection
		anAlgorithmForNudityDetection, err := getAnAlgorithmForNudityDetectionResult(parsedForm.FilePath, parsedForm.debug)
		if err != nil {
			HandleError(w, err)
			return
		}
		res.AnAlgorithmForNudityDetection = anAlgorithmForNudityDetection

		log.Printf("For file %s, anAlgorithmForNudityDetection=%t", parsedForm.SaveAsFilename, anAlgorithmForNudityDetection)
	}

	// remove uploaded file
	removeFile(parsedForm.FilePath)

	// serialize answer
	out, err := json.Marshal(res)
	if err != nil {
		HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// HandleError handles http error.
func HandleError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), 500)
}

// ShowImageForm to upload jpg image.
func ShowImageForm(w http.ResponseWriter) {
	form := `
<html>
	<body>
		<form action="/api/v1/detect" method="post" enctype="multipart/form-data">
			<label>Select image file</label>
			<input type="file" name="image" required accept="image/*"><br/>
			
			<label>
				<input type="checkbox" value="true" name="disableOpenNsfw"> disable open nsfw
			</label><br/>
			<label>
				<input type="checkbox" value="true" name="disableAnAlgorithm"> disable an algorithm
			</label><br/>
			
			<button type="submit">Calc nude scores</button>
		</form>
		<pre>
curl -i -X POST -F "image=@Daddy_Lets_Me_Ride_His_Cock_preview_720p.mp4.jpg" http://localhost:9191/api/v1/detect
		</pre>
	</body>
</html>
`
	w.Write([]byte(form))
}

// ShowPDFForm to upload pdf file.
func ShowPDFForm(w http.ResponseWriter) {
	form := `
<html>
	<body>
		<form action="/api/v1/detect_pdf" method="post" enctype="multipart/form-data">
			<label>Select image file</label>
			<input type="file" name="pdf" required accept="application/pdf"><br/>
			
			<label>
				<input type="checkbox" value="true" name="disableOpenNsfw"> disable open nsfw
			</label><br/>
			<label>
				<input type="checkbox" value="true" name="disableAnAlgorithm"> disable an algorithm
			</label><br/>
			
			<button type="submit">Calc nude scores</button>
		</form>
		<pre>
curl -i -X POST -F "image=@Daddy_Lets_Me_Ride_His_Cock_preview_720p.mp4.pdf" http://localhost:9191/api/v1/detect_pdf
		</pre>
	</body>
</html>
`
	w.Write([]byte(form))
}
