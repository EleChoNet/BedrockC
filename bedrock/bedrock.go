package bedrock

import (
	"bufio"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/exec"
	"time"
)

type BedrockS struct {
	path       string
	executable string
	stdOut     io.Reader
	stdIn      io.Writer
	cmd        *exec.Cmd
}

func (b *BedrockS) Run() error {
	cmd := exec.Command(b.executable)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "无法获取stdout")
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return errors.Wrap(err, "无法获取stdin")
	}
	b.stdOut = stdout
	b.stdIn = stdin
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return errors.Wrap(err, "无法启动bedrock")
	}
	b.cmd = cmd
	return nil
}
func (b *BedrockS) Stop() error {
	err := b.cmd.Process.Kill()
	if err != nil {
		return errors.Wrap(err, "无法关闭bedrock")
	}
	return nil
}

//读取新行，如果没有新行，则等待
func (b *BedrockS) ReadLine(timeout time.Duration) (string, error) {
	reader := bufio.NewReader(b.stdOut)
	Now := time.Now()
waitloop:
	line, err := reader.ReadString('\n')
	if err != nil || err == io.EOF {
		if time.Since(Now) > timeout {
			return "", err
		} else {
			goto waitloop
		}
	} else if err != nil {
		return "", errors.Wrap(err, "无法读取MC控制台输出")
	}
	return line, nil

}
func (b *BedrockS) ReadLinesByTime(timeout time.Duration) (string, error) {
	Now := time.Now()
	var result string
	for {
		if time.Since(Now) > timeout {
			break
		}
		line, err := b.ReadLine(time.Millisecond * 100)
		if err != nil {
			if err == io.EOF {
				continue
			} else {
				return "", err
			}
		}
		result += line
	}
	return result, nil
}

func (b *BedrockS) Command(command string) (string, error) {
	_, err := b.stdIn.Write([]byte(command))

	if err != nil {
		return "", errors.Wrap(err, "无法发送命令")
	}
	result, err := b.ReadLinesByTime(time.Millisecond * 200)
	if err != nil {
		return "", errors.Wrap(err, "无法读取命令结果")
	}
	return result, nil
}
