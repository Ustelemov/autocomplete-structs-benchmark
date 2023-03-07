## Benchmark of search tree realizations which could be using as requests\prefixes storage in **autocomplete task**

Overviewed structs:

- Trie (basic realization)
- Radix trie (a.k.a patrica-tree)
- Double-array trie (minimal prefix tree + cedar)
- SlimTrie (radix trie realization)
- ART (radix trie realization)
- Ternary tree

Data using in benchmark: prefixes starts from "к" letter in sorted and not sorted formats (5_206_618 elements)
(some structs require sorted prefixes as input).

Golang's map using as reference.
Compare structs' results for sorted and unsorted data.

For keys using []bytes or string (depends on realization and struct), for values (mostly) int values.
Double-array and SlimTrie require to prebuild []keys and []values in lexicographical order.
Other structs can be filling iterative way.

Device parameters:

goos: linux
goarch: amd64
cpu: Intel(R) Core(TM) i3-4030U CPU @ 1.90GHz

Benchmark sorting full prefixes-elements slice (using sort.Strings()):

```
benchmark              iter       time/iter   bytes alloc        allocs
---------              ----       ---------   -----------        ------
BenchmarkSortInput-4      1   3493.14 ms/op       24 B/op   1 allocs/op
```

### InsertAll elements, heap in use after filling struct

BenchmarkMapInsertAll/unsorted inuse: 410.523438 MB
BenchmarkMapInsertAll/sorted inuse: 347.593750 MB

BenchmarkCedarInsertAll/unsorted inuse: 1746.015625 MB
BenchmarkCedarInsertAll/sorted inuse: 1746.015625 MB

BenchmarkSlimInsertAll/sorted inuse: 635.656250 MB

BenchmarkDoubleArrayInsertAll/sorted inuse: 541.507812 MB

BenchmarkTernaryTreeInsertAll/unsorted inuse: 982.843750 MB
BenchmarkTernaryTreeInsertAll/sorted inuse: 983.242188 MB

BenchmarkArtSuperfellInsertAll/unsorted inuse: 518.546875 MB
BenchmarkArtSuperfellInsertAll/sorted inuse: 512.140625 MB
BenchmarkArtBobotuInsertAll/unsorted inuse: 556.140625 MB
BenchmarkArtBobotuInsertAll/sorted inuse: 557.093750 MB
BenchmarkArtGammazeroInsertAll/unsorted inuse: 587.031250 MB
BenchmarkArtGammazeroInsertAll/sorted inuse: 588.539062 MB
BenchmarkArtRecoilmeInsertAll/unsorted inuse: 600.070312 MB
BenchmarkArtRecoilmeInsertAll/sorted inuse: 600.289062 MB
BenchmarkArtWenyxuTreeInsertAll/unsorted inuse: 643.773438 MB
BenchmarkArtWenyxuTreeInsertAll/sorted inuse: 644.132812 MB
**BenchmarkArtCurrentPlarInsertAll/unsorted inuse: 765.296875 MB**
**BenchmarkArtCurrentPlarInsertAll/sorted inuse: 775.109375 MB**
BenchmarkArtArriqaaqInsertAll/unsorted inuse: 1056.031250 MB
BenchmarkArtArriqaaqInsertAll/sorted inuse: 1057.226562 MB
BenchmarkArtKellydunnInsertAll/unsorted inuse: 1943.390625 MB
BenchmarkArtKellydunnInsertAll/sorted inuse: 1735.445312 MB

BenchmarkRadixGbrlsnchsInsertAll/unsorted inuse: 718.523438 MB
BenchmarkRadixGbrlsnchsInsertAll/sorted inuse: 723.929688 MB
BenchmarkRadixPatricaTreeInsertAll/unsorted inuse: 1449.695312 MB
BenchmarkRadixPatricaTreeInsertAll/sorted inuse: 1336.187500 MB
BenchmarkRadixWzshimingInsertAll/unsorted inuse: 6671.093750 MB
BenchmarkRadixWzshimingInsertAll/sorted inuse: 6671.085938 MB

BenchmarkTrieRyuanerinInsertAll/unsorted inuse: 700.570312 MB
BenchmarkTrieRyuanerinInsertAll/sorted inuse: 566.023438 MB
BenchmarkTrieBeevikInsertAll/unsorted inuse: 625.968750 MB
BenchmarkTrieBeevikInsertAll/sorted inuse: 596.273438 MB
BenchmarkTrieDghubbleInsertAll/unsorted inuse: 3225.976562 MB
BenchmarkTrieDghubbleInsertAll/sorted inuse: 3228.046875 MB
BenchmarkTrieAcomaguInsertAll/unsorted inuse: 2798.296875 MB
BenchmarkTrieAcomaguInsertAll/sorted inuse: 2798.312500 MB
BenchmarkTrieDerekparkerInsertAll/unsorted inuse: 5905.242188 MB
BenchmarkTrieDerekparkerInsertAll/sorted inuse: 5902.843750 MB

## InsertAll elements benchmark

benchmark                                       iter         time/iter       bytes alloc                allocs
---------                                       ----         ---------       -----------                ------
BenchmarkMapInsertAll/unsorted-4                   1     2274.07 ms/op    496102264 B/op      156828 allocs/op
BenchmarkMapInsertAll/sorted-4                     1     2534.52 ms/op    496100504 B/op      156821 allocs/op

BenchmarkCedarInsertAll/unsorted-4                 1    11647.43 ms/op   2441084480 B/op          76 allocs/op
BenchmarkCedarInsertAll/sorted-4                   1     5766.05 ms/op   2441084480 B/op          76 allocs/op

BenchmarkSlimInsertAll/sorted-4                    1     1540.55 ms/op    708170016 B/op     5206756 allocs/op

BenchmarkDoubleArrayInsertAll/sorted-4             1     4006.02 ms/op   1063273488 B/op     9029693 allocs/op
BenchmarkTernaryTreeInsertAll/unsorted-4           1    13242.66 ms/op   1072563344 B/op    21370165 allocs/op
BenchmarkTernaryTreeInsertAll/sorted-4             1    18911.16 ms/op   1072563488 B/op    21370166 allocs/op

BenchmarkArtGammazeroInsertAll/unsorted-4          1    18997.52 ms/op    641987984 B/op    17510852 allocs/op
BenchmarkArtGammazeroInsertAll/sorted-4            1     4526.48 ms/op    641987984 B/op    17510852 allocs/op
BenchmarkArtBobotuInsertAll/unsorted-4             1    12058.99 ms/op    586779384 B/op     8589669 allocs/op
BenchmarkArtBobotuInsertAll/sorted-4               1     4729.72 ms/op    586779352 B/op     8589669 allocs/op
**BenchmarkArtCurrentPlarInsertAll/unsorted-4        1    18724.40 ms/op    864498160 B/op    22397903 allocs/op**
**BenchmarkArtCurrentPlarInsertAll/sorted-4          1     6890.02 ms/op    863337456 B/op    22385934 allocs/op**
BenchmarkArtArriqaaqInsertAll/unsorted-4           1    31230.86 ms/op   1168427592 B/op    32842822 allocs/op
BenchmarkArtArriqaaqInsertAll/sorted-4             1     7754.50 ms/op   1168424632 B/op    32842822 allocs/op
BenchmarkArtWenyxuTreeInsertAll/unsorted-4         1    19219.43 ms/op    678224528 B/op    11874370 allocs/op
BenchmarkArtWenyxuTreeInsertAll/sorted-4           1     8913.02 ms/op    678224528 B/op    11874370 allocs/op
BenchmarkArtSuperfellInsertAll/unsorted-4          1    15027.74 ms/op    551044864 B/op     9092768 allocs/op
BenchmarkArtSuperfellInsertAll/sorted-4            1     9711.21 ms/op    544104944 B/op     9020476 allocs/op
BenchmarkArtKellydunnInsertAll/unsorted-4          1    34066.20 ms/op   2676054088 B/op   147694921 allocs/op
BenchmarkArtKellydunnInsertAll/sorted-4            1    14122.88 ms/op   2657841664 B/op   145418374 allocs/op
BenchmarkArtRecoilmeInsertAll/unsorted-4           1    34525.78 ms/op    629627200 B/op     9754190 allocs/op
BenchmarkArtRecoilmeInsertAll/sorted-4             1    34237.09 ms/op    629627344 B/op     9754191 allocs/op

BenchmarkRadixPatricaTreeInsertAll/unsorted-4      1    26386.79 ms/op   1624819272 B/op    27633560 allocs/op
BenchmarkRadixPatricaTreeInsertAll/sorted-4        1     6380.76 ms/op   1471543048 B/op    25219033 allocs/op
BenchmarkRadixGbrlsnchsInsertAll/unsorted-4        1    37982.05 ms/op    765939032 B/op    19883868 allocs/op
BenchmarkRadixGbrlsnchsInsertAll/sorted-4          1    13945.76 ms/op    772656568 B/op    20351393 allocs/op
BenchmarkRadixWzshimingInsertAll/unsorted-4        1   101502.45 ms/op   6993935456 B/op     9636841 allocs/op
BenchmarkRadixWzshimingInsertAll/sorted-4          1    43988.20 ms/op   6993935440 B/op     9636841 allocs/op

BenchmarkTrieRyuanerinInsertAll/unsorted-4         1    39758.98 ms/op    829361368 B/op    21361783 allocs/op
BenchmarkTrieRyuanerinInsertAll/sorted-4           1     8831.29 ms/op    747380624 B/op    18675381 allocs/op
BenchmarkTrieBeevikInsertAll/unsorted-4            1    37711.47 ms/op    781958928 B/op    14082547 allocs/op
BenchmarkTrieBeevikInsertAll/sorted-4              1     9137.16 ms/op    662094368 B/op    12304235 allocs/op
BenchmarkTrieDghubbleInsertAll/unsorted-4          1    24334.85 ms/op   3382533688 B/op    56727983 allocs/op
BenchmarkTrieDghubbleInsertAll/sorted-4            1    13549.48 ms/op   3382554080 B/op    56728166 allocs/op
BenchmarkTrieAcomaguInsertAll/unsorted-4           1    30124.79 ms/op   8011876104 B/op    36609603 allocs/op
BenchmarkTrieAcomaguInsertAll/sorted-4             1    14086.33 ms/op   8011876112 B/op    36609604 allocs/op
BenchmarkTrieDerekparkerInsertAll/unsorted-4       1    48751.45 ms/op   6209347088 B/op    73787622 allocs/op
BenchmarkTrieDerekparkerInsertAll/sorted-4         1    36356.23 ms/op   6209345344 B/op    73787607 allocs/op

## Insert one element benchmark

benchmark                                     iter       time/iter   bytes alloc        allocs
---------                                     ----       ---------   -----------        ------
BenchmarkMapInsert/unsorted-4              2272058    466.20 ns/op      109 B/op   0 allocs/op
BenchmarkMapInsert/sorted-4                2656876    447.90 ns/op       93 B/op   0 allocs/op

BenchmarkCedarInsert/unsorted-4             858540   2068.00 ns/op      710 B/op   0 allocs/op
BenchmarkCedarInsert/sorted-4              1365578    969.00 ns/op      446 B/op   0 allocs/op

BenchmarkArtGammazeroInsert/unsorted-4      617863   2812.00 ns/op      145 B/op   3 allocs/op
BenchmarkArtGammazeroInsert/sorted-4       1397941    778.90 ns/op      123 B/op   3 allocs/op
BenchmarkArtBobotuInsert/unsorted-4        1000000   1890.00 ns/op      119 B/op   1 allocs/op
BenchmarkArtBobotuInsert/sorted-4          1362096    895.30 ns/op      112 B/op   1 allocs/op
**BenchmarkArtCurrentPlarInsert/unsorted-4    947798   2946.00 ns/op      171 B/op   4 allocs/op**
**BenchmarkArtCurrentPlarInsert/sorted-4     1000000   1301.00 ns/op      162 B/op   4 allocs/op**
BenchmarkArtWenyxuInsert/unsorted-4         960289   2890.00 ns/op      139 B/op   2 allocs/op
BenchmarkArtWenyxuInsert/sorted-4          1000000   1473.00 ns/op      130 B/op   2 allocs/op
BenchmarkArtSuperfellInsert/unsorted-4     1000000   2247.00 ns/op      117 B/op   1 allocs/op
BenchmarkArtSuperfellInsert/sorted-4       1000000   1699.00 ns/op      103 B/op   1 allocs/op

BenchmarkRadixGbrlsnchsInsert/unsorted-4    878992   4101.00 ns/op      161 B/op   4 allocs/op
BenchmarkRadixGbrlsnchsInsert/sorted-4      952267   2296.00 ns/op      147 B/op   3 allocs/op

BenchmarkTrieBeevikInsert/unsorted-4        702032   4258.00 ns/op      191 B/op   3 allocs/op
BenchmarkTrieBeevikInsert/sorted-4         1000000   1396.00 ns/op      127 B/op   2 allocs/op
BenchmarkTrieRyuanerinInsert/unsorted-4     704898   4689.00 ns/op      177 B/op   4 allocs/op
BenchmarkTrieRyuanerinInsert/sorted-4      1000000   1549.00 ns/op      143 B/op   3 allocs/op

## Search one existing element benchmark

benchmark                                  iter       time/iter   bytes alloc        allocs
---------                                  ----       ---------   -----------        ------
BenchmarkMapExistsSearch-4              7512259    160.80 ns/op        0 B/op   0 allocs/op

BenchmarkCedarExistsSearch-4             685693   1484.00 ns/op        0 B/op   0 allocs/op

BenchmarkSlimExistsSearch-4             7288549    168.70 ns/op        0 B/op   0 allocs/op

BenchmarkDoubleArrayExistsSearch-4       873622   1292.00 ns/op        0 B/op   0 allocs/op

BenchmarkArtBobotuExistsSearch-4         618862   1939.00 ns/op        0 B/op   0 allocs/op
BenchmarkArtSuperfellExistsSearch-4      465730   2471.00 ns/op        0 B/op   0 allocs/op
BenchmarkArtGammazeroExistsSearch-4      369841   3414.00 ns/op        0 B/op   0 allocs/op
**BenchmarkArtCurrentPlarExistsSearch-4    343532   3542.00 ns/op        0 B/op   0 allocs/op**
BenchmarkArtWenyxuExistsSearch-4         303718   3705.00 ns/op        0 B/op   0 allocs/op
BenchmarkRadixGbrlsnchsExistsSearch-4    207688   5763.00 ns/op        0 B/op   0 allocs/op

BenchmarkTrieBeevikExistsSearch-4        251110   5108.00 ns/op        0 B/op   0 allocs/op
BenchmarkTrieRyuanerinExistsSearch-4     142078   8510.00 ns/op        0 B/op   0 allocs/op