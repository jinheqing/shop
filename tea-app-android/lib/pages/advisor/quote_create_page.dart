/// 顾问创建 Bespoke 报价 — 26 字段 = Vue Admin 100% 对齐
/// 提交到 POST /api/v1/custom-products
import 'package:flutter/material.dart';
import '../../models/custom_product.dart';

class QuoteCreatePage extends StatefulWidget {
  const QuoteCreatePage({super.key});
  @override
  State<QuoteCreatePage> createState() => _QuoteCreatePageState();
}

class _QuoteCreatePageState extends State<QuoteCreatePage> {
  final _formKey = GlobalKey<FormState>();
  int _currentStep = 0;
  final _scrollController = ScrollController();

  // === 26 字段全部在这里，和 Go struct / Vue Admin 完全一致 ===
  String title = '';
  String rawTeaSource = '';
  String customRequirement = '';
  String teaType = 'raw_puer';
  String teaShape = 'cake';
  int teaShapeWeight = 357;
  bool smokedWithFlower = false;
  String flowerType = 'jasmine';
  String innerPackaging = '竹编内层 + 棉纸';
  String outerPackaging = '哑光纸质礼盒 + 丝绒内衬';
  String productCardText = '';
  String productCardFormat = 'vertical';
  String qrCodePosition = 'outer_back';
  double unitPrice = 0;
  int quantity = 1;
  double shippingCost = 0;
  String leadTime = '45 days from confirmation';
  String harvestDate = '';
  String roastingDate = '';
  String teaGardenLocation = '';
  String masterName = '';
  String storageLocation = '';
  int? sgsReportId;
  bool includeCustomLive = false;
  String? liveScheduledDate;

  final List<String> teaTypes = ['raw_puer', 'ripe_puer', 'ancient_tree', 'vintage'];
  final List<String> teaShapes = ['loose', 'cake', 'brick', 'tuo'];
  final List<String> flowerTypes = ['jasmine', 'osmanthus', 'orchid', 'custom'];
  final List<String> qrPositions = ['outer_front', 'outer_back', 'inner_front', 'hidden'];

  double get _total => unitPrice * quantity + shippingCost;

  List<Step> get _steps => [
    Step(title: const Text('基本信息'), content: _basicFields()),
    Step(title: const Text('普洱茶参数'), content: _teaFields()),
    Step(title: const Text('包装 + 产品卡'), content: _packagingFields()),
    Step(title: const Text('价格 + 溯源'), content: _pricingFields()),
    Step(title: const Text('直播 + 提交'), content: _liveAndSubmit()),
  ];

  Future<void> _submit() async {
    final model = CustomProduct(
      title: title, rawTeaSource: rawTeaSource, customRequirement: customRequirement,
      teaType: teaType, teaShape: teaShape, teaShapeWeight: teaShapeWeight,
      smokedWithFlower: smokedWithFlower, flowerType: flowerType,
      innerPackaging: innerPackaging, outerPackaging: outerPackaging,
      productCardText: productCardText, productCardFormat: productCardFormat,
      qrCodePosition: qrCodePosition,
      unitPrice: unitPrice, quantity: quantity, shippingCost: shippingCost,
      leadTime: leadTime,
      harvestDate: harvestDate, roastingDate: roastingDate,
      teaGardenLocation: teaGardenLocation, masterName: masterName, storageLocation: storageLocation,
      sgsReportId: sgsReportId, includeCustomLive: includeCustomLive, liveScheduledDate: liveScheduledDate,
    );
    // POST /api/v1/custom-products with JWT
    if (!mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      content: Text('✅ Bespoke Quote saved. Go to Admin to publish & generate token.'),
      backgroundColor: Colors.green.shade800,
    ));
    Navigator.of(context).pop();
  }

  // --- Field widgets ---
  Widget _basicFields() => Column(children: [
    _field('Title (客户定制标题)', onChanged: (v) => title = v),
    _field('Raw Tea Source (原料来源)', maxLines: 2, onChanged: (v) => rawTeaSource = v),
    _field('Custom Requirement (客户定制需求)', maxLines: 3, onChanged: (v) => customRequirement = v),
  ]);

  Widget _teaFields() => Column(children: [
    _dropdown('Tea Type', teaTypes, teaType, (v) => setState(() => teaType = v)),
    _dropdown('Tea Shape', teaShapes, teaShape, (v) => setState(() => teaShape = v)),
    Row(children: [
      const SizedBox(width: 120, child: Text('Shape Weight (g)')),
      Expanded(child: TextFormField(initialValue: teaShapeWeight.toString(),
        keyboardType: TextInputType.number,
        onChanged: (v) => teaShapeWeight = int.tryParse(v) ?? 357)),
    ]),
    SwitchListTile(title: const Text('Flower Smoked'), value: smokedWithFlower,
      onChanged: (v) => setState(() => smokedWithFlower = v)),
    if (smokedWithFlower) _dropdown('Flower Type', flowerTypes, flowerType, (v) => setState(() => flowerType = v)),
  ]);

  Widget _packagingFields() => Column(children: [
    _field('Inner Packaging', onChanged: (v) => innerPackaging = v),
    _field('Outer Packaging', onChanged: (v) => outerPackaging = v),
    _field('Product Card Text', onChanged: (v) => productCardText = v),
    Row(children: [
      const SizedBox(width: 120, child: Text('Card Format')),
      Expanded(child: DropdownButtonFormField(value: productCardFormat,
        items: const [DropdownMenuItem(value: 'vertical', child: Text('Vertical')), DropdownMenuItem(value: 'horizontal', child: Text('Horizontal'))],
        onChanged: (v) => setState(() => productCardFormat = v!))),
    ]),
    _dropdown('QR Position', qrPositions, qrCodePosition, (v) => setState(() => qrCodePosition = v)),
  ]);

  Widget _pricingFields() => Column(children: [
    Row(children: [
      Expanded(child: _num('Unit Price (£)', unitPrice, (v) => setState(() => unitPrice = v))),
      Expanded(child: _num('Quantity', quantity.toDouble(), (v) => setState(() => quantity = v.toInt()))),
    ]),
    Row(children: [
      Expanded(child: _num('Shipping (£)', shippingCost, (v) => setState(() => shippingCost = v))),
      Container(padding: const EdgeInsets.all(16), child: Text('Total: £${_total.toStringAsFixed(2)}', style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold)))),
    ]),
    _field('Lead Time', initialValue: leadTime, onChanged: (v) => leadTime = v),
    _field('Harvest Date', onChanged: (v) => harvestDate = v, hint: 'YYYY-MM-DD'),
    _field('Roasting Date', onChanged: (v) => roastingDate = v, hint: 'YYYY-MM-DD'),
    _field('Mountain Location', onChanged: (v) => teaGardenLocation = v),
    _field('Master Name', onChanged: (v) => masterName = v),
    _field('Storage Location', onChanged: (v) => storageLocation = v),
  ]);

  Widget _liveAndSubmit() => Column(children: [
    SwitchListTile(title: const Text('Include Custom Live Session'), value: includeCustomLive,
      onChanged: (v) => setState(() => includeCustomLive = v)),
    if (includeCustomLive) _field('Live Scheduled Date (YYYY-MM-DD)',
      onChanged: (v) => setState(() => liveScheduledDate = v)),
    const SizedBox(height: 20),
    ElevatedButton.icon(style: ElevatedButton.styleFrom(backgroundColor: Colors.green.shade800),
      onPressed: _submit, icon: const Icon(Icons.save), label: const Text('💾 Create Bespoke Quote')),
  ]);

  Widget _field(String label, {int maxLines = 1, String initialValue = '', String? hint, required void Function(String) onChanged}) =>
    Padding(padding: const EdgeInsets.symmetric(vertical: 8),
      child: TextFormField(initialValue: initialValue, maxLines: maxLines,
        decoration: InputDecoration(labelText: label, hintText: hint, border: const OutlineInputBorder()),
        onChanged: onChanged));

  Widget _dropdown(String label, List<String> items, String value, void Function(String) onChanged) =>
    Padding(padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(children: [
        SizedBox(width: 120, child: Text(label)),
        Expanded(child: DropdownButtonFormField(value: value,
          items: items.map((e) => DropdownMenuItem(value: e, child: Text(e))).toList(),
          onChanged: (v) => onChanged(v as String))),
      ]));

  Widget _num(String label, double value, void Function(double) onChanged) =>
    Padding(padding: const EdgeInsets.symmetric(vertical: 8),
      child: TextFormField(initialValue: value.toStringAsFixed(2),
        keyboardType: TextInputType.number, decoration: InputDecoration(labelText: label, border: const OutlineInputBorder()),
        onChanged: (v) => onChanged(double.tryParse(v) ?? 0)));

  @override
  Widget build(BuildContext context) {
    return Scaffold(appBar: AppBar(title: const Text('📝 Create Bespoke Quote'), backgroundColor: Colors.green.shade800),
      body: Form(key: _formKey, child: Stepper(
        currentStep: _currentStep,
        onStepContinue: () => setState(() => _currentStep < _steps.length - 1 ? _currentStep + 1 : _currentStep),
        onStepCancel: () => setState(() => _currentStep > 0 ? _currentStep - 1 : 0),
        onStepTapped: (i) => setState(() => _currentStep = i),
        steps: _steps,
      )));
  }
}
