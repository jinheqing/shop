package service

import (
	"context"

	"github.com/rs/zerolog/log"

	"tea-system/internal/models"
	"tea-system/internal/repository"
)

// stateTriggers — 订单-直播状态联动占位
// Step 9 的 order_state_machine 会调用这里；暂不接真实实现。
var (
	liveRoomRepo *repository.LiveRoomRepo
)

// InitStateTriggers — 注入依赖（在 main.go 初始化时调用一次）
func InitStateTriggers(lrRepo *repository.LiveRoomRepo) {
	liveRoomRepo = lrRepo
}

// OnOrderReadyForProduction — 订单进入 ready_for_production 时
// 自动创建一个 delivery_inspection 类型的直播间，绑定 order_id。
// 当前是占位 stub，只打日志；后续会接入真实仓库创建。
func OnOrderReadyForProduction(orderID uint64) {
	log.Info().Uint64("order_id", orderID).Msg("state_trigger: OnOrderReadyForProduction → would create delivery_inspection live_room (STUB)")

	if liveRoomRepo == nil {
		return
	}

	ctx := context.Background()
	// 先确认没有活跃的同订单 delivery_inspection 直播间
	existing, err := liveRoomRepo.GetActiveByOrderID(ctx, orderID)
	if err != nil {
		log.Warn().Err(err).Uint64("order_id", orderID).Msg("state_trigger: GetActiveByOrderID failed")
		return
	}
	if existing != nil {
		log.Info().Uint64("order_id", orderID).Uint64("room_id", existing.ID).Msg("state_trigger: active delivery_inspection room already exists")
		return
	}

	// TODO: 这里创建直播间，push_source=camera_rtmp，room_type=delivery_inspection
	// 后续接入 live_room_service.Create
	log.Debug().Uint64("order_id", orderID).Msg("state_trigger: stub — not creating room yet")
}

// OnLiveRoomEnded — 直播间结束时
// 如果 room_type=delivery_inspection 且 order_id!=nil → 推进订单状态
// 当前是占位 stub。
func OnLiveRoomEnded(roomID string) {
	log.Info().Str("room_id", roomID).Msg("state_trigger: OnLiveRoomEnded → would advance order state if delivery_inspection (STUB)")

	if liveRoomRepo == nil {
		return
	}

	ctx := context.Background()
	room, err := liveRoomRepo.GetByRoomID(ctx, roomID)
	if err != nil {
		log.Warn().Err(err).Str("room_id", roomID).Msg("state_trigger: GetByRoomID failed")
		return
	}
	if room.RoomType != models.RoomTypeDeliveryInspection || room.OrderID == nil {
		return
	}

	log.Info().
		Uint64("order_id", *room.OrderID).
		Str("room_id", room.RoomID).
		Msg("state_trigger: delivery_inspection ended → would call order_state_machine.AdvanceOrder (STUB)")
	// TODO: 调用 order_state_machine 推进状态（比如 ready_for_production → producing 或 inspecting）
}
