@REM 编译WEB界面
cd ./web/admin/assets/web
call web_build.bat
cd ../../../../
set file=./dist/fastsearch
echo "%file%"

@REM  编译Linux版本
SET CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -ldflags "-s -w -extldflags '-static'"  -o "%file%" . 
upx -6 %file%

@REM 编译WINDOWS版本
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-s -w -extldflags '-static'"  -o "%file%.exe" .
upx -6 %file%.exe

@REM 编译完成