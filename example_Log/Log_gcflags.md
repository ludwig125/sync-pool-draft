```
[~/go/src/github.com/ludwig125/sync-pool/example_Log] $go build -gcflags '-m'
# github.com/ludwig125/sync-pool/example_Log
./example_Log.go:24:6: can inline timeNow
./example_Log.go:25:18: inlining call to time.Unix
./example_Log.go:25:18: inlining call to time.unixTime
./example_Log.go:30:9: inlining call to bytes.(*Buffer).Reset
./example_Log.go:32:23: inlining call to timeNow
./example_Log.go:32:23: inlining call to time.Unix
./example_Log.go:32:23: inlining call to time.unixTime
./example_Log.go:32:29: inlining call to time.Time.UTC
./example_Log.go:32:29: inlining call to time.(*Time).setLoc
./example_Log.go:32:29: inlining call to time.(*Time).stripMono
./example_Log.go:32:29: inlining call to time.(*Time).sec
./example_Log.go:37:30: inlining call to bytes.(*Buffer).Bytes
./example_Log.go:44:6: can inline main
./example_Log.go:15:7: can inline glob..func1
./example_Log.go:30:9: Log ignoring self-assignment in bytes.b.buf = bytes.b.buf[:int(0)]
./example_Log.go:28:10: leaking param: w
./example_Log.go:28:23: key does not escape
./example_Log.go:28:28: val does not escape
./example_Log.go:39:12: ... argument does not escape
./example_Log.go:19:13: new(bytes.Buffer) escapes to heap
<autogenerated>:1: .this does not escape
[~/go/src/github.com/ludwig125/sync-pool/example_Log] $
```

```
[~/go/src/github.com/ludwig125/sync-pool/example_Log] $go build -gcflags '-m -m'
# github.com/ludwig125/sync-pool/example_Log
./example_Log.go:24:6: can inline timeNow with cost 55 as: func() time.Time { return time.Unix(1136214245, 0) }
./example_Log.go:25:18: inlining call to time.Unix func(int64, int64) time.Time { if time.nsec < int64(0) || time.nsec >= int64(1000000000) { var time.n·4 int64; time.n·4 = <N>; time.n·4 = time.nsec / int64(1000000000); time.sec += time.n·4; time.nsec -= time.n·4 * int64(1000000000); if time.nsec < int64(0) { time.nsec += int64(1000000000); time.sec-- } }; return time.unixTime(time.sec, int32(time.nsec)) }
./example_Log.go:25:18: inlining call to time.unixTime func(int64, int32) time.Time { return time.Time literal }
./example_Log.go:28:6: cannot inline Log: function too complex: cost 769 exceeds budget 80
./example_Log.go:30:9: inlining call to bytes.(*Buffer).Reset method(*bytes.Buffer) func() { bytes.b.buf = bytes.b.buf[:int(0)]; bytes.b.off = int(0); bytes.b.lastRead = bytes.readOp(0) }
./example_Log.go:32:23: inlining call to timeNow func() time.Time { return time.Unix(1136214245, 0) }
./example_Log.go:32:23: inlining call to time.Unix func(int64, int64) time.Time { if time.nsec < int64(0) || time.nsec >= int64(1000000000) { var time.n·4 int64; time.n·4 = <N>; time.n·4 = time.nsec / int64(1000000000); time.sec += time.n·4; time.nsec -= time.n·4 * int64(1000000000); if time.nsec < int64(0) { time.nsec += int64(1000000000); time.sec-- } }; return time.unixTime(time.sec, int32(time.nsec)) }
./example_Log.go:32:23: inlining call to time.unixTime func(int64, int32) time.Time { return time.Time literal }
./example_Log.go:32:29: inlining call to time.Time.UTC method(time.Time) func() time.Time { time.t.setLoc(&time.utcLoc); return time.t }
./example_Log.go:32:29: inlining call to time.(*Time).setLoc method(*time.Time) func(*time.Location) { if time.loc == &time.utcLoc { time.loc = nil }; time.t.stripMono(); time.t.loc = time.loc }
./example_Log.go:32:29: inlining call to time.(*Time).stripMono method(*time.Time) func() { if time.t.wall & uint64(9223372036854775808) != uint64(0) { time.t.ext = time.t.sec(); time.t.wall &= uint64(1073741823) } }
./example_Log.go:32:29: inlining call to time.(*Time).sec method(*time.Time) func() int64 { if time.t.wall & uint64(9223372036854775808) != uint64(0) { return int64(59453308800) + int64(time.t.wall << uint(1) >> uint(31)) }; return time.t.ext }
./example_Log.go:37:30: inlining call to bytes.(*Buffer).Bytes method(*bytes.Buffer) func() []byte { return bytes.b.buf[bytes.b.off:] }
./example_Log.go:44:6: can inline main with cost 63 as: func() { Log(os.Stdout, "path", "/search?q=flowers") }
./example_Log.go:15:7: can inline glob..func1 with cost 5 as: func() interface {} { return new(bytes.Buffer) }
./example_Log.go:30:9: Log ignoring self-assignment in bytes.b.buf = bytes.b.buf[:int(0)]
./example_Log.go:28:10: parameter w leaks to {heap} with derefs=0:
./example_Log.go:28:10:   flow: {heap} = w:
./example_Log.go:28:10:     from w.Write(([]byte)(~R0)) (call parameter) at ./example_Log.go:37:22
./example_Log.go:28:10: leaking param: w
./example_Log.go:28:23: key does not escape
./example_Log.go:28:28: val does not escape
./example_Log.go:39:12: ... argument does not escape
./example_Log.go:19:13: new(bytes.Buffer) escapes to heap:
./example_Log.go:19:13:   flow: ~r0 = &{storage for new(bytes.Buffer)}:
./example_Log.go:19:13:     from new(bytes.Buffer) (spill) at ./example_Log.go:19:13
./example_Log.go:19:13:     from new(bytes.Buffer) (interface-converted) at ./example_Log.go:19:13
./example_Log.go:19:13:     from return new(bytes.Buffer) (return) at ./example_Log.go:19:3
./example_Log.go:19:13: new(bytes.Buffer) escapes to heap
<autogenerated>:1: .this does not escape
[~/go/src/github.com/ludwig125/sync-pool/example_Log] $
```
