import 'package:json_annotation/json_annotation.dart';

/// 100% 对齐 Go struct custom_products — 26 字段用户输入 + 自动生成字段
@JsonSerializable()
class CustomProduct {
  final int? id;
  final String? product_token;      // 自动生成 64 位
  final int? version;               // 自动版本号
  final bool? is_bespoke;
  final bool? non_refundable;

  // === 用户输入 25 字段 ===
  final String title;                           // 1
  @JsonKey(name: 'raw_tea_source') final String rawTeaSource;  // 2
  @JsonKey(name: 'custom_requirement') final String customRequirement; // 3
  @JsonKey(name: 'tea_type') final String teaType;             // 4: raw_puer / ripe_puer / ancient_tree / vintage
  @JsonKey(name: 'tea_shape') final String teaShape;           // 5: loose / cake / brick / tuo
  @JsonKey(name: 'tea_shape_weight') final int teaShapeWeight; // 6
  @JsonKey(name: 'smoked_with_flower') final bool smokedWithFlower; // 7
  @JsonKey(name: 'flower_type') final String? flowerType;      // 8
  @JsonKey(name: 'inner_packaging') final String innerPackaging;  // 9
  @JsonKey(name: 'outer_packaging') final String outerPackaging; // 10
  @JsonKey(name: 'product_card_text') final String? productCardText; // 11
  @JsonKey(name: 'product_card_format') final String productCardFormat; // 12
  @JsonKey(name: 'qr_code_position') final String qrCodePosition; // 13
  final double unitPrice;                                      // 14
  final int quantity;                                          // 15
  final double shippingCost;                                   // 16
  final String leadTime;                                       // 17
  @JsonKey(name: 'harvest_date') final String harvestDate;     // 18
  @JsonKey(name: 'roasting_date') final String roastingDate;   // 19
  @JsonKey(name: 'tea_garden_location') final String teaGardenLocation; // 20
  @JsonKey(name: 'master_name') final String masterName;       // 21
  @JsonKey(name: 'storage_location') final String storageLocation; // 22
  @JsonKey(name: 'sgs_report_id') final int? sgsReportId;      // 23
  final bool includeCustomLive;                                // 24
  final String? liveScheduledDate;                             // 25

  final double? totalAmount;
  final String? status;  // draft / published / archived
  final String? sku;     // auto-gen

  CustomProduct({
    this.id, this.product_token, this.version = 1, this.is_bespoke = true, this.non_refundable = true,
    required this.title, required this.rawTeaSource, required this.customRequirement,
    required this.teaType, required this.teaShape, this.teaShapeWeight = 357,
    this.smokedWithFlower = false, this.flowerType,
    required this.innerPackaging, required this.outerPackaging,
    this.productCardText, this.productCardFormat = 'vertical',
    this.qrCodePosition = 'outer_back',
    required this.unitPrice, required this.quantity, this.shippingCost = 0,
    required this.leadTime,
    required this.harvestDate, required this.roastingDate,
    required this.teaGardenLocation, required this.masterName, required this.storageLocation,
    this.sgsReportId, this.includeCustomLive = false, this.liveScheduledDate,
    this.totalAmount, this.status, this.sku,
  });

  Map<String, dynamic> toJson() => {
    'title': title, 'raw_tea_source': rawTeaSource, 'custom_requirement': customRequirement,
    'tea_type': teaType, 'tea_shape': teaShape, 'tea_shape_weight': teaShapeWeight,
    'smoked_with_flower': smokedWithFlower, 'flower_type': flowerType,
    'inner_packaging': innerPackaging, 'outer_packaging': outerPackaging,
    'product_card_text': productCardText, 'product_card_format': productCardFormat,
    'qr_code_position': qrCodePosition,
    'unit_price': unitPrice, 'quantity': quantity, 'shipping_cost': shippingCost,
    'lead_time': leadTime,
    'harvest_date': harvestDate, 'roasting_date': roastingDate,
    'tea_garden_location': teaGardenLocation, 'master_name': masterName,
    'storage_location': storageLocation, 'sgs_report_id': sgsReportId,
    'include_custom_live': includeCustomLive, 'live_scheduled_date': liveScheduledDate,
  };
}
