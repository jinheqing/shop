package models

// AllModels — 业务库所有 model（不含 audit_logs，那个是独立 PostgreSQL 实例）
var AllModels = []interface{}{
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
	&LiveRoom{},
	&SgsReport{},
	&Node{},
	&TranslationSession{},
	&SiteContent{},
}

// AuditModels — 审计库 model（独立 PostgreSQL 实例）
var AuditModels = []interface{}{
	&AuditLog{},
}
