source /etc/profile
rm -rf bin/*
buildlist=("main" "cmdmain" "toolmain" "othermain")
for item in ${buildlist[*]}
do
 echo $item
 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o "bin/$tem"   "$item.go"
done
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/report_v2 report_v2_main.go