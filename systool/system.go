package systool

import (
	"github.com/ulricqin/goutils/filetool"
	"github.com/ulricqin/goutils/convertor"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"net"
	"strings"
	"math/rand"
	"time"
	"bytes"
)

func WritePidFile(pidFilePath string) error {
	if pidFilePath == "" {
		panic("parameter pidFilePath is blank")
	}

	_, err := filetool.WriteStringToFile(pidFilePath, fmt.Sprintf("%d", os.Getpid()))
	if err != nil {
		return err
	}

	return nil
}


func LocalIP() (string, error) {
	addr, err := net.ResolveUDPAddr("udp", "1.2.3.4:1")
	if err != nil {
		return "", err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return "", err
	}

	defer conn.Close()

	host, _, err := net.SplitHostPort(conn.LocalAddr().String())
	if err != nil {
		return "", err
	}

	// host = "10.180.2.66"
	return host, nil
}

func LocalDnsName() (hostname string, err error) {
	var ip string
	ip, err = LocalIP()
	if err != nil {
		return
	}

	cmd := exec.Command("host", ip)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return
	}

	tmp := out.String()
	arr := strings.Split(tmp, ".\n")

	if len(arr) > 1 {
		content := arr[0]
		arr = strings.Split(content, " ")
		return arr[len(arr)-1], nil
	}

	err = fmt.Errorf("parse host %s fail", ip)
	return
}

func GrabEphemeralPort() (port uint16, err error) {
	var listener net.Listener
	var portStr string
	var p int

	listener, err = net.Listen("tcp", ":0")
	if err != nil {
		return
	}
	defer listener.Close()

	_, portStr, err = net.SplitHostPort(listener.Addr().String())
	if err != nil {
		return
	}

	p, err = strconv.Atoi(portStr)
	port = uint16(p)

	return
}

func URandom() string {
	f, _ := os.Open("/dev/urandom")
	b := make([]byte, 16)
	f.Read(b)
	f.Close()

	return fmt.Sprintf("%x", b)
}

func GenerateRandomSeed() int64 {
	return convertor.BytesToInt64([]byte(URandom()))
}

func SleepRandomDuration(t int) {
	r := rand.New(rand.NewSource(GenerateRandomSeed()))
	d := time.Duration(r.Intn(t)) * time.Millisecond
	time.Sleep(d)
}


