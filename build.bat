@echo off
setlocal EnableExtensions EnableDelayedExpansion

set "APP_NAME=examtopics-downloader"
set "ENTRY=.\cmd\main.go"
set "DIST_DIR=dist"
set "ARCH=%~1"
set "COMPRESS=%~2"

if "%ARCH%"=="" set "ARCH=amd64"
if /I "%COMPRESS%"=="" set "COMPRESS=auto"

call :banner

if /I not "%ARCH%"=="amd64" if /I not "%ARCH%"=="arm64" (
  call :error "Unsupported architecture: %ARCH% (use amd64 or arm64)"
  exit /b 1
)

call :step "Checking Go..."
where go >nul 2>&1
if errorlevel 1 (
  call :error "Go is not installed or not in PATH."
  exit /b 1
)
for /f "delims=" %%v in ('go version') do set "GO_VERSION=%%v"
echo     !GO_VERSION!

call :step "Preparing output directory..."
if not exist "%DIST_DIR%" mkdir "%DIST_DIR%"
set "OUT_FILE=%DIST_DIR%\%APP_NAME%-windows-%ARCH%.exe"
if exist "%OUT_FILE%" del /q "%OUT_FILE%" >nul 2>&1

call :step "Building %APP_NAME%-windows-%ARCH%.exe ..."
set "GOOS=windows"
set "GOARCH=%ARCH%"
set "CGO_ENABLED=0"
go build -trimpath -ldflags "-s -w -buildid=" -o "%OUT_FILE%" "%ENTRY%"
if errorlevel 1 (
  call :error "Build failed."
  exit /b 1
)

for %%F in ("%OUT_FILE%") do set "SIZE_BEFORE=%%~zF"

set "UPX_USED=0"
if /I "%COMPRESS%"=="off" (
  call :step "Compression disabled (arg: off)."
) else (
  call :step "Checking optional UPX compression..."
  where upx >nul 2>&1
  if errorlevel 1 (
    echo     UPX not found. Skipping compression.
  ) else (
    upx --best --lzma "%OUT_FILE%" >nul
    if errorlevel 1 (
      echo     UPX found but compression failed. Keeping uncompressed binary.
    ) else (
      echo     UPX compression complete.
      set "UPX_USED=1"
    )
  )
)

for %%F in ("%OUT_FILE%") do set "SIZE_AFTER=%%~zF"
set /a "KB_AFTER=(SIZE_AFTER+1023)/1024"

if "!UPX_USED!"=="1" (
  set /a "SAVED=SIZE_BEFORE-SIZE_AFTER"
  if !SIZE_BEFORE! GTR 0 (
    set /a "PCT=(SAVED*100)/SIZE_BEFORE"
  ) else (
    set "PCT=0"
  )
)

for %%I in ("%OUT_FILE%") do set "OUT_ABS=%%~fI"

echo.
echo ===============================================================
echo Build completed successfully.
echo Target      : windows/%ARCH%
echo Output      : !OUT_ABS!
echo Final size  : !KB_AFTER! KB
if "!UPX_USED!"=="1" (
  echo Size saved  : !PCT!%% ^(!SAVED! bytes^)
)
echo ===============================================================
echo.
echo Usage:
echo   build.bat [amd64^|arm64] [auto^|off]
echo.
echo Share this file:
echo   %OUT_FILE%
echo.

exit /b 0

:banner
echo ===============================================================
echo ExamTopics Downloader - Windows Builder
echo ===============================================================
echo.
exit /b 0

:step
echo [*] %~1
exit /b 0

:error
echo [X] %~1
exit /b 0
