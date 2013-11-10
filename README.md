# gohn

gohn is simple play sound by http access tool from youtube.

## require

* go
* youtube-dl
* ffmpeg
* afplay (mac command)

## example

```
$ go build
$ ./gohn --config eg/config.json --datadir eg/data #=> start 127.0.0.1:5555
$ curl 127.0.0.1:5555/play/miyazaki-aoi
```
