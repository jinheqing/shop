import 'package:flutter/material.dart';

class SlowPresetsPage extends StatefulWidget {
  const SlowPresetsPage({super.key});

  @override
  State<SlowPresetsPage> createState() => _SlowPresetsPageState();
}

class _SlowPresetsPageState extends State<SlowPresetsPage> {
  final _presets = <Map<String, dynamic>>[
    {'id': 1, 'name': '云南省临沧市临翔区邦东乡曼岗村茶园', 'location': '云南临沧', 'status': 'live', 'camera_rtmp_url': 'rtmp://192.168.1.10/live/iceland', 'start_time': '2026-03-01 06:00'},
    {'id': 2, 'name': '班章春茶林', 'location': '云南勐海', 'status': 'live', 'camera_rtmp_url': 'rtmp://192.168.1.11/live/banzhang', 'start_time': '2026-03-15 05:30'},
    {'id': 3, 'name': '云南省普洱市澜沧拉祜族自治县惠民镇景迈村茶园', 'location': '云南澜沧', 'status': 'idle', 'camera_rtmp_url': 'rtmp://192.168.1.12/live/jingmai', 'start_time': null},
  ];

  void _toggleStatus(int i) {
    setState(() {
      _presets[i]['status'] = _presets[i]['status'] == 'live' ? 'idle' : 'live';
    });
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(_presets[i]['status'] == 'live' ? '▶ Started streaming' : '⏸ Stopped')));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('🌱 24/7 Slow Cameras'), backgroundColor: Colors.green.shade800),
      body: ListView.separated(
        itemCount: _presets.length,
        separatorBuilder: (_, __) => const Divider(height: 1),
        itemBuilder: (_, i) {
          final p = _presets[i];
          return ListTile(
            leading: CircleAvatar(backgroundColor: p['status'] == 'live' ? Colors.green : Colors.grey.shade300,
              child: Icon(p['status'] == 'live' ? Icons.videocam : Icons.videocam_off, color: Colors.white)),
            title: Text(p['name'], style: const TextStyle(fontWeight: FontWeight.w600)),
            subtitle: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Text('${p['location']}', style: const TextStyle(fontSize: 12)),
              if (p['start_time'] != null) Text('⏱ ${p['start_time']}', style: TextStyle(fontSize: 11, color: Colors.grey[600])),
              const SizedBox(height: 2),
              Text(p['camera_rtmp_url'], style: const TextStyle(fontSize: 10, fontFamily: 'monospace', color: Colors.grey)),
            ]),
            trailing: SizedBox(width: 110, child: Row(children: [
              Expanded(child: ElevatedButton.icon(
                onPressed: () => _toggleStatus(i),
                icon: Icon(p['status'] == 'live' ? Icons.stop : Icons.play_arrow, size: 16),
                label: Text(p['status'] == 'live' ? 'Stop' : 'Go'),
                style: ElevatedButton.styleFrom(backgroundColor: p['status'] == 'live' ? Colors.red : Colors.green),
              )),
            ])),
          );
        },
      ),
      floatingActionButton: FloatingActionButton.extended(onPressed: () {}, backgroundColor: Colors.green.shade800,
        icon: const Icon(Icons.add), label: const Text('Add Camera')),
    );
  }
}
