import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:tea_app_android/services/api.dart';

class LivePage extends StatefulWidget {
  final String room;
  const LivePage({super.key, required this.room});
  @override
  State<LivePage> createState() => _LivePageState();
}

class _LivePageState extends State<LivePage> {
  String? token;
  String? rtmpUrl;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final d = await ApiService().api!.post('/livekit/token-for-obs',
        data: {'room_name': widget.room, 'identity': 'mobile-advisor'});
      if (mounted) setState(() {
        token = d.data['token'];
        rtmpUrl = d.data['obs_rtmp_url'];
      });
    } catch {}
  }

  @override
  Widget build(ctx) => Scaffold(
    appBar: AppBar(title: Text('📷 ${widget.room}'), leading: BackButton(onPressed: () => ctx.pop())),
    body: Padding(padding: const EdgeInsets.all(16), child: Column(crossAxisAlignment: CrossAxisAlignment.stretch, children: [
      Container(height: 220, color: Colors.black87, alignment: Alignment.center,
        child: const Column(mainAxisSize: MainAxisSize.min, children: [Icon(Icons.videocam, color: Colors.white38, size: 64), Text('LiveKit Video Feed', style: TextStyle(color: Colors.white54))])),
      const SizedBox(height: 20),
      Card(child: Padding(padding: const EdgeInsets.all(12), child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        const Text('OBS Stream URL', style: TextStyle(fontWeight: FontWeight.w600)),
        SelectableText(rtmpUrl ?? 'rtmp://localhost:1935/...', style: const TextStyle(fontFamily: 'monospace', fontSize: 11)),
        const SizedBox(height: 8),
        const Text('LiveKit Token (OBS Broadcast)', style: TextStyle(fontWeight: FontWeight.w600)),
        SelectableText(token ?? '—', style: const TextStyle(fontFamily: 'monospace', fontSize: 11)),
      ]))),
      const Spacer(),
      FilledButton.icon(icon: const Icon(Icons.broadcast_on_home), label: const Text('Start Streaming'), onPressed: token != null ? () {} : null),
    ])),
  );
}
