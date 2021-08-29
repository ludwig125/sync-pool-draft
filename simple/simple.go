package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	pool := &sync.Pool{ // <1> poolの定義
		New: func() interface{} { // Poolから最初にGetした時はこのNew関数が呼ばれる
			return &[]int{}
		},
	}

	// append 10
	l := pool.Get().(*[]int) // <2> poolから取得。[]int{}が取れる
	fmt.Println("got slice", *l)
	(*l) = append((*l), 10)
	fmt.Println("after append", *l)
	pool.Put(l) // <3> poolに戻す

	// append 20
	l = pool.Get().(*[]int) // <4> poolから取得。[]int{10}が取れる
	fmt.Println("got slice", *l)
	(*l) = append((*l), 20)
	fmt.Println("after append", *l)
	pool.Put(l) // poolに戻す

	// ガベージコレクションをしてpoolの中身を消す
	runtime.GC() // <5> GCをすると一次的なキャッシュのPoolの中身は消える

	// append 30
	l = pool.Get().(*[]int) // <6> poolから取得。[]int{}が取れる
	fmt.Println("got slice", *l)
	(*l) = append((*l), 30)
	fmt.Println("after append", *l)
	pool.Put(l) // poolに戻す
}
