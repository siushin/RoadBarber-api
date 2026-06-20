package handler

// 审核状态常量（避免在多个文件中重复定义）
const (
	models_AuditApproved = 1 // 通过
	models_AuditPending  = 2 // 待审核
	models_AuditRejected = 3 // 拒绝
)
