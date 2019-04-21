### adult-image-detector
Use deep neural networks and other algos for detect nude images

### Usage
Exec:
```
curl -i -X POST -F "image=@Daddy_Lets_Me_Ride_His_Cock_preview_720p.mp4.jpg" http://localhost:9191/api/v1/detect
```
Result:
```json
{
  "open_nsfw_score":0.11577397,
  "image_name":"Daddy_Lets_Me_Ride_His_Cock_preview_720p.mp4.jpg"
}
```

### Docker
#### Run
```
docker run -p 9191:9191 grinat0/adult-image-detector
```

#### Build
```
docker build -t adult-image-detector .
```

#### Development
```
cd docker/dev
docker-compose up
```

#### Test
```
cd docker/test
docker-compose up
```

### Install to heroku
Fork, create app and change stack to container
```
heroku stack:set container
```

### Requirements
Go 1.11

opencv 3.4.1

### Development without docker
Install dependencies:
```
# get package manager
go get -u github.com/kardianos/govendor

# install dependencies
govendor sync
```

Clone submodules:
```
git submodule init
git submodule update
```

Install opencv 3.4.1

Run hot reload with fresh:
```
go get github.com/pilu/fresh
fresh
```
