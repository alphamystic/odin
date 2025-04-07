package main

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"net"
	"net/netip"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

func SMain() {
	_, _ = os.Stdout.WriteString(string(net.ParseIP("127.0.0.1").To4()))
	_, _ = os.Stdout.Write([]byte{'\n'})
	serve()
}

type writer struct {
	write func(p []byte) (n int, err error)
}

func (w writer) Write(p []byte) (n int, err error) {
	return w.write(p)
}

func fmtHex(b []byte) string {
	s := strings.Builder{}
	s.WriteString("{")
	for i := range b {
		s.WriteString("0x")
		s.WriteString(hex.EncodeToString([]byte{b[i]}))
		if i < len(b)-1 {
			s.WriteString(", ")
		}
	}
	s.WriteString("}")
	return s.String()
}

const (
	red         = "\033[31m"
	yellow      = "\033[33m"
	orange      = "\033[38;5;208m"
	purple      = "\033[35m"
	blue        = "\033[34m"
	green       = "\033[32m"
	gray        = "\033[90m"
	reset       = "\033[0m"
	errPrefix   = red + "(FATAL) "
	warnPrefix  = orange + "(WARN!️️) " + reset
	startPrefix = green + "(START) " + reset
	closePrefix = red + "(CLOSE) " + reset
	infoPrefix  = gray + "(DEBUG) " + reset
	writePrefix = yellow + "(W--->)" + reset
	readPrefix  = purple + "(<---R)" + reset
	finPrefix   = green + "(FINAL) " + reset
)

func handle(c net.Conn) {
	_ = c.SetDeadline(time.Now().Add(time.Duration(5) * time.Second))

	log := func(s string, isErr ...string) {
		dest := os.Stderr
		prefix := "[" + c.RemoteAddr().String() + "]\t"
		if len(isErr) > 0 && len(isErr[0]) > 0 {
			prefix = prefix[:len(prefix)-1]
			prefix += isErr[0]
		}
		_, _ = dest.Write([]byte(prefix))
		_, _ = dest.Write([]byte{byte('\t')})
		_, _ = dest.Write([]byte(s))
		if len(isErr) != 1 || isErr[0] != "" {
			_, _ = dest.Write([]byte{byte('\n')})
		}
		_, _ = dest.Write([]byte(reset))
	}

	var finished = false

	log0 := func(s string) {
		log(s, infoPrefix)
	}

	logRead := func(s string) {
		log(s, readPrefix)
	}

	logWrite := func(s string) {
		log(s, writePrefix)
	}

	log1 := func(s string) {
		log(s, errPrefix)
	}

	log2 := func(s string) {
		log(s, warnPrefix)
	}

	logFin := func(s string) {
		if finished {
			log(s, finPrefix)
		} else {
			log(s, closePrefix)
		}
	}

	logWriter := writer{write: func(p []byte) (n int, err error) {
		log(string(p))
		return len(p), nil
	}}

	log(gray + "-----------------------------------------------")

	log("connection established", startPrefix)
	log("")

	defer logFin("connection closed")

	log(gray + "reading 2 bytes of protocol negotiation data...")

	buf := make([]byte, 2)
	var r int
	var e error
	r, e = c.Read(buf)
	if e != nil {
		log1(e.Error())
		return
	}

	enc := hex.Dumper(logWriter)

	dump := func(buf []byte) {
		log0("dumping buffer:")
		_, _ = enc.Write(buf)
	}

	if r < 2 {
		log1("short read:")
		dump(buf[:r])
		_ = c.Close()
		return
	}

	head := 0

	if buf[head] != 0x05 {
		log1("bad version")
		dump(buf)
		_ = c.Close()
		return
	}

	logRead("\tproto version:\t" +
		blue + "0x" + hex.EncodeToString([]byte{buf[head]}) +
		reset + gray + "\t(int: " + strconv.Itoa(int(buf[head])) + ")",
	)

	head++

	numMethods := int(buf[head])
	if numMethods < 1 {
		log1("no methods")
		dump(buf)
		_ = c.Close()
		return
	}
	if numMethods > 255 {
		log1("too many methods (>255)")
		dump(buf)
		_ = c.Close()
		return
	}

	logRead("\tauth offer ct:\t" +
		blue + "0x" + hex.EncodeToString([]byte{buf[head]}) +
		reset + gray + "\t(int: " + strconv.Itoa(int(buf[head])) + ")",
	)
	log("")

	if buf[head] > 1 {
		log(gray + "expanding buffer for auth methods...")
		buf = append(buf, make([]byte, numMethods)...)
	}

	var authMethods = map[string]bool{
		"anonymous": false,
		"gss-api":   false,
		"user/pass": false,
	}

	updateAuthMethods := func(b byte) {
		logRead("\tauth methods +=\t" + blue + "0x" + hex.EncodeToString([]byte{b}) +
			reset + gray + "\t(int: " + strconv.Itoa(int(b)) + ")" + reset)
		switch {
		case b > 0x02:
			switch {
			case b > 0x02 && b < 0x7f:
				log2("iana assigned auth method used by client: " + strconv.Itoa(int(b)))
			case b > 0x7f && b < 0xfe:
				log2("reserved auth method used by client: " + strconv.Itoa(int(b)))
			case b == 0xff:
				log2("no acceptable auth methods (0xff) sent by client")
			default:
				log2("unknown auth method: " + strconv.Itoa(int(b)))
			}
		default:
			switch b {
			case 0x00:
				authMethods["anonymous"] = true
			case 0x01:
				authMethods["gss-api"] = true
			case 0x02:
				authMethods["user/pass"] = true
			default:
				log2("unknown auth method: " + strconv.Itoa(int(b)))
			}
		}
	}

	printAuthMethods := func() {
		log("")
		hdr := gray + "----- client auth methods -----"
		log(hdr)
		var res = "N/A"
		for k, v := range authMethods {
			res = red + strconv.FormatBool(v) + reset
			if v {
				res = green + strconv.FormatBool(v) + reset
			}
			log0("   " + k + ":\t\t" + res)
		}
		log(gray + strings.Repeat("-", len(hdr)-5))
		log("")
	}

	miniBuf := make([]byte, 1)
	oldHead := head + 1
	for head++; head-oldHead < numMethods; head++ {
		var e error
		if r, e = c.Read(miniBuf); e != nil {
			println(e.Error())
			_ = c.Close()
			return
		}
		if r < 1 {
			println("short read")
			dump(buf)
			_ = c.Close()
			return
		}
		updateAuthMethods(miniBuf[0])
		copy(buf[head:], miniBuf)
		miniBuf = slices.Delete(miniBuf, 0, 0)
	}

	printAuthMethods()

	if !authMethods["anonymous"] {
		log1("does not support anonymous auth")
		// 0xff no acceptable auth methods
		resp := []byte{0x05, 0xff}
		logWrite(red + fmtHex(resp) + gray + " (no acceptable auth methods)")
		_, _ = c.Write(resp)
		_ = c.Close()
		return
	}

	resp := []byte{0x05, 0x00}
	logWrite(green + fmtHex(resp) + gray + "\t(successful auth)")
	log("")

	written, e := c.Write(resp)
	if e != nil {
		log1(e.Error())
		_ = c.Close()
		return
	}

	if written != 2 {
		log1("short write")
		_ = c.Close()
		return
	}

	log(gray + "reading 10 bytes of request data...")

	buf = append(buf, make([]byte, 10)...)

	r, e = c.Read(buf[head:])

	if e != nil {
		log1(e.Error())
		_ = c.Close()
		return
	}

	if r < 10 {
		log1("short read")
		dump(buf[:head+r])
		_ = c.Close()
		return
	}

	if buf[head] != 0x05 {
		log1("bad version")
		dump(buf[head:])
		_ = c.Close()
		return
	}

	logRead("\tproto version:\t" +
		blue + "0x" + hex.EncodeToString([]byte{buf[head]}) +
		reset + gray + "\t(int: " + strconv.Itoa(int(buf[head])) + ")",
	)

	head++

	if buf[head] != 0x01 {
		log1("bad command")
		dump(buf[head:])
		_ = c.Close()
		return
	}

	logRead("\tcommand:\t" +
		blue + "0x" + hex.EncodeToString([]byte{buf[head]}) +
		reset + gray + "\t(int: " + strconv.Itoa(int(buf[head])) + ") (connect)",
	)

	head++

	if buf[head] != 0x00 {
		log1("reserved header not zero")
		dump(buf[head:])
		_ = c.Close()
		return
	}

	logRead(gray + "\treserved:\t" + blue + "0x" + hex.EncodeToString([]byte{buf[head]}))
	log("")

	head++

	if buf[head] != 0x01 {
		log1("bad address type, only ipv4 address supported")
		dump(buf[head:])
		resp := []byte{0x05, 0x08}
		logWrite(red + fmtHex(resp) + reset + gray + " (bad address type)")
		_, _ = c.Write(resp)
		_ = c.Close()
		return
	}

	target := net.IP{
		buf[head+1], buf[head+2], buf[head+3], buf[head+4],
	}

	logRead("\ttarget addr:\t" +
		blue + fmtHex([]byte(target)),
	)
	log("\t\t" + gray + "   \\--> (int: " +
		strconv.Itoa(int(binary.BigEndian.Uint32(target))) +
		") (" + target.String() + ")",
	)

	head += 5

	portSlice := []byte{buf[head], buf[head+1]}
	var port uint16
	for i := range portSlice {
		if portSlice[i] < 0 || portSlice[i] > 255 {
			log1("bad port")
			dump(buf[head:])
			resp := []byte{0x05, 0x01}
			logWrite(red + fmtHex(resp) + reset + gray + " (bad port)")
			_, _ = c.Write(resp)
			_ = c.Close()
			return
		}
		if i == 0 {
			port = uint16(portSlice[i])
		} else {
			port = port<<8 | uint16(portSlice[i])
		}
	}

	logRead("\ttarget port:\t" + blue + fmtHex(portSlice))
	log("\t\t" + gray + "   \\--> (int: " + strconv.Itoa(int(port)) + ")")

	targetStr := target.String() + ":" + strconv.Itoa(int(port))

	ap, err := netip.ParseAddrPort(targetStr)
	if err != nil {
		log1(err.Error())
		dump(buf[head:])
		resp := []byte{0x05, 0x01}
		logWrite(red + fmtHex(resp) + reset + gray + " (general failure)")
		_, _ = c.Write(resp)
		_ = c.Close()
		return
	}
	targetHost := ap.String()

	log("")
	log(gray + "beginning connection to target host...")

	log("\tconnecting to "+targetHost+"... ", "")

	var conn net.Conn
	if conn, e = net.DialTimeout("tcp", targetHost, time.Duration(5)*time.Second); e != nil {
		_, _ = os.Stderr.Write([]byte(red + "failed" + reset + "\n"))
		log("")
		log1(e.Error())
		errResp := []byte{0x05, 0x01}
		logWrite(red + fmtHex(errResp) + reset + gray + " (general failure)")
		_, _ = c.Write(errResp)
		_ = c.Close()
		return
	}

	_, _ = os.Stderr.Write([]byte(green + "success!" + reset + "\n"))
	log("")

	localAddr := c.LocalAddr().(*net.TCPAddr).IP.To4()
	localPortUint16 := uint16(c.LocalAddr().(*net.TCPAddr).Port)
	localPortBytes := []byte{byte(localPortUint16 >> 8), byte(localPortUint16)}

	resp = []byte{0x05, 0x00, 0x00, 0x01}
	logWrite(green + fmtHex(resp) + reset + gray + " (success)")
	written, e = c.Write(resp)
	if e != nil {
		log1(e.Error())
		return
	}
	if written != 4 {
		log1("short write")
		return
	}

	logWrite(blue + fmtHex(localAddr) + reset + gray + " (my addr: " + localAddr.String() + ")")
	written, e = c.Write(localAddr)
	if e != nil {
		log1(e.Error())
		return
	}
	if written != 4 {
		log1("short write")
		dump(localAddr[:written])
		return
	}

	logWrite(blue + fmtHex(localPortBytes) + reset + gray + "\t\t (my port: " + strconv.Itoa(int(localPortUint16)) + ")")
	log("")

	written, e = c.Write(localPortBytes)

	if e != nil {
		log1(e.Error())
		return
	}
	if written != 2 {
		log1("short write")
		dump(localPortBytes[:written])
		return
	}

	log(gray+"beginning proxy i/o goroutines...", "")

	defer func() { _ = conn.Close() }()
	var totalRead, totalWritten int
	totalRead, totalWritten, e = pipe(c, conn)
	switch {
	case errors.Is(e, io.EOF):
		finished = true
		_, _ = os.Stderr.WriteString(green + " EOF " + reset + "\n")
	case e == nil:
		finished = true
		_, _ = os.Stderr.WriteString(green + " FIN " + reset + "\n")
	default:
		_, _ = os.Stderr.WriteString(red + " ERR " + reset + "\n")
		log1(e.Error())
	}
	log0("bytes read: " + strconv.Itoa(totalRead) + "\tbytes written: " + strconv.Itoa(totalWritten))
	log("")
}

func pipe(socksClient net.Conn, target net.Conn) (totalRead int, totalWritten int, err error) {
	defer func() { _ = target.Close() }()
	// caller closes socksClient
	// defer func() { _ = socksClient.Close() }()

	outBuf := make([]byte, 1024)
	inBuf := make([]byte, 1024)
	eChan := make(chan error, 2)

	ctx, cancel := context.WithCancel(context.Background())

	totalRead = 0
	totalWritten = 0

	go func() {
		defer cancel()
		for {
			_ = target.SetDeadline(time.Now().Add(time.Duration(5) * time.Second))
			select {
			case <-ctx.Done():
				return
			default:
			}
			n, e := target.Read(inBuf)
			if e != nil {
				eChan <- e
				return
			}
			if n == 0 {
				eChan <- io.EOF
				return
			}
  CreatedAt string
  UpdatedAt string
}

// 	aptname 	code 	id 	description 	active 	created_at 	updated_at

func CreateAPT(a Apt)error{
			totalRead += n

			_, e = socksClient.Write(inBuf[:n])
			if e != nil {
				eChan <- e
				return
			}
		}
	}()

	go func() {
		defer cancel()
		for {
			_ = socksClient.SetDeadline(time.Now().Add(time.Duration(5) * time.Second))
			select {
			case <-ctx.Done():
				return
			default:
			}
			n, e := socksClient.Read(outBuf)
			if e != nil {
				eChan <- e
				return
			}
			if n == 0 {
				eChan <- io.EOF
				return
			}

			n, e = target.Write(outBuf[:n])
			totalWritten += n

			if e != nil {
				eChan <- e
				return
			}
		}
	}()

	select {
	case e := <-eChan:
		return totalRead, totalWritten, e
	case <-ctx.Done():
		return totalRead, totalWritten, nil
	}
}

func serve() {
	if len(os.Args) < 2 {
		return
	}

	var l net.Listener
	var e error
	go func() {
		l, e = net.Listen("tcp", os.Args[1])
		if e != nil {
			println(e.Error())
			os.Exit(1)
		}
	}()

	time.Sleep(time.Duration(5) * time.Millisecond)
	println("listening on " + os.Args[1])

	for {
		c, e := l.Accept()
		if e != nil {
			println(e.Error())
			continue
		}
		go handle(c)
	}
}
