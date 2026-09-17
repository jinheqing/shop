/// 核心服务：顾问/茶农用 livekit_client 原生 WebRTC 推流
/// 完全对齐设计文档 Section 4 Flutter 直播推流核心代码
class LiveService {
  final Room _room = Room();
  final String livekitUrl;

  LiveService({this.livekitUrl = 'wss://live.ourdomain.com'});

  /// Go Live — 设计文档的核心方法
  Future<void> startLive({required String token}) async {
    await _room.connect(livekitUrl, token);
    await _room.localParticipant?.setCameraEnabled(true);
    await _room.localParticipant?.setMicrophoneEnabled(true);
  }

  Future<void> stopLive() async => await _room.disconnect();

  /// 获取 LiveKit Host Token（从 Go 后端 POST /livekit/token）
  Future<String?> fetchHostToken(int liveRoomId, int staffId) async {
    // GET /api/v1/livekit/token 由 api_service 调
    return null; // 实际实现通过 Dio + JWT
  }
}

/// IM WebSocket 实时通道
class ChatWebSocketService {
  WebSocketChannel? _channel;
  final String wsUrl; // ws://localhost:8080/ws/im

  ChatWebSocketService({this.wsUrl = 'ws://localhost:8080/ws/im'});

  Future<void> connect({required String token, required int conversationId}) async {
    _channel = WebSocketChannel.connect(
      Uri.parse('$wsUrl?token=$token&conversation_id=$conversationId'),
    );
  }

  void send(String type, Map<String, dynamic> payload) {
    _channel?.sink.add(jsonEncode({'type': type, ...payload}));
  }

  Stream<Map<String, dynamic>> get messages =>
      _channel!.stream.map((e) => Map<String, dynamic>.from(jsonDecode(e as String)));

  Future<void> dispose() async { await _channel?.sink.close(); }
}

import 'dart:convert';
import 'package:web_socket_channel/web_socket_channel.dart';
