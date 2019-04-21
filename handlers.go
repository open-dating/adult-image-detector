package main

import (
	"encoding/json"
	"net/http"
)

type ImageScoringResult struct {
	AppVersion string `json:"app_version"`
	OpenNsfwScore float32 `json:"open_nsfw_score"`
	ImageName string `json:"image_name"`
}

// save uploaded image and get scoring
func ProceedImage(w http.ResponseWriter, r *http.Request)  {
	// save image
	filePath, imageName, err := SaveUploadFile(r)
	if err != nil {
		HandleError(w, err)
		return
	}

	// get yahoo open nfsw score
	openNsfwScore, err := GetOpenNsfwScore(filePath)
	if err != nil {
		HandleError(w, err)
		return
	}

	// prepare result
	res := ImageScoringResult{
		OpenNsfwScore: openNsfwScore,
		ImageName: imageName,
		AppVersion: VERSION,
	}

	// remove uploaded file
	RemoveFile(filePath)

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
func ShowForm(w http.ResponseWriter)  {
	form := `
<html>
	<body>
		<form action="/api/v1/detect" method="post" enctype="multipart/form-data">
			<label>Select image file</label>
			<input type="file" name="image" required accept="image/*">
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
