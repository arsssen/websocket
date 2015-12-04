package websocket

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io/ioutil"
	"strings"
)

func (c *Conn) compressMessage(p []byte) (out []byte, err error) {
	compressionAlgo := strings.TrimSuffix(strings.Split(c.compression, " ")[0], ";")
	switch compressionAlgo {
	case "permessage-deflate":
		out, err = compressMessageWithFlate(p)
	default:
		out, err = nil, fmt.Errorf("Compression not supported: '%s'", compressionAlgo)
	}
	return
}

func (c *Conn) decompressMessage(p []byte) (out []byte, err error) {
	compressionAlgo := strings.TrimSuffix(strings.Split(c.compression, " ")[0], ";")
	switch compressionAlgo {
	case "permessage-deflate":
		out, err = decompressFlateMessage(p)
	default:
		out, err = nil, fmt.Errorf("Compression not supported: '%s'", compressionAlgo)
	}
	return
}

func decompressFlateMessage(p []byte) (out []byte, err error) {
	d := append(p, 0x00, 0x00, 0xff, 0xff)
	in := bytes.NewBuffer(d)
	decompressor := flate.NewReader(in)
	defer decompressor.Close()

	out, err = ioutil.ReadAll(decompressor)
	if err != nil {
		if err.Error() != "unexpected EOF" {
			return
		} else {
			err = nil
			//fmt.Println("Fix 'unexpected EOF' on decompression!")
		}
	}
	return
}

func compressMessageWithFlate(data []byte) (out []byte, err error) {
	var (
		buff = new(bytes.Buffer)
		flt  *flate.Writer
	)

	if flt, err = flate.NewWriter(buff, 1); err != nil {
		return
	}
	defer flt.Close()

	if _, err = flt.Write(data); err != nil {
		return
	}

	if err = flt.Flush(); err != nil {
		return
	}

	out = buff.Bytes()
	if out[len(out)-1] != 0x00 {
		out = append(out, 0x00)
	}

	out = out[:len(out)-5]

	return
}
