import 'package:flutter/material.dart';
import 'dart:convert';
import 'package:web_socket_channel/web_socket_channel.dart';

class ChatPage extends StatefulWidget {
  final int conversationId;
  final String token;
  const ChatPage({super.key, required this.conversationId, required this.token});

  @override
  State<ChatPage> createState() => _ChatPageState();
}

class _ChatPageState extends State<ChatPage> {
  final _controller = TextEditingController();
  final _messages = <Map<String, dynamic>>[];
  WebSocketChannel? _channel;

  @override
  void initState() {
    super.initState();
    _connect();
  }

  void _connect() {
    try {
      _channel = WebSocketChannel.connect(
        Uri.parse('ws://localhost:8080/api/v1/ws/im?token=${widget.token}&conv_id=${widget.conversationId}'),
      );
      _channel!.stream.listen((event) {
        try {
          final m = jsonDecode(event as String) as Map<String, dynamic>;
          setState(() => _messages.add(m));
        } catch (_) {}
      });
    } catch (_) {
      // Fallback: load demo messages
      setState(() {
        _messages.addAll([
          {'content': '你好！我是曼岗村的王师傅 🍵', 'translated_content': 'Hello! I\'m Master Wang from Mangang Village', 'sender_role': 'advisor'},
          {'content': 'Hi Master Wang — how are the trees this season?', 'sender_role': 'user'},
          {'content': '今年雨水好，茶芽肥壮。预计3月底开始采摘 🌱', 'translated_content': 'Great rain this year, plump buds. Harvest starts late March', 'sender_role': 'advisor'},
        ]);
      });
    }
  }

  void _send() {
    final text = _controller.text.trim();
    if (text.isEmpty) return;
    final msg = {'content': text, 'sender_role': 'user', 'created_at': DateTime.now().toIso8601String()};
    setState(() { _messages.add(msg); _controller.clear(); });
    _channel?.sink.add(jsonEncode({
      'type': 'message', 'content': text, 'conv_id': widget.conversationId,
    }));
  }

  @override
  void dispose() { _controller.dispose(); _channel?.sink.close(); super.dispose(); }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('💬 Tea Advisor'), backgroundColor: Colors.green.shade800),
      body: Column(
        children: [
          Expanded(child: ListView.builder(
            padding: const EdgeInsets.all(12),
            itemCount: _messages.length,
            itemBuilder: (_, i) {
              final m = _messages[i];
              final mine = m['sender_role'] == 'user';
              return Align(
                alignment: mine ? Alignment.centerRight : Alignment.centerLeft,
                child: Container(
                  constraints: BoxConstraints(maxWidth: MediaQuery.of(context).size.width * 0.72),
                  margin: const EdgeInsets.symmetric(vertical: 4),
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: mine ? Colors.green.shade800 : Colors.grey.shade200,
                    borderRadius: BorderRadius.circular(16),
                  ),
                  child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                    Text(m['content'] ?? '', style: TextStyle(color: mine ? Colors.white : Colors.black87)),
                    if (m['translated_content'] != null && !mine) ...[
                      const SizedBox(height: 4),
                      Text('🇬🇧 ${m['translated_content']}', style: TextStyle(fontSize: 11, color: Colors.grey[600])),
                    ],
                  ]),
                ),
              );
            },
          )),
          Container(padding: const EdgeInsets.all(8), decoration: const BoxDecoration(color: Colors.white, border: Border(top: BorderSide(color: Colors.grey))),
            child: Row(children: [
              Expanded(child: TextField(controller: _controller, decoration: const InputDecoration(hintText: 'Type message...', border: OutlineInputBorder()))),
              const SizedBox(width: 8),
              IconButton(onPressed: _send, icon: const Icon(Icons.send, color: Colors.green), color: Colors.green.shade800),
            ])),
        ],
      ),
    );
  }
}
