del /Q airi
del /Q airi_linux.tgz
set GOOS=linux
go build
tar czf airi_linux.tgz airi airi.toml.example
del /Q airi
