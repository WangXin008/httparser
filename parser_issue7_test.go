package httparser

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

/*
func TestServerParserRequest_Trailer(t *testing.T) {

	for _, testData := range [][]byte{
		[]byte("POST / HTTP/1.1\r\nHost: localhost:1235\r\nUser-Agent: Go-http-client/1.1\r\nTransfer-Encoding: chunked\r\nTrailer: Md5,Size\r\nAccept-Encoding: gzip\r\n\r\n4\r\nbody\r\n0\r\nMd5: 841a2d689ad86bd1611447453c22c6fc\r\nSize: 4\r\n\r\n"),
	} {
		testParser(t, testData)
	}
}
*/

func TestServerParserRequest(t *testing.T) {

	for _, testData := range [][]byte{
		[]byte("POST /echo HTTP/1.1\r\n\r\n"),
		[]byte("POST /echo HTTP/1.1\r\nHost: localhost:8080\r\nConnection: close \r\nAccept-Encoding : gzip \r\n\r\n"),
		[]byte("POST /echo HTTP/1.1\r\nHost: localhost:8080\r\nConnection: close \r\nContent-Length :  5\r\nAccept-Encoding : gzip \r\n\r\nhello"),
		[]byte("POST / HTTP/1.1\r\nHost: localhost:1235\r\nUser-Agent: Go-http-client/1.1\r\nTransfer-Encoding: chunked\r\nAccept-Encoding: gzip\r\n\r\n4\r\nbody\r\n0\r\n\r\n"),
		[]byte("GET /test HTTP/1.1\r\n" +
			"User-Agent: curl/7.18.0 (i486-pc-linux-gnu) libcurl/7.18.0 OpenSSL/0.9.8g zlib/1.2.3.3 libidn/1.1\r\n" +
			"Host: 0.0.0.0=5000\r\n" +
			"Accept: */*\r\n" +
			"\r\n"),
		[]byte("GET /favicon.ico HTTP/1.1\r\n" +
			"Host: 0.0.0.0=5000\r\n" +
			"User-Agent: Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9) Gecko/2008061015 Firefox/3.0\r\n" +
			"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\n" +
			"Accept-Language: en-us,en;q=0.5\r\n" +
			"Accept-Encoding: gzip,deflate\r\n" +
			"Accept-Charset: ISO-8859-1,utf-8;q=0.7,*;q=0.7\r\n" +
			"Keep-Alive: 300\r\n" +
			"Connection: keep-alive\r\n" +
			"\r\n"),
		[]byte("GET /dumbluck HTTP/1.1\r\n" +
			"aaaaaaaaaaaaa:++++++++++\r\n" +
			"\r\n"),
		[]byte("GET /forums/1/topics/2375?page=1#posts-17408 HTTP/1.1\r\n" +
			"\r\n"),
		[]byte("GET /get_no_headers_no_body/world HTTP/1.1\r\n" +
			"\r\n"),
	} {
		testParser(t, testData)
	}
}

func testParser(t *testing.T, data []byte) error {
	setting := Setting{
		MessageBegin:    func(*Parser) {},
		URL:             func(*Parser, []byte) {},
		Status:          func(*Parser, []byte) {},
		HeaderField:     func(*Parser, []byte) {},
		HeaderValue:     func(*Parser, []byte) {},
		HeadersComplete: func(*Parser) {},
		Body:            func(*Parser, []byte) {},
		MessageComplete: func(*Parser) {},
	}
	p := New(REQUEST)

	var remain []byte
	for i := 0; i < len(data); i++ {
		b := append(remain, data[i:i+1]...)
		n, err := p.Execute(&setting, b)
		if err != nil {
			t.Fatal(fmt.Errorf("%v success, %v", n, err))
		}
		if n < len(b) {
			remain = append([]byte{}, b[n:]...)
		}
	}

	nRequest := 0
	data = append(data, data...)
	setting = Setting{
		MessageBegin:    func(*Parser) {},
		URL:             func(*Parser, []byte) {},
		Status:          func(*Parser, []byte) {},
		HeaderField:     func(*Parser, []byte) {},
		HeaderValue:     func(*Parser, []byte) {},
		HeadersComplete: func(*Parser) {},
		Body:            func(*Parser, []byte) {},
		MessageComplete: func(*Parser) {
			nRequest++
		},
	}

	tBegin := time.Now()
	loop := 10
	for i := 0; i < loop; i++ {
		tmp := data
		var remain []byte

		for len(tmp) > 0 {

			nRead := int(rand.Intn(len(tmp)) + 1)
			readBuf := append(remain, tmp[:nRead]...)

			n, err := p.Execute(&setting, readBuf)
			if err != nil {
				//fmt.Printf("remain:(%s):readBuf(%s) err:%v\n", remain, readBuf, err)
				t.Fatal(fmt.Errorf("%v success, %v", n, err))
				//return nil
			}

			//fmt.Printf("---> n = %d, %d\n", n, len(readBuf))

			if n < len(readBuf) {
				remain = append([]byte{}, readBuf[n:]...)
			}

			tmp = tmp[nRead:]
		}

		fmt.Printf("remain:(%s)\n", remain)

		if nRequest != (i+1)*2 {
			return fmt.Errorf("nRequest: %v, %v", i, nRequest)
		}
	}

	tUsed := time.Since(tBegin)
	fmt.Printf("%v loops, %v s used, %v ns/op, %v req/s\n", loop, tUsed.Seconds(), tUsed.Nanoseconds()/int64(loop), float64(loop)/tUsed.Seconds())

	return nil
}
