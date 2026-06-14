pkill -f main.go
pkill -f go-build
sleep 1
rm -rf output.log
nohup go run cmd/main.go > output.log 2>&1 &
