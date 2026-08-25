# array-af

array-af 是均匀直线阵（Uniform Linear Array）的阵因子核算工具。给定阵元数 N、阵元间距 d、工作波长 λ 与相邻阵元之间的线性相位梯度 β，它计算阵因子 AF(θ)=|sin(Nψ/2)/sin(ψ/2)|（ψ=kd·cosθ+β，k=2π/λ），并给出主瓣指向、半功率波束宽度（HPBW）、可见区内的栅瓣与方向性近似。输入边界：N 必须是 ≥2 的整数，d 与 λ 必须为正；非法输入一律返回带 error 字段的 JSON 错误体，不会返回静默错值。输出能力边界：方向性仅在半波间距侧射时近似为 N；主瓣（|ψ|=0 的点）落在可见区之外时如实报告「主瓣不可见」，不强行给角度。

## 启动

```bash
go run . -http :8080
```

打开 http://localhost:8080 进入交互页面，或直接用 curl 调用 API。

## API

两个求解接口均为 POST，请求与响应均为 JSON。

### POST /api/af

单次阵因子核算。请求：

```json
{
  "N": 8,
  "d": 0.5,
  "lambda": 1.0,
  "beta": 0.0,
  "element": "iso"
}
```

`beta` 为弧度（也可用 `beta_deg` 直接给度）；`element` 为 `iso`（各向同性，元因子恒为 1）或 `dipole`（短偶极 sinθ），缺省 `iso`。响应含 `mainlobe`（主瓣角与可见性）、`hpbw`（两个半功率点与宽度）、`grating`（栅瓣是否存在及角度）、`nulls`（第一对零点）、`directivity`、`af_peak`，以及 θ∈[0°,180°] 的采样点列 `points`（含阵因子、元因子、总方向图与归一化极坐标 x/y，供前端画图）。

### POST /api/scan

β 从起始值扫到终止值的主瓣表。请求：

```json
{
  "N": 8,
  "d": 0.5,
  "lambda": 1.0,
  "steps": 32
}
```

`beta_start_deg` 缺省 0（侧射），`beta_end_deg` 缺省 −kd（端射）。响应含逐行 `rows`（β、主瓣角、可见性、HPBW、是否栅瓣）与 `summary`（主瓣是否从侧射移向端射、HPBW 是否变宽、栅瓣是否出现）。

### GET /api/examples

返回内置算例列表与内容，供页面一键加载。

## 算例

- `example/broadside-8.json`：8 元半波侧射阵（N=8、d=λ/2、β=0）。主瓣在 90°，第一对零点在 arccos(±1/4)≈75.5°/104.5°，AF 峰值等于 N=8。

curl 验证：

```bash
curl -s -X POST http://localhost:8080/api/af \
  -H 'Content-Type: application/json' \
  --data-binary @example/broadside-8.json
```

非法输入示例（N=1）：

```bash
curl -s -X POST http://localhost:8080/api/af \
  -H 'Content-Type: application/json' \
  -d '{"N":1,"d":0.5,"lambda":1.0}'
```

返回 HTTP 400 与 `{"error":"invalid input: N: 1 violates must be an integer >= 2"}`。

## 关键约定

- **角度定义**：θ 从阵列轴线方向（θ=0，端射）量起，θ=π/2 为侧射，θ∈[0,π] 为可见区；空间相位差 ψ(θ)=kd·cosθ+β。β=0 时主瓣在侧射（90°），β=−kd 时主瓣在端射（0°）。
- **元因子**：默认各向同性（恒 1）；选择 `dipole` 时总方向图乘以 sinθ。主瓣、HPBW、栅瓣的判定基于阵因子 AF，元因子只调制总方向图曲线。
- **HPBW 扫描步长**：从主瓣角向两侧以 0.05° 步长扫描，找到阵因子首次降到峰值 1/√2 以下的点，再做线性插值取精确角度；两侧角度之差即 HPBW。端射束的另半侧在可见区外时，对应边界按可见区端点处理并标记 clipped。
- **栅瓣判定**：可见区（含端点）内存在 |ψ|=2π|m|（m≠0）的 θ 即判为有栅瓣。侧射半波阵（d=λ/2、β=0）无栅瓣；d 超过 λ 后栅瓣进入可见区。
- **方向性**：仅在半波间距且侧射时近似 D≈N，不引入与 N 无关的常数。

## 构建与测试

```bash
go build ./...
go test ./...
```

## 实现

领域内核在 `internal/`：`geometry`（校验、波数与 ψ、阵因子与元因子）、`beam`（主瓣、HPBW、栅瓣、零点、方向性）、`scan`（β 扫描）、`web`（HTTP 控制台与静态页托管）。`main.go` 只负责接线与资源嵌入。
