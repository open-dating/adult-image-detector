### adult-image-detector
[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/open-dating/adult-image-detector) 
[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/opendating/adult-image-detector)](https://hub.docker.com/repository/docker/opendating/adult-image-detector)


Use deep neural networks and other algos for detect nude images. 

[Try detection](https://adult-image-detector.herokuapp.com/)

### Usage
Exec:
```
curl -i -X POST -F "image=@Daddy_Lets_Me_Ride_His_Cock_preview_720p.mp4.jpg" http://localhost:9191/api/v1/detect
```
Result:
```json
{
  "app_version":"0.2.0",
  "open_nsfw_score":0.81577397,
  "an_algorithm_for_nudity_detection": true,
  "image_name":"Daddy_Lets_Me_Ride_His_Cock_preview_720p.mp4.jpg"
}
```

### Docker
#### Run
```
docker run -p 9191:9191 opendating/adult-image-detector
```

#### Build
```
git clone https://github.com/open-dating/adult-image-detector --recursive
docker build -t adult-image-detector .
```

#### Development
```
git clone https://github.com/open-dating/adult-image-detector --recursive
cd docker/dev
docker-compose up
```

#### Test
```
cd docker/test
docker-compose up
```

### Install to heroku
Use deploy button or:

fork, create app and change stack to container
```
heroku stack:set container
```

### Requirements
Go 1.11

opencv 3.4.1

### Development without docker
Recursive clone that repo:
```
git clone https://github.com/open-dating/adult-image-detector --recursive
```
or manually install submodules:
```
git submodule init
git submodule update
```

Install dependencies:
```
# get package manager
go get -u github.com/kardianos/govendor

# install dependencies
govendor sync
```

Install opencv 3.4.1

Run hot reload with fresh:
```
go get github.com/pilu/fresh
fresh
```
