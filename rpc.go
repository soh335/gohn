package gohn

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Executor struct {
	DataDir string
}

type Source struct {
	VideoUrl string
	Offset   string
	Duration string
}

type M4A struct {
	FileName string
}

type Video struct {
	Id string
}

type PlayResponse struct {
	Msg string
}

func (e *Executor) Convert(source *Source, m4a *M4A) error {
	log.Println("call convert", source, m4a)
	filename, err := e.getM4A(source)
	if err != nil {
		return err
	} else {
		m4a.FileName = filename
		return nil
	}
}

func (e *Executor) Play(m4a *M4A, res *PlayResponse) error {
	log.Println("call play", m4a)
	path, err := exec.LookPath("afplay")
	if err != nil {
		return err
	}

	go func() {
		cmd := exec.Command(path, m4a.FileName)
		log.Println(cmd.Args)
		err = cmd.Run()

		if err != nil {
			log.Println(err)
		}
	}()

	res.Msg = "ok"

	return nil
}

func StartRpcServer(host string, port string, dataDir string) {
	executor := new(Executor)
	executor.DataDir = dataDir

	server := rpc.NewServer()
	server.Register(executor)

	addr := net.JoinHostPort(host, port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("start rpc server", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

func (e *Executor) getM4A(source *Source) (string, error) {
	id, downloadUrl, err := e.getIdAndUrlFromVideoUrl(source.VideoUrl)
	if err != nil {
		return "", err
	}

	video := &Video{id}

	if !e.isExistM4A(video) {
		log.Println("does not exist m4a. start downloading")
		file, err := ioutil.TempFile(os.TempDir(), "gohn")
		if err != nil {
			return "", err
		}
		defer func() {
			err := os.Remove(file.Name())
			if err != nil {
				log.Println(err)
			}
		}()

		if err := downloadMP4(downloadUrl, file); err != nil {
			return "", err
		}
		if err := extractM4A(file, source.Offset, source.Duration, e.getM4APath(video)); err != nil {
			return "", err
		}
	}

	return e.getM4APath(video), nil
}

func (e *Executor) isExistM4A(v *Video) bool {
	_, err := os.Stat(e.getM4APath(v))
	if err != nil {
		return false
	}
	return true
}

func (e *Executor) getM4APath(v *Video) string {
	return filepath.Join(e.DataDir, v.Id+".m4a")
}

func (e *Executor) getIdAndUrlFromVideoUrl(videoUrl string) (string, string, error) {
	path, err := exec.LookPath("youtube-dl")
	if err != nil {
		return "", "", err
	}

	cmd := exec.Command(path, "--get-id", "--get-url", videoUrl)
	log.Println(cmd.Args)

	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()

	if err != nil {
		return "", "", err
	}

	ret := strings.Split(out.String(), "\n")
	id := ret[0]
	downloadUrl := ret[1]

	return id, downloadUrl, nil
}

func downloadMP4(url string, file *os.File) error {

	resp, err := http.Get(url)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	io.Copy(file, resp.Body)

	return nil
}

func extractM4A(file *os.File, offset string, duration string, output string) error {
	path, err := exec.LookPath("ffmpeg")
	if err != nil {
		return err
	}

	if offset == "" {
		offset = "0"
	}

	args := []string{"-ss", offset}

	if duration != "" {
		args = append(args, "-t", duration)
	}

	args = append(args, "-i", file.Name(), "-vn", "-y", "-acodec", "copy", output)

	cmd := exec.Command(path, args...)
	log.Println(cmd.Args)

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
