import 'package:flutter/material.dart';

class FarmerOrderDetailPage extends StatelessWidget {
  final int orderId;
  const FarmerOrderDetailPage({super.key, required this.orderId});

  final _timeline = const ['received', 'producing', 'roasting', 'ready_for_delivery', 'pending_inspection', 'inspected'];

  @override
  Widget build(BuildContext context) {
    const current = 2;
    return Scaffold(
      appBar: AppBar(title: Text('Order #$orderId'), backgroundColor: Colors.green.shade800),
      body: ListView(children: [
        Container(color: Colors.green.shade50, padding: const EdgeInsets.all(16), child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          const Text('ORD-20260917-398535', style: TextStyle(fontFamily: 'monospace', fontSize: 13, color: Colors.green)),
          const SizedBox(height: 8),
          const Text('冰岛古树饼 · 357g × 12', style: TextStyle(fontSize: 20, fontWeight: FontWeight.w600)),
          const SizedBox(height: 4),
          const Text('Customer: Mr. James L. · London, UK', style: TextStyle(color: Colors.grey)),
        ])),
        _section('📅 Progress Timeline', Column(children: [
          for (var i = 0; i < _timeline.length; i++)
            ListTile(
              leading: Icon(i <= current ? Icons.check_circle : Icons.radio_button_unchecked, color: i <= current ? Colors.green : Colors.grey),
              title: Text(_timeline[i]),
              trailing: i == current ? const Chip(label: Text('CURRENT')) : null,
            ),
        ])),
        _section('🍵 Production Checklist', Column(children: [
          _checkbox('Picked from Iceland Old Village (GPS verified)', true),
          _checkbox('Withered for 12h at 24°C', true),
          _checkbox('Rolled by Master Wang 8x', true),
          _checkbox('Sun-dried 3 days', false),
          _checkbox('Steamed & pressed 357g', false),
          _checkbox('Stored in raw ageing facility 7d', false),
        ])),
        _section('📷 Live Camera Check', Column(children: [
          Row(children: [Expanded(child: Container(height: 140, color: Colors.black, alignment: Alignment.center, child: const Column(mainAxisSize: MainAxisSize.min, children: [
            Icon(Icons.videocam, color: Colors.green[400], size: 32),
            SizedBox(height: 4),
            Text('LIVE · 冰岛晒场', style: TextStyle(color: Colors.green[400])),
          ])))]),
          const SizedBox(height: 8),
          ElevatedButton.icon(onPressed: () {}, icon: const Icon(Icons.videocam), label: const Text('Open Live Camera'),
            style: ElevatedButton.styleFrom(backgroundColor: Colors.green.shade800)),
        ])),
      ]),
    );
  }

  Widget _section(String title, Widget body) {
    return Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
      Padding(padding: const EdgeInsets.fromLTRB(16, 20, 16, 8), child: Text(title, style: const TextStyle(fontSize: 14, fontWeight: FontWeight.w700, color: Colors.black54))),
      Padding(padding: const EdgeInsets.symmetric(horizontal: 16), child: Card(child: Padding(padding: const EdgeInsets.all(4), child: body))),
    ]);
  }

  Widget _checkbox(String label, bool value) {
    return ListTile(leading: Icon(value ? Icons.check_circle : Icons.radio_button_unchecked, color: value ? Colors.green : Colors.grey),
      title: Text(label, style: TextStyle(fontSize: 13, decoration: value ? TextDecoration.lineThrough : null, color: value ? Colors.grey : null)));
  }
}
