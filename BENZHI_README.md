Go 实现的均匀直线阵阵因子分析服务，默认启动 HTTP 在 :8080 提供主瓣指向、半功率波束宽度与栅瓣计算接口及交互式前端页面。

## 构建与启动

```bash
go build -o array-af .
./array-af              # 启动 HTTP 服务 :8080，打开 http://localhost:8080
```

## 评测镜像

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh array-af
```
