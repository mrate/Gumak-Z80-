# Builds gumak.dll and generates .lib file from gumak.def,
# run inside VS Developer Console.

Write-Host "Building GO lib..."
go build -o lib\gumak.dll -buildmode=c-shared .\gumak.go

Write-Host "Generating .lib file..."
lib /def:gumak.def /OUT:lib\gumak.lib /MACHINE:x64

Write-Host "Done"