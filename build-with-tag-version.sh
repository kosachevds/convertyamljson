VERSION=$(git describe --tags --abbrev=0) || VERSION="v0.0.0"
go build -ldflags="-X convertyamljson/cmd.Version=$VERSION"