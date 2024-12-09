build:
    go build main.go -o opc
scp-debug :build
    scp -P 10001 opc abc@zks.today:~/
    ssh abc@zks.today -p 10001 './opc'