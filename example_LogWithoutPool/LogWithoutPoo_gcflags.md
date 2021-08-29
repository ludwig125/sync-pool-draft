```
[~/go/src/github.com/ludwig125/sync-pool/example_LogWithoutPool] $go build -gcflags '-m'
# github.com/ludwig125/sync-pool/example_LogWithoutPool
./example_LogWithoutPool.go:12:6: can inline timeNow
./example_LogWithoutPool.go:13:18: inlining call to time.Unix
./example_LogWithoutPool.go:13:18: inlining call to time.unixTime
./example_LogWithoutPool.go:20:23: inlining call to timeNow
./example_LogWithoutPool.go:20:23: inlining call to time.Unix
./example_LogWithoutPool.go:20:23: inlining call to time.unixTime
./example_LogWithoutPool.go:20:29: inlining call to time.Time.UTC
./example_LogWithoutPool.go:20:29: inlining call to time.(*Time).setLoc
./example_LogWithoutPool.go:20:29: inlining call to time.(*Time).stripMono
./example_LogWithoutPool.go:20:29: inlining call to time.(*Time).sec
./example_LogWithoutPool.go:25:30: inlining call to bytes.(*Buffer).Bytes
./example_LogWithoutPool.go:30:6: can inline main
./example_LogWithoutPool.go:17:21: leaking param: w
./example_LogWithoutPool.go:17:34: key does not escape
./example_LogWithoutPool.go:17:39: val does not escape
./example_LogWithoutPool.go:18:7: &bytes.Buffer literal does not escape
./example_LogWithoutPool.go:26:12: ... argument does not escape
<autogenerated>:1: .this does not escape
[~/go/src/github.com/ludwig125/sync-pool/example_LogWithoutPool] $
```


```
[~/go/src/github.com/ludwig125/sync-pool/example_LogWithoutPool] $go build -gcflags '-m -m'
# github.com/ludwig125/sync-pool/example_LogWithoutPool
./example_LogWithoutPool.go:12:6: can inline timeNow with cost 55 as: func() time.Time { return time.Unix(1136214245, 0) }
./example_LogWithoutPool.go:13:18: inlining call to time.Unix func(int64, int64) time.Time { if time.nsec < int64(0) || time.nsec >= int64(1000000000) { var time.n·4 int64; time.n·4 = <N>; time.n·4 = time.nsec / int64(1000000000); time.sec += time.n·4; time.nsec -= time.n·4 * int64(1000000000); if time.nsec < int64(0) { time.nsec += int64(1000000000); time.sec-- } }; return time.unixTime(time.sec, int32(time.nsec)) }
./example_LogWithoutPool.go:13:18: inlining call to time.unixTime func(int64, int32) time.Time { return time.Time literal }
./example_LogWithoutPool.go:17:6: cannot inline LogWithoutPool: function too complex: cost 628 exceeds budget 80
./example_LogWithoutPool.go:20:23: inlining call to timeNow func() time.Time { return time.Unix(1136214245, 0) }
./example_LogWithoutPool.go:20:23: inlining call to time.Unix func(int64, int64) time.Time { if time.nsec < int64(0) || time.nsec >= int64(1000000000) { var time.n·4 int64; time.n·4 = <N>; time.n·4 = time.nsec / int64(1000000000); time.sec += time.n·4; time.nsec -= time.n·4 * int64(1000000000); if time.nsec < int64(0) { time.nsec += int64(1000000000); time.sec-- } }; return time.unixTime(time.sec, int32(time.nsec)) }
./example_LogWithoutPool.go:20:23: inlining call to time.unixTime func(int64, int32) time.Time { return time.Time literal }
./example_LogWithoutPool.go:20:29: inlining call to time.Time.UTC method(time.Time) func() time.Time { time.t.setLoc(&time.utcLoc); return time.t }
./example_LogWithoutPool.go:20:29: inlining call to time.(*Time).setLoc method(*time.Time) func(*time.Location) { if time.loc == &time.utcLoc { time.loc = nil }; time.t.stripMono(); time.t.loc = time.loc }
./example_LogWithoutPool.go:20:29: inlining call to time.(*Time).stripMono method(*time.Time) func() { if time.t.wall & uint64(9223372036854775808) != uint64(0) { time.t.ext = time.t.sec(); time.t.wall &= uint64(1073741823) } }
./example_LogWithoutPool.go:20:29: inlining call to time.(*Time).sec method(*time.Time) func() int64 { if time.t.wall & uint64(9223372036854775808) != uint64(0) { return int64(59453308800) + int64(time.t.wall << uint(1) >> uint(31)) }; return time.t.ext }
./example_LogWithoutPool.go:25:30: inlining call to bytes.(*Buffer).Bytes method(*bytes.Buffer) func() []byte { return bytes.b.buf[bytes.b.off:] }
./example_LogWithoutPool.go:30:6: can inline main with cost 63 as: func() { LogWithoutPool(os.Stdout, "path", "/search?q=flowers") }
./example_LogWithoutPool.go:17:21: parameter w leaks to {heap} with derefs=0:
./example_LogWithoutPool.go:17:21:   flow: {heap} = w:
./example_LogWithoutPool.go:17:21:     from w.Write(([]byte)(~R0)) (call parameter) at ./example_LogWithoutPool.go:25:22
./example_LogWithoutPool.go:17:21: leaking param: w
./example_LogWithoutPool.go:17:34: key does not escape
./example_LogWithoutPool.go:17:39: val does not escape
./example_LogWithoutPool.go:18:7: &bytes.Buffer literal does not escape
./example_LogWithoutPool.go:26:12: ... argument does not escape
<autogenerated>:1: .this does not escape
[~/go/src/github.com/ludwig125/sync-pool/example_LogWithoutPool] $
```