package repeater

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
)

const (
	maxHostNameLen          = 250
	rfbPortOffset           = 5900
	rfbProtoVerMsgLen       = 12
	repeaterVerMsg          = "RFB 000.000\n"
	repeaterVerKeepAliveMsg = "REP 000.000\n"
)

var reRepeaterID = regexp.MustCompile(`^ID:(\d+)`)

func sendRepeaterVersion(w io.Writer) (err error) {

	_, err = io.WriteString(w, repeaterVerMsg)

	return
}

func fetchHostInfo(r io.Reader) (token string, err error) {

	buf := bytes.NewBuffer(nil)

	if _, err = io.CopyN(buf, r, maxHostNameLen); err != nil {
		return
	}

	if m := reRepeaterID.FindSubmatch(buf.Bytes()); len(m) == 2 {
		if idx := bytes.IndexByte(m[1], 0); idx < 0 {
			token = string(m[1])
		} else {
			token = string(m[1][:idx])
		}
		return
	}

	err = fmt.Errorf("Invalid HostName")

	return
}
