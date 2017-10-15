package baloo

/// Come from https://github.com/emicklei/forest/
// LICENSE MIT https://github.com/emicklei/forest/blob/master/LICENSE.txt

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"testing"
)

var scanStackForFile = true
var logf_func = logf

const noStackOffset = 0

// Logf adds the actual file:line information to the log message
func Logf(t *testing.T, format string, args ...interface{}) {
	logf_func(t, noStackOffset, "\n"+format, args...)
}

func logfatal(t *testing.T, format string, args ...interface{}) {
	logf_func(t, noStackOffset, format, args...)
	t.FailNow()
}

// Error is equivalent to Log followed by Fail.
func logerror(t *testing.T, args ...interface{}) {
	logerrorf(t, "\terror: "+tabify("%s")+"\n", args)
}

func logerrorf(t *testing.T, format string, args ...interface{}) {
	logf_func(t, noStackOffset, format, args...)
	t.Fail()
}

func logf(t *testing.T, stackOffset int, format string, args ...interface{}) {
	var file string
	var line int
	var ok bool
	if scanStackForFile {
		offset := 0
		outside := false
		for !outside {
			_, file, line, ok = runtime.Caller(2 + offset)
			outside = !strings.Contains(file, "/baloo/")
			offset++
		}
	} else {
		_, file, line, ok = runtime.Caller(2)
	}
	if ok {
		// Truncate file name at last file name separator.
		if index := strings.LastIndex(file, "/"); index >= 0 {
			file = file[index+1:]
		} else if index = strings.LastIndex(file, "\\"); index >= 0 {
			file = file[index+1:]
		}
	} else {
		file = "???"
		line = 1
	}
	t.Logf("<-- %s:%d "+format, append([]interface{}{file, line}, args...)...)
}

// Dump is a convenient method to log the full contents of a request and its response.
func Dump(t *testing.T, resp *http.Response) {
	// dump request
	var buffer bytes.Buffer
	buffer.WriteString("\n")
	buffer.WriteString(fmt.Sprintf("%v %v\n", resp.Request.Method, resp.Request.URL))
	for k, v := range resp.Request.Header {
		if len(k) > 0 {
			buffer.WriteString(fmt.Sprintf("%s : %v\n", k, strings.Join(v, ",")))
		}
	}
	if resp == nil {
		buffer.WriteString("-- no response --")
		Logf(t, buffer.String())
		return
	}
	// dump response
	buffer.WriteString(fmt.Sprintf("\n%s\n", resp.Status))
	for k, v := range resp.Header {
		if len(k) > 0 {
			buffer.WriteString(fmt.Sprintf("%s : %v\n", k, strings.Join(v, ",")))
		}
	}
	if resp.Body != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			if resp.StatusCode/100 == 3 {
				// redirect closes body ; nothing to read
				buffer.WriteString("\n")
			} else {
				buffer.WriteString(fmt.Sprintf("unable to read body:%v", err))
			}
		} else {
			if len(body) > 0 {
				buffer.WriteString("\n")
			}
			buffer.WriteString(string(body))
		}
		resp.Body.Close()
		// put the body back for re-reads
		resp.Body = ioutil.NopCloser(bytes.NewReader(body))
	}
	buffer.WriteString("\n")
	Logf(t, buffer.String())
}

func tabify(format string) string {
	if strings.HasPrefix(format, "\n") {
		return strings.Replace(format, "\n", "\n\t\t", 1)
	}
	return format
}
