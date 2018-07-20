package httpClient

import (
	"bufio"
	"io"
	"io/ioutil"
	"net/http"
	"pisc"
	"strings"
)

var ModHTTPRequests = pisc.Module{
	Author:    "Andrew Owen",
	Name:      "HTTPRequests",
	License:   "MIT",
	DocString: "Wrapper around Go's standard lib HTTP client",
	Load:      loadHTTPClient,
}

type httpErr string

func (e httpErr) Error() string {
	return string(e)
}

func doSmallHTTPReq(m *pisc.Machine) error {

	// TODO: Fix bugs
	opts := m.PopValue().(pisc.Dict)
	url := m.PopValue().String()
	verb := m.PopValue().String()

	var body io.Reader = nil
	if str, found := opts["body"]; found {
		body = strings.NewReader(str.String())
	}

	req, err := http.NewRequest(verb, url, body)
	if err != nil {
		return err
	}
	if headers, ok := opts["header"]; ok {
		headersVec, ok := headers.(pisc.Vector)
		if !ok {
			return httpErr("Headers, if present, need to be an array of string pairs")
		}
		for _, val := range headersVec {
			inner := val.(pisc.Vector)
			req.Header.Add(inner[0].String(), inner[1].String())
		}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	reader := resp.Body
	wrapper := pisc.MakeReader(bufio.NewReader(reader))
	wrapper["close"] = pisc.GoFunc(func(m *pisc.Machine) error {
		return reader.Close()
	})

	resp.Cookies()

	result := pisc.Dict{
		"status-code":  pisc.Integer(resp.StatusCode),
		"reply-reader": wrapper,
		// TOOD: Expose headers, cookies, other things
		"content-str": pisc.GoFunc(func(m *pisc.Machine) error {
			body, err := ioutil.ReadAll(reader)
			if err != nil {
				return err
			}
			m.PushValue(pisc.String(string(body)))
			// Close the reader since we have the content loaded down
			return reader.Close()
		}),
	}

	m.PushValue(result)
	return nil
}

func loadHTTPClient(m *pisc.Machine) error {
	m.AddGoWordWithStack(
		"do-http-req",
		"( verb:str url:str options:dict[ ->headers? ->body? ] -- reply:dict[.status-code, .reply-reader, .content-str, ] ) ",
		`Use an HTTP verb, a URL, and a set of options to make an http request, 
		recieving the reply with a status code, reply reader, and a call for the content as a string`,
		doSmallHTTPReq)
	// Todo: look into adding a full suite of HTTP tools
	return nil
}
