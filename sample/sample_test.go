package main

import (
	"testing"
)

func TestGzip(t *testing.T) {
	data := `https://pkg.go.dev/compress/gzip
Documentation
Overview
Package gzip implements reading and writing of gzip format compressed files, as specified in RFC 1952.`

	// Poolを正しく使わないと前にPutした値をGetで取ってきてしまうミスがあり得る
	// そのため、２回実行しても同じ結果であることを確認している
	for i := 0; i < 3; i++ {
		t.Run("Gzip_and_Gunzip", func(t *testing.T) {
			res, err := GzipWithBytesBufferPool([]byte(data))
			if err != nil {
				t.Fatal(err)
			}
			got, err := GunzipWithBytesBufferPool(res)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != data {
				t.Errorf("got: %s, want: %s", string(got), data)
			}
		})
		t.Run("Gzip_and_Gunzip2", func(t *testing.T) {
			res, err := GzipWithBytesBufferPool([]byte(data))
			if err != nil {
				t.Fatal(err)
			}
			got, err := GunzipWithBytesBufferPool2(res)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != data {
				t.Errorf("got: %s, want: %s", string(got), data)
			}
		})
	}
	// for i := 0; i < 3; i++ {
	// 	t.Run("Gzip_and_Gunzip", func(t *testing.T) {
	// 		res, err := GzipWithBytesBufferPool([]byte(data))
	// 		if err != nil {
	// 			t.Fatal(err)
	// 		}
	// 		got, err := Gunzip(res)
	// 		if err != nil {
	// 			t.Fatal(err)
	// 		}
	// 		if string(got) != data {
	// 			t.Errorf("got: %s, want: %s", string(got), data)
	// 		}
	// 	})
	// 	t.Run("Gzip_and_Gunzip2", func(t *testing.T) {
	// 		res, err := GzipWithBytesBufferPool([]byte(data))
	// 		if err != nil {
	// 			t.Fatal(err)
	// 		}
	// 		got, err := Gunzip(res)
	// 		if err != nil {
	// 			t.Fatal(err)
	// 		}
	// 		if string(got) != data {
	// 			t.Errorf("got: %s, want: %s", string(got), data)
	// 		}
	// 	})
	// }
}
