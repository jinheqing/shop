class LiveRoom {
  final String roomId;
  final String roomName;
  final String roomType; // slow_preset / obs_tasting / open_calendar / customer_request / delivery_inspection / custom_private
  final String pushSource; // camera_rtmp / app_webrtc / obs_rtmp
  final String? obsRtmpUrl;
  final String? obsRtmpKey;
  final String status;    // configuring / scheduled / live / offline / ended
  final String? livekitToken;
  final int? peakViewers;

  LiveRoom({
    required this.roomId, required this.roomName, required this.roomType,
    required this.pushSource, this.obsRtmpUrl, this.obsRtmpKey,
    required this.status, this.livekitToken, this.peakViewers,
  });

  factory LiveRoom.fromJson(Map<String, dynamic> j) => LiveRoom(
    roomId: j['room_id'] ?? '', roomName: j['room_name'] ?? '',
    roomType: j['room_type'] ?? '', pushSource: j['push_source'] ?? '',
    obsRtmpUrl: j['obs_rtmp_url'], obsRtmpKey: j['obs_rtmp_key'],
    status: j['status'] ?? '', livekitToken: j['livekit_token_for_host'],
    peakViewers: j['peak_viewers'],
  );
}
