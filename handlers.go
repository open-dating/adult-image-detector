package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type ImageScoringResult struct {
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
}

func proceedPDF(w http.ResponseWriter, r *http.Request) {
	parsedForm, err := HandleUploadFileForm(r)
	if err != nil {
		HandleError(w, err)
		return
	}

	if parsedForm.FileExt != strings.ToLower("pdf") {
		ProceedImage(w, r)
		return
	}

	images, dir, err := getImagesFromPDF(parsedForm.FilePath)
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

	for _, v := range images {
		var r testResult
		// r.ImageName = v
		if !parsedForm.disableOpenNsfw {
			openNsfwScore, err := GetOpenNsfwScore(v)
			if err != nil {
				continue
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
	RemoveFile(parsedForm.FilePath)
	RemoveFile(dir)

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
func ProceedImage(w http.ResponseWriter, r *http.Request) {
	// save image and get options
	parsedForm, err := HandleUploadFileForm(r)
	if err != nil {
		HandleError(w, err)
		return
	}

	log.Printf("Uploaded file %s, saved as %s", parsedForm.Filename, parsedForm.SaveAsFilename)

	res := ImageScoringResult{
		ImageName:  parsedForm.Filename,
		AppVersion: VERSION,
	}

	if parsedForm.disableOpenNsfw == false {
		// get yahoo open nfsw score
		openNsfwScore, err := GetOpenNsfwScore(parsedForm.FilePath)
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
	RemoveFile(parsedForm.FilePath)

	// serialize answer
	js, err := json.Marshal(res)
	if err != nil {
		HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// error handling
func HandleError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), 500)
}

// show form for upload file
func ShowForm(w http.ResponseWriter) {
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
