// Package geometry 提供均匀直线阵(ULA)的几何与阵因子内核。
//
// 阵元沿直线排列，θ 从阵列轴线方向(θ=0，端射)量起，
// 侧射方向为 θ=π/2，θ∈[0,π] 即可见区。
package geometry

// ElementType 描述单一阵元的辐射模式。
type ElementType string

const (
	// ElementIso 各向同性元因子，恒为 1。
	ElementIso ElementType = "iso"
	// ElementDipole 短偶极元因子 sinθ，θ 为相对轴线的角度。
	ElementDipole ElementType = "dipole"
)

// Params 是一次阵因子核算的完整输入。
type Params struct {
	// N 是阵元数，必须 >= 2。
	N int
	// D 是相邻阵元间距，与 Lambda 同一长度单位，必须 > 0。
	D float64
	// Lambda 是工作波长，必须 > 0。
	Lambda float64
	// Beta 是相邻阵元之间的线性相位梯度，单位弧度。
	// 正值使波束向 +θ 侧偏转，负值向 −θ 侧；侧射 β=0，
	// 端射 β=−kd。
	Beta float64
	// Element 是元因子类型，空串按 iso 处理。
	Element ElementType
}

// ElementKind 返回规范化后的元因子类型。
func (p Params) ElementKind() ElementType {
	if p.Element == ElementDipole {
		return ElementDipole
	}
	return ElementIso
}

// FromDegrees 返回一份 Beta 已从度转换到弧度的副本。
func (p Params) FromDegrees() Params {
	p.Beta = DegToRad(p.Beta)
	return p
}

// Copy 返回一份独立副本。
func (p Params) Copy() Params { return p }

// DefaultThetaSteps 是 θ 采样密度的默认值（见 sample 包）。
const DefaultThetaSteps = 361
