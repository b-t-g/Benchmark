# Timescale Benchmark Tool

## How To run
The project should run simply with `docker-compose up --build` on the first run and then simply `docker-compose up` on subsequent runs.

For running just the CLI (independently of Docker Compose), there is an environment variable called `TIMESCALE_CONNECTION_STRING` which, if present, the program will use as the connection string by which to connect to the database. If not, the program will use a pre-defined connection string (and will log this to the database).

The required flags are:
1. `--query-params` or `-q` a required flag which specifies the file from which to generate queries.
2. `--number-of-workers` or `-n` an optional flag which specifies the number of worker threads to use, defaults to 4 workers.

## Tradeoffs
### Computational and Space Complexity
Improvements in computational complexity were measured against implementation complexity and readability; I deferred to the latter when the difference was not egregious (such as O(nlog(n)) vs O(n)) and the resulting code would either be harder to implement (correctly) or harder to read and modify.

My reasoning for this is primarily due to:
1. This tool primarily consists of making database queries and therefore, in most cases, the run time should be largely due to those queries. Not just due to network latency but also any I/O done by the database.
2. Without extensive profiling, it is possible to optimize not something in the critical path, which is time better spent improving other areas in the case.
3. The primary case for space considerations is my decision to pre-process the query params file, keep it and memory, this simplifies the implementation considerably and is elaborated upon below.

In particular, the median is computed via sorting all of the query durations and getting the median either from the middle (when the array has an odd number of elements) or by taking the average of the two "middle elements" when the array as an even number of elements. I considered the ["quickselect" algorithm](https://en.wikipedia.org/wiki/Quickselect), but decided against as there is no standard library package for it and implementing it myself would have complicated the implementation considerably.

While it is possible, in principle, to not read the entire query params file into memory, the issues of coordinating the main thread to get each worker thread a query when they're ready, while not reading too much of the file to get an appreciable space complexity benefit, while also ensuring that there are no threads being overly idle was a tradeoff that did not seem worth it.
