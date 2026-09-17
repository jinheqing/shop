class Staff {
  final int id;
  final String name;
  final String email;
  final String role; // admin / supervisor / advisor / tea_farmer / operations
  final String? token;
  final bool mfaEnabled;

  Staff({required this.id, required this.name, required this.email, required this.role, this.token, this.mfaEnabled = false});

  factory Staff.fromJson(Map<String, dynamic> j) => Staff(
    id: j['id'] ?? 0, name: j['name'] ?? '', email: j['email'] ?? '',
    role: j['role'] ?? '', token: j['access_token'] ?? j['token'], mfaEnabled: j['mfa_enabled'] ?? false,
  );
}
