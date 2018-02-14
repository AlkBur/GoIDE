package util

import (
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"github.com/AlkBur/GoIDE/log"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

const (
	vary            = "Vary"
	acceptEncoding  = "Accept-Encoding"
	contentEncoding = "Content-Encoding"
	contentType     = "Content-Type"
	contentLength   = "Content-Length"
	typeGzip        = "gzip"
)

var (
	compressionPool = &sync.Pool{New: func() interface{} { return gzip.NewWriter(nil) }}
)

type H map[string]interface{}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

// Result.
type Result struct {
	Result bool        `json:"result"` // return result
	Code   string      `json:"code"`   // return code
	Msg    string      `json:"msg"`    // message
	Data   interface{} `json:"data"`   // data object
}

func getGzipWriter(w http.ResponseWriter) (gz *gzip.Writer) {
	r := compressionPool.Get()
	if r != nil {
		gz = r.(*gzip.Writer)
		gz.Reset(w)
	} else {
		gz = gzip.NewWriter(w)
	}
	return
}

func putGzipWriter(gw *gzip.Writer) {
	compressionPool.Put(gw)
}

func GzipHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !shouldCompress(r) {
			f(w, r)
			return
		}
		w.Header().Set(contentEncoding, typeGzip)
		w.Header().Set(vary, acceptEncoding)

		gz := getGzipWriter(w)
		defer func() {
			w.Header().Set(contentLength, "0")
			gz.Close()
			putGzipWriter(gz)
		}()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}

		f(gzr, r)
	}
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	if "" == w.Header().Get(contentType) {
		// If no content type, apply sniffing algorithm to un-gzipped body.
		w.Header().Set(contentType, http.DetectContentType(b))
	}

	return w.Writer.Write(b)
}

func shouldCompress(req *http.Request) bool {
	if !strings.Contains(req.Header.Get(acceptEncoding), "gzip") {
		return false
	}
	extension := filepath.Ext(req.URL.Path)
	if len(extension) < 4 { // fast path
		return true
	}

	switch extension {
	case ".png", ".gif", ".jpeg", ".jpg":
		return false
	default:
		return true
	}
}

func LoadTemplate(name ...string) *template.Template {
	return template.Must(template.New("").Delims("[[", "]]").ParseFiles(name...))
}

// MarshalXML allows type H to be used with xml.Marshal.
func (h H) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: "",
		Local: "map",
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for key, value := range h {
		elem := xml.StartElement{
			Name: xml.Name{Space: "", Local: key},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

func NewResult() *Result {
	return &Result{
		Result: true,
		Code:   "0",
		Msg:    "",
		Data:   nil,
	}
}

func (res *Result) Send(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(res)
	if err != nil {
		log.Error(err)
		return
	}
	w.Write(data)
}
