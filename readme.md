#

# Windows 

add environment: CGO_ENABLED=1

# If got error: "cgo: C compiler "gcc" not found: exec: "gcc": executable file not found in %PATH%"

To install the GCC compiler:

Download and install 7-Zip.

Download [Mingw-w64](https://github.com/niXman/mingw-builds-binaries/releases/tag/14.2.0-rt_v12-rev0) . You want the file with “posix-seh” in its name.

Extract the archive with 7-Zip, then move the mingw64 directory to the root of your C:\ drive.

C:\MINGW64
    ├───bin
    ├───etc
    ├───include
    ├───lib
    ├───libexec
    ├───licenses
    ├───opt
    ├───share
    └───x86_64-w64-mingw32
Update the PATH environment variable

Right click on the Start menu icon, choose Run, then paste:

rundll32.exe sysdm.cpl,EditEnvironmentVariables

Under “User variables” select PATH, press Edit, press New, then paste:

C:\mingw64\bin

Verify that gcc is installed and in your PATH. Open a cmd window then type:

gcc --version

Now compile the extended version of Hugo. Change to the directory containing the Hugo codebase, then type:

go install -tags extended

Be patient; it will take a while.

Verify the location of the executable by typing where hugo. This should display something like:

C:\Users\joe\go\bin\hugo.exe

Verify the version by typing hugo version. This should display something like:

hugo v0.106.0-DEV-52ea07d2eb9c581015f85c0d2d6973640a62281b+extended windows/amd64 BuildDate=2022-11-01T17:45:34Z

# Start command: gotbot.exe -tg-bot-token 'TOKEN'

# Linux

Need install for linux:
    apt-get install build-essential

# Start command: gotbot -tg-bot-token 'TOKEN'