BUILD_DIR := "build"
DEBUG_DIR := BUILD_DIR / "debug"
RELEASE_DIR := BUILD_DIR / "release"
BIN_DIR := "bin"
DIST_DIR := "dist"
VERSION := "0.5.0"

APP := "dnsck"
APP_BIN := if os() == "windows" { APP + ".exe" } else { APP }

# Default recipe (this list)
default:
    @echo "App = {{APP}}"
    @echo "Executable = {{APP_BIN}}"
    @echo "OS: {{os()}}, OS Family: {{os_family()}}, architecture: {{arch()}}"
    @just --list

# Delete binaries
clean:
    -rm -rf {{BUILD_DIR}}
    -rm -rf {{BIN_DIR}}
    -rm -rf {{DIST_DIR}}

# Quick run test of a release build (help and google.com)
run: release
    {{RELEASE_DIR}}/{{APP_BIN}} --version
    {{RELEASE_DIR}}/{{APP_BIN}} --help
    {{RELEASE_DIR}}/{{APP_BIN}} google.com
    {{RELEASE_DIR}}/{{APP_BIN}} tls-v1-0.badssl.com:1010
    {{RELEASE_DIR}}/{{APP_BIN}} expired.badssl.com www.burgaud.com burgaud.com go.dev

# Create the dist directory used to collect the packaged zip file for releases
dist:
    mkdir {{DIST_DIR}}

# Build app
build:
    go build -o {{DEBUG_DIR}}/{{APP_BIN}} main.go

# Build release version
release:
    go build -o {{RELEASE_DIR}}/{{APP}} -ldflags="-s -w" main.go
    -upx {{RELEASE_DIR}}/{{APP}}

# Cross compile versions for multiple target platforms
cross-compile: clean dist build-linux build-win build-mac build-mac-arm

# Cross compile for Linux
build-linux:
    GOOS=linux GOARCH=amd64 go build -o {{BIN_DIR}}/linux/{{APP_BIN}} main.go
    zip -j {{DIST_DIR}}/{{APP_BIN}}_linux-amd64_{{VERSION}}.zip {{BIN_DIR}}/linux/{{APP_BIN}}

# Cross compile for windows
build-win:
    GOOS=windows GOARCH=amd64 go build -o {{BIN_DIR}}/windows/{{APP}}.exe main.go
    zip -j {{DIST_DIR}}/{{APP}}_windows-amd64_{{VERSION}}.zip {{BIN_DIR}}/windows/{{APP}}.exe

# Cross compile for Mac amd64
build-mac:
    GOOS=darwin GOARCH=amd64 go build -o bin/macosx-amd64/{{APP_BIN}} main.go
    zip -j {{DIST_DIR}}/{{APP}}_macosx-amd64_{{VERSION}}.zip {{BIN_DIR}}/macosx-amd64/{{APP_BIN}}

# Cross compile for Mac arm64
build-mac-arm:
    GOOS=darwin GOARCH=arm64 go build -o bin/macosx-arm64/{{APP_BIN}} main.go
    zip -j {{DIST_DIR}}/{{APP}}_macosx-arm64_{{VERSION}}.zip {{BIN_DIR}}/macosx-arm64/{{APP_BIN}}

# Push and tag changes to github
push:
    git push
    git tag -a {{VERSION}} -m 'Version {{VERSION}}'
    git push origin --tags