# gohn

gohn is simple play youtube sound by http access tool. it is easy to hit a bell, "gohn".

## require

* go
* youtube-dl
* ffmpeg
* play sound command ( default is afplay )

## example

```
$ go build
$ ./gohn --config eg/config.json --datadir eg/data #=> start 127.0.0.1:5555
$ curl 127.0.0.1:5555/play/miyazaki-aoi
```
## git commit hook

```
$ echo "#/bin/sh\ncurl 127.0.0.1:5555/play/miyazaki-aoi" > .git/hooks/post-commit
$ chmod +x .git/hooks/post-commit
$ touch test
$ git add test
$ git ci -m 'add test'
```

* http://hisaichi5518.hatenablog.jp/entry/2013/11/03/193719

## on network change

```
$ perl -MCocoa::EventLoop -MCocoa::NetworkChange -e 'on_network_change(sub{ system(qw|curl 127.0.0.1:5555/play/miyazaki-aoi-return-home|); }, sub{}); Cocoa::EventLoop->run;'
```
