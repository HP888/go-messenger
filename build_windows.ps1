# This build script should be used when your system is Windows otherwise this may cause some problems when you want to build another project for other system.
# Set environment variable to windows
$env:GOOS = "windows"

# Build application
go build