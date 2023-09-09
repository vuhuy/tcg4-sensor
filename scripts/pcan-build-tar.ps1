# Build and tar the PCAN-USB client example on and for Windows x64.
# Tested with Node.js 18.17.1 and Nexe 4.0.0-rc.2.

# Go to example directory
$previousPwd = $PWD
Set-Location -ErrorAction Stop -LiteralPath $PSScriptRoot
Set-Location ..\examples\pcan-client-example\

# Remove old builds
$outputFile = "..\..\build\pcan-client-example.exe"

if (Test-Path $outputFile) {
    Remove-Item $outputFile
}

# Run nexe
$command = "nexe index.js --build --resource node_modules\@csllc\cs-pcan-usb\binding\Release\node-v108-win32-x64\cs_pcan_usb.node\ --resource node_modules\@csllc\cs-pcan-usb\binding\Release\node-v108-win32-x64\PCANBasic.dll --resource index.html --output ..\..\build\pcan-client-example.exe"
Invoke-Expression $command

# Tar compiled binary and dependencies
# 4.0.0-rc.2 does not return a zero exit code on success
if (Test-Path $outputFile) {
    $tarCommand = "tar -czvf ..\..\build\pcan-client-example-win.tar.gz README.md -C ..\..\build\  pcan-client-example.exe -C ..\ LICENSE -C examples\pcan-client-example\ node_modules\@csllc\cs-pcan-usb\binding\Release\node-v108-win32-x64\*"
    
    Invoke-Expression $tarCommand
    Remove-Item $outputFile
    Write-Host "Build and tar completed!"
}
else {
    Write-Host "Error: build failed"
}

# Restore current working directory
$previousPwd | Set-Location