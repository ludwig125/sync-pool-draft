package main

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/google/go-cmp/cmp"
)

func Gzip(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gw := gzip.NewWriter(&b)
	if _, err := gw.Write(data); err != nil {
		return nil, fmt.Errorf("failed to gzip Write: %v", err)
	}
	if err := gw.Close(); err != nil {
		return nil, fmt.Errorf("failed to Close gzip Writer: %v", err)
	}

	return b.Bytes(), nil
}

var pool = &sync.Pool{
	New: func() interface{} {
		fmt.Println("called New")
		return &bytes.Buffer{}
	},
}

func GunzipWithBytesBufferPool(data []byte) ([]byte, error) {
	gr, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to gzip.NewReader: %v", err)
	}
	defer gr.Close()

	buf := pool.Get().(*bytes.Buffer)
	defer pool.Put(buf)
	buf.Reset()

	data, err = ioutil.ReadAll(gr)
	if err != nil {
		return nil, fmt.Errorf("failed to ReadAll: %v", err)
	}
	buf.Write(data)

	return buf.Bytes(), nil
}

func GunzipWithBytesBufferPool2(data []byte) ([]byte, error) {
	gr, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to gzip.NewReader: %v", err)
	}
	// r := pool.Get().(*bytes.Buffer) // ここは
	// fmt.Printf("buf %#v\n", r)

	buf := pool.Get().(*bytes.Buffer)
	defer pool.Put(buf)
	buf.Reset()
	fmt.Printf("buf in Gunzip: %#v\n", buf)

	if _, err := io.Copy(buf, gr); err != nil {
		return nil, fmt.Errorf("failed to io.Copy: %v", err)
	}

	if err := gr.Close(); err != nil {
		return nil, fmt.Errorf("failed to Close gzip Reader: %v", err)
	}

	return buf.Bytes(), nil

	// buf2 := *buf
	// // var buf2 bytes.Buffer
	// if _, err := io.Copy(&buf2, gr); err != nil {
	// 	return nil, fmt.Errorf("failed to io.Copy: %v", err)
	// }
	// if err := gr.Close(); err != nil {
	// 	return nil, fmt.Errorf("failed to Close gzip Reader: %v", err)
	// }

	// return buf2.Bytes(), nil
}

func GzipWithBytesBufferPool(data []byte) ([]byte, error) {
	buf := pool.Get().(*bytes.Buffer)
	defer pool.Put(buf)
	buf.Reset()

	fmt.Printf("\nbuf in Gzip: %#v\n", buf)
	fmt.Println("buf in Gzip: ", buf.Bytes())
	// defer func() {
	// 	fmt.Printf("2 buf in Gzip: %#v\n", buf.String())
	// }()

	gz := gzip.NewWriter(buf)
	if _, err := gz.Write(data); err != nil {
		return nil, fmt.Errorf("failed to gzip Write: %v", err)
	}
	if err := gz.Close(); err != nil {
		return nil, fmt.Errorf("failed to gzip Close: %v", err)
	}

	return buf.Bytes(), nil
}

func Gunzip(data []byte) ([]byte, error) {
	gr, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to gzip.NewReader: %v", err)
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, gr); err != nil {
		return nil, fmt.Errorf("failed to io.Copy: %v", err)
	}
	if err := gr.Close(); err != nil {
		return nil, fmt.Errorf("failed to Close gzip Reader: %v", err)
	}

	return buf.Bytes(), nil
}

// func main() {
// 	file, err := os.Create("empty.gz")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()

// 	data := `https://pkg.go.dev/compress/gzip
// 	Documentation
// 	Overview
// 	Package gzip implements reading and writing of gzip format compressed files, as specified in RFC 1952.`

// 	for i := 0; i < 3; i++ {
// 		fmt.Println("num i:", i)
// 		res, err := GzipWithBytesBufferPool([]byte(data))
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		// if _, err := io.Copy(file, bytes.NewReader(res)); err != nil {
// 		// 	log.Fatalf("failed to io.Copy: %v", err)
// 		// }

// 		// buf := pool.Get().(*bytes.Buffer)
// 		// fmt.Printf("buf: %#v\n", buf.String())
// 		// pool.Put(buf)

// 		r, err := GunzipWithBytesBufferPool2(res)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Println(string(r))
// 	}
// }

func main() {
	d := `https://pkg.go.dev/compress/gzip
	Documentation
	Overview
	Package gzip implements reading and writing of gzip format compressed files, as specified in RFC 1952.`
	data := []byte(d)

	// Buffer pool Gzip
	buf1 := pool.Get().(*bytes.Buffer)
	buf1.Reset()

	gz1 := gzip.NewWriter(buf1)
	if _, err := gz1.Write(data); err != nil {
		log.Fatalf("failed to gzip Write: %v", err)
	}
	if err := gz1.Close(); err != nil {
		log.Fatalf("failed to gzip Close: %v", err)
	}
	pool.Put(buf1)
	res1 := buf1.Bytes()

	// Buffer pool Gunzip
	gr2, err := gzip.NewReader(bytes.NewBuffer(res1))
	if err != nil {
		log.Fatalf("failed to gzip.NewReader: %v", err)
	}

	buf2 := pool.Get().(*bytes.Buffer)
	buf2.Reset()
	// fmt.Printf("buf2 in Gunzip: %#v\n", buf2)

	if _, err := io.Copy(buf2, gr2); err != nil {
		// if _, err := MYCopy(buf2, gr2); err != nil {
		log.Fatalf("failed to io.Copy: %v", err)
	}

	if err := gr2.Close(); err != nil {
		log.Fatalf("failed to Close gzip Reader: %v", err)
	}
	res2 := buf2.Bytes()

	gr2.Close()
	pool.Put(buf2)
	fmt.Println(string(res2))
	if diff := cmp.Diff(string(res2), string(data)); diff != "" {
		fmt.Printf("got: %v,want: %v, diff: %s\n", string(res2), string(data), diff)
	}
	if diff := cmp.Diff(res2, data); diff != "" {
		fmt.Printf("raw got: %v,want: %v, diff: %s\n", res2, data, diff)
	}
	if !bytes.Equal(res2, data) {
		fmt.Println("diff!")
	}

	// r := pool.Get().(*bytes.Buffer) // ここは
	// fmt.Printf("buf %#v\n", r)
	// r2 := pool.Get().(*bytes.Buffer) // ここは
	// fmt.Printf("buf %#v\n", r2)

	// Buffer pool Gzip
	buf3 := pool.Get().(*bytes.Buffer)
	// fmt.Printf("buf3 %#v\n", buf3)
	buf3.Reset()

	gz3 := gzip.NewWriter(buf3)
	if _, err := gz3.Write(data); err != nil {
		log.Fatalf("failed to gzip Write: %v", err)
	}
	if err := gz3.Close(); err != nil {
		log.Fatalf("failed to gzip Close: %v", err)
	}
	pool.Put(buf3)
	res3 := buf3.Bytes()

	// Buffer pool Gunzip
	gr4, err := gzip.NewReader(bytes.NewBuffer(res3))
	if err != nil {
		log.Fatalf("failed to gzip.NewReader: %v", err)
	}

	buf4 := pool.Get().(*bytes.Buffer)
	buf4.Reset()
	// r := pool.Get().(*bytes.Buffer) // ここは
	// fmt.Printf("buf %#v\n", r)
	// if !bytes.Equal(buf4.Bytes(), r.Bytes()) {
	// 	fmt.Println("koko diff!")
	// }
	// fmt.Printf("buf4 in Gunzip: %#v\n", buf4)

	if written, err := io.Copy(buf4, gr4); err != nil {
		// if _, err := MyCopy(buf4, gr4); err != nil {
		fmt.Println("written", written)
		log.Fatalf("failed to io.Copy: %v", err)
	}

	if err := gr4.Close(); err != nil {
		log.Fatalf("failed to Close gzip Reader: %v", err)
	}
	res4 := buf4.Bytes()

	gr4.Close()
	pool.Put(buf4)
	fmt.Println(string(res4))
	if diff := cmp.Diff(string(res4), string(data)); diff != "" {
		fmt.Printf("got: %v,want: %v, diff: %s\n", string(res4), string(data), diff)
	}
	if diff := cmp.Diff(res4, data); diff != "" {
		fmt.Printf("raw got: %v,want: %v, diff: %s\n", res4, data, diff)
	}
	if !bytes.Equal(res4, data) {
		fmt.Println("diff!")
	}

}

func fn1() {
	d := `https://pkg.go.dev/compress/gzip
	Documentation
	Overview
	Package gzip implements reading and writing of gzip format compressed files, as specified in RFC 1952.`
	data := []byte(d)

	// Buffer pool Gzip
	buf := pool.Get().(*bytes.Buffer)
	buf.Reset()

	gz := gzip.NewWriter(buf)
	if _, err := gz.Write(data); err != nil {
		log.Fatalf("failed to gzip Write: %v", err)
	}
	if err := gz.Close(); err != nil {
		log.Fatalf("failed to gzip Close: %v", err)
	}
	pool.Put(buf)
	res := buf.Bytes()

	// Buffer pool Gunzip
	gr, err := gzip.NewReader(bytes.NewBuffer(res))
	if err != nil {
		log.Fatalf("failed to gzip.NewReader: %v", err)
	}
	// r := pool.Get().(*bytes.Buffer) // ここは
	// fmt.Printf("buf %#v\n", r)

	buf2 := pool.Get().(*bytes.Buffer)
	buf2.Reset()
	fmt.Printf("buf2 in Gunzip: %#v\n", buf2)

	if _, err := io.Copy(buf2, gr); err != nil {
		log.Fatalf("failed to io.Copy: %v", err)
	}

	if err := gr.Close(); err != nil {
		log.Fatalf("failed to Close gzip Reader: %v", err)
	}
	pool.Put(buf2)
	res2 := buf2.Bytes()
	fmt.Println(string(res2))
	if diff := cmp.Diff(string(res2), string(data)); diff != "" {
		fmt.Printf("got: %v,want: %v, diff: %s\n", string(res2), string(data), diff)
	}
	if diff := cmp.Diff(res2, data); diff != "" {
		fmt.Printf("raw got: %v,want: %v, diff: %s\n", res2, data, diff)
	}
	if !bytes.Equal(res2, data) {
		fmt.Println("diff!")
	}

	// // 通常のGzip
	// var b3 bytes.Buffer
	// gw3 := gzip.NewWriter(&b3)
	// if _, err := gw3.Write(data); err != nil {
	// 	log.Fatalf("failed to gzip Write: %v", err)
	// }
	// if err := gw3.Close(); err != nil {
	// 	log.Fatalf("failed to Close gzip Writer: %v", err)
	// }
	// res3 := b3.Bytes()

	// Buffer pool Gzip
	buf3 := pool.Get().(*bytes.Buffer)
	buf3.Reset()

	// fmt.Printf("buf3 in Gunzip: %#v\n", buf3)

	// gz3 := gzip.NewWriter(buf3)
	gz3 := NewWriter(buf3)
	// gz3.Reset(buf3)
	// // 念のため再度チェック
	buf32 := &bytes.Buffer{}
	if !bytes.Equal(buf3.Bytes(), buf32.Bytes()) {
		fmt.Println("diff!")
	}
	// gz32 := gzip.NewWriter(buf32)
	// if !bytes.Equal(gz3.Header, buf32.Bytes()) {
	// 	fmt.Println("diff!")
	// }

	newRes2 := make([]byte, len(res2))
	copy(newRes2, res2)
	if !bytes.Equal(newRes2, res2) {
		fmt.Println("diff origin!")
	}

	fmt.Printf("buf before gz3.Write: %#v\n", buf3.Bytes())

	// // if _, err := gz3.Write(newRes2); err != nil { // これなら大丈夫
	if _, err := gz3.Write(res2); err != nil { // ここがなんか変
		// if _, err := gz3.Write(data); err != nil { // これならOK
		log.Fatalf("failed to gzip Write: %v", err)
	}
	// digest := uint32(0)
	// digest = crc32.Update(digest, crc32.IEEETable, res2)
	// fmt.Println("digest", digest)
	// gz3.Writeの前後でres2の中身が変わっている
	if !bytes.Equal(newRes2, res2) {
		fmt.Println("diff after!")
	}
	// if diff := cmp.Diff(string(newRes2), string(res2)); diff != "" {
	// 	fmt.Printf("got: %v,want: %v, diff: %s\n", string(newRes2), string(res2), diff)
	// }

	if err := gz3.Close(); err != nil {
		log.Fatalf("failed to gzip Close: %v", err)
	}
	pool.Put(buf3)
	res3 := buf3.Bytes()
	if diff := cmp.Diff(string(res3), string(res)); diff != "" {
		fmt.Printf("after res 3got: %v,want: %v, diff: %s\n", string(res3), string(res), diff)
	}
	// if !bytes.Equal(res3, res) {
	// 	fmt.Println("diff!")
	// }

	// // 通常のGunzip
	// gr4, err := gzip.NewReader(bytes.NewBuffer(res3))
	// if err != nil {
	// 	log.Fatalf("failed to gzip.NewReader: %v", err)
	// }
	// var buf4 bytes.Buffer
	// if _, err := io.Copy(&buf4, gr4); err != nil {
	// 	log.Fatalf("failed to io.Copy: %v", err)
	// }
	// if err := gr4.Close(); err != nil {
	// 	log.Fatalf("failed to Close gzip Reader: %v", err)
	// }
	// res4 := buf4.Bytes()
	// // fmt.Println(string(res4))
	// if diff := cmp.Diff(string(res4), string(data)); diff != "" {
	// 	fmt.Printf("got: %v,want: %v, diff: %s\n", string(res4), string(data), diff)
	// }

	// // Gunzip
	// gr4, err := gzip.NewReader(bytes.NewBuffer(res3))
	// if err != nil {
	// 	log.Fatalf("failed to gzip.NewReader: %v", err)
	// }
	// // r := pool.Get().(*bytes.Buffer) // ここは
	// // fmt.Printf("buf %#v\n", r)

	// buf4 := pool.Get().(*bytes.Buffer)
	// buf4.Reset()
	// fmt.Printf("buf4 in Gunzip: %#v\n", buf4)

	// if _, err := io.Copy(buf4, gr4); err != nil {
	// 	log.Fatalf("failed to io.Copy: %v", err)
	// }
	// // d4, err := ioutil.ReadAll(gr4)
	// // if err != nil {
	// // 	log.Fatalf("failed to ReadAll: %v", err)
	// // }
	// // buf4.Write(d4)

	// if err := gr4.Close(); err != nil {
	// 	log.Fatalf("failed to Close gzip Reader: %v", err)
	// }
	// pool.Put(buf4)
	// res4 := buf4.Bytes()
	// fmt.Println(string(res4))

	// if diff := cmp.Diff(string(res2), string(res4)); diff != "" {
	// 	fmt.Printf("got: %v,want: %v, diff: %s\n", string(res2), string(res4), diff)
	// }

	// for i := 0; i < 3; i++ {
	// 	var buf bytes.Buffer
	// 	zw := gzip.NewWriter(&buf)

	// 	_, err := zw.Write([]byte("A long time ago in a galaxy far, far away..."))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	if err := zw.Close(); err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	zr, err := gzip.NewReader(&buf)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	// fmt.Printf("Name: %s\nComment: %s\nModTime: %s\n\n", zr.Name, zr.Comment, zr.ModTime.UTC())

	// 	var buf2 bytes.Buffer
	// 	if _, err := io.Copy(&buf2, zr); err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	if err := zr.Close(); err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	fmt.Printf("buf2: %#v\n", buf2.String())
	// }

	// for i := 0; i < 3; i++ {
	// 	buf := &bytes.Buffer{}
	// 	buf.Reset()
	// 	zw := gzip.NewWriter(buf)

	// 	_, err := zw.Write([]byte("A long time ago in a galaxy far, far away..."))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	if err := zw.Close(); err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	zr, err := gzip.NewReader(buf)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	// fmt.Printf("Name: %s\nComment: %s\nModTime: %s\n\n", zr.Name, zr.Comment, zr.ModTime.UTC())
	// 	// buf.Reset()

	// 	if _, err := io.Copy(buf, zr); err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	if err := zr.Close(); err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	fmt.Printf("buf: %#v\n", buf.String())

	// 	// var buf2 bytes.Buffer
	// 	// if _, err := io.Copy(&buf2, zr); err != nil {
	// 	// 	log.Fatal(err)
	// 	// }

	// 	// if err := zr.Close(); err != nil {
	// 	// 	log.Fatal(err)
	// 	// }

	// 	// fmt.Printf("buf2: %#v\n", buf2.String())
	// }

	// // for i := 0; i < 3; i++ {
	// // buf2 := pool.Get().(*bytes.Buffer)
	// // fmt.Printf("buf2: %#v\n", buf2.String())

	// res, err := GzipWithBytesBufferPool([]byte(data))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // _ = pool.Get().(*bytes.Buffer) // ここだけやってもだめ

	// _, err = GunzipWithBytesBufferPool2(res)
	// if err != nil {
	// 	log.Fatal("GunzipWithBytesBufferPool2", err)
	// }

	// fmt.Printf("res %#v\n", res)

	// // _ = pool.Get().(*bytes.Buffer) // 1回目のGunzipのあとに必要

	// res2, err := GzipWithBytesBufferPool([]byte(data))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // _ = pool.Get().(*bytes.Buffer) // ここでも大丈夫
	// // fmt.Printf("buf2: %#v\n", buf2.String())

	// fmt.Printf("res2 %#v\n", res2)
	// // fmt.Println(reflect.DeepEqual(res, res2))
	// // fmt.Println("res==res2", string(res) == string(res2))
	// // _, err = GunzipWithBytesBufferPool2(res2)

	// _, err = GunzipWithBytesBufferPool2(res2)
	// if err != nil {
	// 	log.Fatal("GunzipWithBytesBufferPool2 2", err)
	// }
	// // fmt.Printf("got: %#v\n", string(got))

}

// These constants are copied from the flate package, so that code that imports
// "compress/gzip" does not also have to import "compress/flate".
const (
	// NoCompression      = flate.NoCompression
	BestSpeed          = flate.BestSpeed
	BestCompression    = flate.BestCompression
	DefaultCompression = flate.DefaultCompression
	HuffmanOnly        = flate.HuffmanOnly
)

const (
	gzipID1     = 0x1f
	gzipID2     = 0x8b
	gzipDeflate = 8
	// flagText    = 1 << 0
	// flagHdrCrc  = 1 << 1
	// flagExtra   = 1 << 2
	// flagName    = 1 << 3
	// flagComment = 1 << 4
)

// The gzip file stores a header giving metadata about the compressed file.
// That header is exposed as the fields of the Writer and Reader structs.
//
// Strings must be UTF-8 encoded and may only contain Unicode code points
// U+0001 through U+00FF, due to limitations of the GZIP file format.
type Header struct {
	Comment string    // comment
	Extra   []byte    // "extra data"
	ModTime time.Time // modification time
	Name    string    // file name
	OS      byte      // operating system type
}

var le = binary.LittleEndian

// A Writer is an io.WriteCloser.
// Writes to a Writer are compressed and written to w.
type Writer struct {
	Header      // written at first call to Write, Flush, or Close
	w           io.Writer
	level       int
	wroteHeader bool
	compressor  *flate.Writer
	digest      uint32 // CRC-32, IEEE polynomial (section 8)
	size        uint32 // Uncompressed size (section 2.3.1)
	closed      bool
	buf         [10]byte
	err         error
}

// NewWriter returns a new Writer.
// Writes to the returned writer are compressed and written to w.
//
// It is the caller's responsibility to call Close on the Writer when done.
// Writes may be buffered and not flushed until Close.
//
// Callers that wish to set the fields in Writer.Header must do so before
// the first call to Write, Flush, or Close.
func NewWriter(w io.Writer) *Writer {
	z, _ := NewWriterLevel(w, DefaultCompression)
	return z
}

// NewWriterLevel is like NewWriter but specifies the compression level instead
// of assuming DefaultCompression.

// The compression level can be DefaultCompression, NoCompression, HuffmanOnly
// or any integer value between BestSpeed and BestCompression inclusive.
// The error returned will be nil if the level is valid.
func NewWriterLevel(w io.Writer, level int) (*Writer, error) {
	if level < HuffmanOnly || level > BestCompression {
		return nil, fmt.Errorf("gzip: invalid compression level: %d", level)
	}
	z := new(Writer)
	z.init(w, level)
	return z, nil
}

func (z *Writer) init(w io.Writer, level int) {
	compressor := z.compressor
	if compressor != nil {
		compressor.Reset(w)
	}
	*z = Writer{
		Header: Header{
			OS: 255, // unknown
		},
		w:          w,
		level:      level,
		compressor: compressor,
	}
}

// Reset discards the Writer z's state and makes it equivalent to the
// result of its original state from NewWriter or NewWriterLevel, but
// writing to w instead. This permits reusing a Writer rather than
// allocating a new one.
func (z *Writer) Reset(w io.Writer) {
	z.init(w, z.level)
}

// writeBytes writes a length-prefixed byte slice to z.w.
func (z *Writer) writeBytes(b []byte) error {
	if len(b) > 0xffff {
		return errors.New("gzip.Write: Extra data is too large")
	}
	le.PutUint16(z.buf[:2], uint16(len(b)))
	_, err := z.w.Write(z.buf[:2])
	if err != nil {
		return err
	}
	_, err = z.w.Write(b)
	return err
}

// writeString writes a UTF-8 string s in GZIP's format to z.w.
// GZIP (RFC 1952) specifies that strings are NUL-terminated ISO 8859-1 (Latin-1).
func (z *Writer) writeString(s string) (err error) {
	// GZIP stores Latin-1 strings; error if non-Latin-1; convert if non-ASCII.
	needconv := false
	for _, v := range s {
		if v == 0 || v > 0xff {
			return errors.New("gzip.Write: non-Latin-1 header string")
		}
		if v > 0x7f {
			needconv = true
		}
	}
	if needconv {
		b := make([]byte, 0, len(s))
		for _, v := range s {
			b = append(b, byte(v))
		}
		_, err = z.w.Write(b)
	} else {
		_, err = io.WriteString(z.w, s)
	}
	if err != nil {
		return err
	}
	// GZIP strings are NUL-terminated.
	z.buf[0] = 0
	_, err = z.w.Write(z.buf[:1])
	return err
}

// Write writes a compressed form of p to the underlying io.Writer. The
// compressed bytes are not necessarily flushed until the Writer is closed.
func (z *Writer) Write(p []byte) (int, error) {
	fmt.Println("Write!!!!")

	if z.err != nil {
		return 0, z.err
	}

	p2 := make([]byte, len(p))
	copy(p2, p)
	if !bytes.Equal(p2, p) {
		fmt.Println("diff origin in Write!")
	}

	var n int
	// Write the GZIP header lazily.
	if !z.wroteHeader {
		z.wroteHeader = true
		z.buf = [10]byte{0: gzipID1, 1: gzipID2, 2: gzipDeflate}
		if z.Extra != nil {
			z.buf[3] |= 0x04
		}
		if z.Name != "" {
			z.buf[3] |= 0x08
		}
		if z.Comment != "" {
			z.buf[3] |= 0x10
		}
		if !bytes.Equal(p2, p) {
			fmt.Println("diff origin in Write1!")
		}
		if z.ModTime.After(time.Unix(0, 0)) {
			// Section 2.3.1, the zero value for MTIME means that the
			// modified time is not set.
			le.PutUint32(z.buf[4:8], uint32(z.ModTime.Unix()))
		}
		if z.level == BestCompression {
			z.buf[8] = 2
		} else if z.level == BestSpeed {
			z.buf[8] = 4
		}
		z.buf[9] = z.OS
		if !bytes.Equal(p2, p) {
			fmt.Println("diff origin in Write1-5!")
		}
		_, z.err = z.w.Write(z.buf[:10]) // ここ？
		fmt.Println("z.buf", z.buf, z.buf[:10])
		if !bytes.Equal(p2, p) {
			fmt.Println("diff origin in Write2!")
		}
		if diff := cmp.Diff(p2, p); diff != "" {
			fmt.Printf("got: %v,want: %v, diff: %s\n", p2, p, diff)
		}
		// if diff := cmp.Diff(string(p2), string(p)); diff != "" {
		// 	fmt.Printf("got: %v,want: %v, diff: %s\n", string(p2), string(p), diff)
		// }
		if z.err != nil {
			return 0, z.err
		}

		if z.Extra != nil {
			z.err = z.writeBytes(z.Extra)
			if z.err != nil {
				return 0, z.err
			}
		}
		if z.Name != "" {
			z.err = z.writeString(z.Name)
			if z.err != nil {
				return 0, z.err
			}
		}
		if z.Comment != "" {
			z.err = z.writeString(z.Comment)
			if z.err != nil {
				return 0, z.err
			}
		}
		if z.compressor == nil {
			z.compressor, _ = flate.NewWriter(z.w, z.level)
		}
	}
	if !bytes.Equal(p2, p) {
		fmt.Println("diff origin in Write3!")
	}

	z.size += uint32(len(p))
	z.digest = crc32.Update(z.digest, crc32.IEEETable, p)
	n, z.err = z.compressor.Write(p)

	fmt.Println("koko")
	if !bytes.Equal(p2, p) {
		fmt.Println("diff after in Write4!")
	}
	return n, z.err
}

// Close closes the Writer by flushing any unwritten data to the underlying
// io.Writer and writing the GZIP footer.
// It does not close the underlying io.Writer.
func (z *Writer) Close() error {
	if z.err != nil {
		return z.err
	}
	if z.closed {
		return nil
	}
	z.closed = true
	if !z.wroteHeader {
		z.Write(nil)
		if z.err != nil {
			return z.err
		}
	}
	z.err = z.compressor.Close()
	if z.err != nil {
		return z.err
	}
	le.PutUint32(z.buf[:4], z.digest)
	le.PutUint32(z.buf[4:8], z.size)
	_, z.err = z.w.Write(z.buf[:8])
	return z.err
}

// Implementations must not retain p.
type Reader interface {
	Read(p []byte) (n int, err error)
}

type IoWriter interface {
	Write(p []byte) (n int, err error)
}

type WriterTo interface {
	WriteTo(w IoWriter) (n int64, err error)
}

type ReaderFrom interface {
	ReadFrom(r Reader) (n int64, err error)
}

type LimitedReader struct {
	R Reader // underlying reader
	N int64  // max bytes remaining
}

func (l *LimitedReader) Read(p []byte) (n int, err error) {
	if l.N <= 0 {
		return 0, io.EOF
	}
	if int64(len(p)) > l.N {
		p = p[0:l.N]
	}
	n, err = l.R.Read(p)
	l.N -= int64(n)
	return
}

func MyCopy(dst IoWriter, src Reader) (written int64, err error) {
	fmt.Println("my Copy")
	return copyBuffer(dst, src, nil)
}

// copyBuffer is the actual implementation of Copy and CopyBuffer.
// if buf is nil, one is allocated.
func copyBuffer(dst IoWriter, src Reader, buf []byte) (written int64, err error) {
	// If the reader has a WriteTo method, use it to do the copy.
	// Avoids an allocation and a copy.
	if wt, ok := src.(WriterTo); ok {
		return wt.WriteTo(dst)
	}
	// Similarly, if the writer has a ReadFrom method, use it to do the copy.
	if rt, ok := dst.(ReaderFrom); ok {
		return rt.ReadFrom(src)
	}
	if buf == nil {
		size := 32 * 1024
		if l, ok := src.(*LimitedReader); ok && int64(size) > l.N {
			if l.N < 1 {
				size = 1
			} else {
				size = int(l.N)
			}
		}
		buf = make([]byte, size)
	}
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = ErrShortWrite
				break
			}
		}
		if er != nil {
			// if er != EOF {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}

var ErrShortWrite = errors.New("short write")
