# array-af 是均匀直线阵的阵因子核算工具，计算主瓣指向、半功率波束宽度与可见区栅瓣

## 构建 / 运行 / 测试

```text
go build ./...
go run . -http :8080      # 打开 http://localhost:8080；算例见 example/broadside-8.json
go test ./...
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
