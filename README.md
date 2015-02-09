AurSir Runtime
==============

Install from source
-------------------

Linux:

- Install ZeroMQ v4.0.X from zeromq.org
- clone the source into GOPATH/src/
- go install aursirrt/src

Windows

- Get the Ruby DevKit (32 or 64 bit, depending on what version of Go you installed, file names here are for 64bit) from http://rubyinstaller.org/downloads . It gives you a fully working gcc environment
- Get the ZeroMQ sources
- in the source folder, go to builds\mingw32 and open a command promt
- start the devkitvars.bat (It's in the root folder where you installed DevKit)
- type "make all -f Makefile.mingw32"
- after a short time you get a lot of files. Important are "libzmq.dll" and "libzmq.dll.a"
- copy "libzmq.dll" and "libgcc_s_sjlj-1.dll" and "libstdc++-6.dll" either somewhere into your PATH or next to the binary that you will build later. The later two DLLs are found in "\mingw\bin"
- copy "libzmq.dll.a" into "\mingw\x86_64-w64-mingw32\lib"
- copy "zmq.h" and "zmq_utils.h" from "\include" to "\mingw\x86_64-w64-mingw32\include"
- go install aursirrt/src

if the install fails because it cannot find zmq.h, try the following

- $env:C_INCLUDE_PATH = "C:\rubydevkit\mingw\x86_64-w64-mingw32\include"
- $env:LIBRARY_PATH = "C:\rubydevkit\mingw\x86_64-w64-mingw32\lib"

Replace the paths with your own.