package models

// AllModels — 业务库所有 model（不含 audit_logs，那个是独立 PostgreSQL 实例）
// 注意顺序：有外键的表必须排在被引用表之后
var AllModels = []interface{}{
	// ===== 无外键引用的基础表 =====
	&User{},
	&Staff{},
	&Conversation{},
	&ConversationParticipant{},
	&Message{},
	&CustomProduct{},
	&Order{},
	&PaymentTransaction{},
	&Invoice{},
	&DeclarationLedger{},
	&ForeignExchangeLedger{},
	&SgsReport{},
	&Node{},
	&TranslationSession{},
	&SiteContent{},
	&DSARRequest{},
	&CookieConsentLog{},
	&UserGroup{},
	&UserGroupMember{},
	&ShortLink{},
	&Recording{},
	&Referral{},

	// ===== 有外键引用的表放最后 =====
	&LiveRoom{},
}

// AuditModels — 审计库 model（独立 PostgreSQL 实例）
var AuditModels = []interface{}{
	&AuditLog{},
}
