| Data Size                     | Algorithm + Compression Mode | n (iterations) | Time (ns/op)   | Speed (MB/s)   | Memory (B/op) | Allocations (allocs/op) |
|-------------------------------|------------------------------|----------------|----------------|----------------|---------------|--------------------------|
| **Small (~130 bytes)**        | gzip BestCompression         | 84 250         | 13 008         | 8.61           | 240           | 3                        |
|                               | gzip BestSpeed               | 865 723        | 1 384          | 80.90          | 498           | 4                        |
|                               | zstd SpeedBestCompression    | 248 517        | 4 752          | 23.57          | 112           | 1                        |
|                               | zstd SpeedBetterCompression  | 567 555        | 2 089          | 53.61          | 112           | 1                        |
|                               | zstd SpeedFastest            | 5 928 451      | 183.2          | 611.31         | 336           | 2                        |
| **Medium (~335 bytes)**       | gzip BestCompression         | 69 074         | 16 122         | 20.78          | 384           | 3                        |
|                               | gzip BestSpeed               | 261 393        | 4 701          | 71.27          | 886           | 4                        |
|                               | zstd SpeedBestCompression    | 132 940        | 9 096          | 36.83          | 352           | 1                        |
|                               | zstd SpeedBetterCompression  | 314 161        | 3 757          | 89.16          | 352           | 1                        |
|                               | zstd SpeedFastest            | 431 554        | 2 768          | 121.03         | 352           | 1                        |
| **Large Compressible (~1.15 GB)** | gzip BestCompression     | 4 957          | 242 034 000    | 537.11         | 2 135         | 5                        |
|                               | gzip BestSpeed               | 43 628         | 26 219 000     | 4 958.23       | 1 923         | 5                        |
|                               | zstd SpeedBestCompression    | 7 285          | 166 550        | 780.55         | 131 072       | 1                        |
|                               | zstd SpeedBetterCompression  | 79 058         | 12 804         | 10 152.82      | 131 072       | 1                        |
|                               | zstd SpeedFastest            | 84 102         | 12 812         | 10 146.44      | 131 072       | 1                        |