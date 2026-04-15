# 定义伪目标，不生成文件，直接执行
.PHONY: all build run gotool clean help

# 定义变量
BINARY="bluebell"

# 输入make后直接执行的命令
all: gotool build runexe

# 对应的shell命令
build:
	go build -o bluebell.exe

linux:
	go env -w CGO_ENABLED=0 GOOS=linux GOARCH=amd64
	go build -o bluebell
	go env -w CGO_ENABLED=1 GOOS=windows GOARCH=amd64

runexe:
	./${BINARY}

run:
	@go run ./

gotool:
	go fmt ./
	go vet ./

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

help:
	@echo "make - 格式化 Go 代码, 编译生成二进制文件, 并执行"
	@echo "make build - 编译 Go 代码, 生成二进制文件"
	@echo "make run - 直接运行 Go 代码"
	@echo "make clean - 移除二进制文件和 vim swap files"
	@echo "make gotool - 运行 Go 工具 'fmt' and 'vet'"
	@echo "make linux - 格式化 Go 代码, 编译生成linux二进制文件"