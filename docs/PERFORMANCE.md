## Performance Improvement

The `static mode` of `routine` improves performance by more than `20%` compared to the `normal mode`.

The benchmark results are as follows:

### :computer:Normal Mode

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

### :rocket:Static Mode

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
