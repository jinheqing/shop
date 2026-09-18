import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

class SGSReportsPage extends StatelessWidget {
  const SGSReportsPage({super.key});

  final _reports = const [
    {'batch': 'MG-2026-003', 'product': '邦东古树饼', 'date': '2026-02-18', 'status': 'passed', 'file': 'SGS_MG_2026.pdf'},
    {'batch': 'BZ-2026-007', 'product': '班章单株', 'date': '2026-02-22', 'status': 'passed', 'file': 'SGS_BZ_2026.pdf'},
    {'batch': 'JM-2026-012', 'product': '景迈乔木', 'date': '2026-02-25', 'status': 'pending', 'file': null},
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('🔬 SGS Reports'), backgroundColor: Colors.green.shade800),
      body: ListView.separated(
        itemCount: _reports.length,
        separatorBuilder: (_, __) => const SizedBox(height: 8),
        itemBuilder: (_, i) {
          final r = _reports[i];
          return Card(margin: const EdgeInsets.symmetric(horizontal: 12), child: Padding(padding: const EdgeInsets.all(12),
            child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Row(children: [
                Expanded(child: Text(r['product'] as String, style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 16))),
                Container(padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                  decoration: BoxDecoration(color: r['status'] == 'passed' ? Colors.green.shade100 : Colors.orange.shade100, borderRadius: BorderRadius.circular(12)),
                  child: Text(r['status'] == 'passed' ? '✅ Passed' : '⏳ Testing', style: TextStyle(fontSize: 12, color: r['status'] == 'passed' ? Colors.green.shade800 : Colors.orange.shade800)),
                ),
              ]),
              const SizedBox(height: 6),
              Text('Batch: ${r['batch']} · Reported ${r['date']}', style: TextStyle(fontSize: 12, color: Colors.grey[600])),
              const SizedBox(height: 4),
              const Text('Pesticides · Heavy metals · Microbiology · Flavonoid profile', style: TextStyle(fontSize: 11, color: Colors.grey)),
              const SizedBox(height: 8),
              Row(children: [
                Expanded(child: OutlinedButton.icon(onPressed: () {}, icon: const Icon(Icons.remove_red_eye, size: 16), label: const Text('View Results'))),
                const SizedBox(width: 8),
                Expanded(child: ElevatedButton.icon(onPressed: r['file'] != null ? () {} : null, icon: const Icon(Icons.download, size: 16), label: const Text('PDF'))),
              ]),
            ])),
          );
        },
      ),
    );
  }
}
