package service

import (
	"fmt"

	"tea-system/internal/models"
)

// OrderStateMachine — 订单状态机（基于 models.OrderState* 常量）
//
// 合法转移图：
//
//	ordering → paid, cancelled
//	paid → pending_declaration, producing, cancelled
//	pending_declaration → producing, pending_customs, cancelled
//	producing → ready_for_production, cancelled
//	ready_for_production → pending_declaration, pending_customs, cancelled
//	pending_customs → customs_clear, cancelled
//	customs_clear → shipped, cancelled
//	shipped → completed, disputed, cancelled
//	disputed → cancelled, completed
//
// 终态：completed / cancelled
// 关键约束：
//   - paid 不能直接跳到 completed（必须经过 生产 / 报关 / 清关 / 运输）
//   - ready_for_production 不能跳过 pending_declaration / pending_customs 直达 shipped
//   - shipped 才能进 completed
type OrderStateMachine struct{}

var stateMachine = &OrderStateMachine{}

// NewOrderStateMachine — 构造（目前无状态，返回单例即可）
func NewOrderStateMachine() *OrderStateMachine {
	return stateMachine
}

// 合法转移表：from → []allowed_to
var transitions = map[string][]string{
	models.OrderStateOrdering: {
		models.OrderStatePaid,
		models.OrderStateCancelled,
	},
	models.OrderStatePaid: {
		models.OrderStatePendingDeclaration,
		models.OrderStateProducing,
		models.OrderStateCancelled,
	},
	models.OrderStatePendingDeclaration: {
		models.OrderStateProducing,
		models.OrderStatePendingCustoms,
		models.OrderStateCancelled,
	},
	models.OrderStateProducing: {
		models.OrderStateReadyForProduction,
		models.OrderStateCancelled,
	},
	models.OrderStateReadyForProduction: {
		models.OrderStatePendingDeclaration,
		models.OrderStatePendingCustoms,
		models.OrderStateCancelled,
	},
	models.OrderStatePendingCustoms: {
		models.OrderStateCustomsClear,
		models.OrderStateCancelled,
	},
	models.OrderStateCustomsClear: {
		models.OrderStateShipped,
		models.OrderStateCancelled,
	},
	models.OrderStateShipped: {
		models.OrderStateCompleted,
		models.OrderStateDisputed,
		models.OrderStateCancelled,
	},
	models.OrderStateDisputed: {
		models.OrderStateCancelled,
		models.OrderStateCompleted,
	},
	// 终态不可再转移
	models.OrderStateCompleted: {},
	models.OrderStateCancelled: {},
}

// CanTransition — 校验 from → to 是否合法
func (sm *OrderStateMachine) CanTransition(from, to string) bool {
	if from == to {
		return false
	}
	allowed, ok := transitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

// IsTerminal — 判断是否终态
func (sm *OrderStateMachine) IsTerminal(state string) bool {
	return state == models.OrderStateCompleted || state == models.OrderStateCancelled
}

// AllowedTransitions — 返回当前状态允许的下一状态列表
func (sm *OrderStateMachine) AllowedTransitions(from string) []string {
	allowed, ok := transitions[from]
	if !ok {
		return nil
	}
	out := make([]string, len(allowed))
	copy(out, allowed)
	return out
}

// ValidateTransitionError — 转移不合法时的友好错误
type ValidateTransitionError struct {
	From string
	To   string
}

func (e *ValidateTransitionError) Error() string {
	return fmt.Sprintf("illegal order state transition: %s → %s", e.From, e.To)
}

// CheckTransition — 同 CanTransition，但失败时返回带上下文的 error
func (sm *OrderStateMachine) CheckTransition(from, to string) error {
	if sm.CanTransition(from, to) {
		return nil
	}
	return &ValidateTransitionError{From: from, To: to}
}
