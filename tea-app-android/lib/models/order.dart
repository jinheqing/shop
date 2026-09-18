class Order {
  final int? id;
  final String orderNo;
  final int userId;
  final int staffId;
  final int customProductId;
  final Map<String, dynamic>? customProductSnapshot;
  final String state;
  final double unitPrice;
  final int quantity;
  final double shippingCost;
  final double totalAmount;
  final String hsCode;
  final String countryOfOrigin;
  final Map<String, dynamic>? billingAddressSnapshot;
  final Map<String, dynamic>? deliveryAddressSnapshot;
  final int? liveRoomId;
  final String? createdAt;

  Order({
    this.id, required this.orderNo, required this.userId, required this.staffId,
    required this.customProductId, this.customProductSnapshot,
    required this.state, required this.unitPrice, required this.quantity,
    required this.shippingCost, required this.totalAmount,
    required this.hsCode, required this.countryOfOrigin,
    this.billingAddressSnapshot, this.deliveryAddressSnapshot,
    this.liveRoomId, this.createdAt,
  });

  factory Order.fromJson(Map<String, dynamic> j) => Order(
    id: j['id'], orderNo: j['order_no'] ?? '', userId: j['user_id'] ?? 0,
    staffId: j['staff_id'] ?? 0, customProductId: j['custom_product_id'] ?? 0,
    customProductSnapshot: j['custom_product_snapshot'], state: j['state'] ?? '',
    unitPrice: (j['unit_price'] ?? 0).toDouble(), quantity: j['quantity'] ?? 1,
    shippingCost: (j['shipping_cost'] ?? 0).toDouble(), totalAmount: (j['total_amount'] ?? 0).toDouble(),
    hsCode: j['hs_code'] ?? '0902.20', countryOfOrigin: j['country_of_origin'] ?? 'China',
    billingAddressSnapshot: j['billing_address_snapshot'],
    deliveryAddressSnapshot: j['delivery_address_snapshot'],
    liveRoomId: j['live_room_id'], createdAt: j['created_at'],
  );
}
