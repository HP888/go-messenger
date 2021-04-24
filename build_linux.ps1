# This build script should be used when your system is Windows otherwise this may cause some problems when you want to build another project for other system.
# Set environment variable to linux
$env:GOOS = "linux"

# Build application
go build

# Set environment variable back to windows
$env:GOOS = "windows"