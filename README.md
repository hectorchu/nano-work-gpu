nano-work-gpu
=============

This is a command-line app for generating Nano Proof-of-Work using the GPU.

Usage:

    -b int
          benchmark: number of iterations
    -d value
          difficulty: 8-byte hex string (default fffffff800000000)
    -r value
          root: 32-byte hex string

Example:

    nano-work-gpu -r 718CC2121C3E641059BC1C2CFC45666C99E8AE922F7A807B7D07B62C995D79E2

Generate a PoW for the hash `718CC2121C3E641059BC1C2CFC45666C99E8AE922F7A807B7D07B62C995D79E2` using the default network difficulty.

    nano-work-gpu -b 10

Run a benchmark with 10 iterations.

Building on Android
-------------------

The following environment variables are probably required:

    export CGO_LDFLAGS="-L/system/vendor/lib64 -Wl,-rpath -Wl,/system/vendor/lib64"
