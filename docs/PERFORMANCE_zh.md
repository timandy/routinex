## 性能提升

`routine` 的`静态模式`相比`普通模式`性能提升超过`20%`。

基准测试结果如下：

### :computer:普通模式

```text
> go test -bench='.' -run='^$' -a
...
pkg: github.com/timandy/routine
cpu: AMD Ryzen 7 8845HS w/ Radeon 780M Graphics
...
BenchmarkGoid-4                                 424359808                2.798 ns/op           0 B/op          0 allocs/op
BenchmarkThreadLocal-4                          23663457                48.39 ns/op            8 B/op          0 allocs/op
BenchmarkThreadLocalWithInitial-4               23303187                48.22 ns/op            8 B/op          0 allocs/op
BenchmarkInheritableThreadLocal-4               24854239                47.97 ns/op            8 B/op          0 allocs/op
BenchmarkInheritableThreadLocalWithInitial-4    21850201                48.64 ns/op            8 B/op          0 allocs/op
BenchmarkGohack-4                               596646250                1.920 ns/op           0 B/op          0 allocs/op
```

### :rocket:静态模式

```text
> go test -bench='.' -run='^$' -a -toolexec='routinex -v'
...
pkg: github.com/timandy/routine
cpu: AMD Ryzen 7 8845HS w/ Radeon 780M Graphics
...
BenchmarkGoid-4                                 544578021                2.197 ns/op           0 B/op          0 allocs/op
BenchmarkThreadLocal-4                          39003333                27.54 ns/op            8 B/op          0 allocs/op
BenchmarkThreadLocalWithInitial-4               37230658                28.75 ns/op            8 B/op          0 allocs/op
BenchmarkInheritableThreadLocal-4               42434164                28.43 ns/op            8 B/op          0 allocs/op
BenchmarkInheritableThreadLocalWithInitial-4    41751183                29.28 ns/op            8 B/op          0 allocs/op
BenchmarkGohack-4                               1000000000               1.045 ns/op           0 B/op          0 allocs/op
```
