pwd=`pwd`
cd ./web/admin/assets/web
# npm run build
cd $pwd
file=fastsearch
GOOS=linux go build -ldflags "-s -w -extldflags '-static'"  -o $file . 
upx -6 $file