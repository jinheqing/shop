import 'package:flutter/material.dart';

class UserLoginPage extends StatefulWidget {
  const UserLoginPage({super.key});

  @override
  State<UserLoginPage> createState() => _UserLoginPageState();
}

class _UserLoginPageState extends State<UserLoginPage> {
  int _step = 0; // 0: enter email, 1: sent, 2: verify
  final _emailCtrl = TextEditingController();
  final _tokenCtrl = TextEditingController();
  String _savedEmail = '';
  bool _loading = false;

  Future<void> _requestLink() async {
    setState(() { _loading = true; _savedEmail = _emailCtrl.text; });
    // POST /api/v1/user/magic-link/request
    await Future.delayed(const Duration(milliseconds: 500));
    setState(() { _loading = false; _step = 1; });
  }

  Future<void> _verify() async {
    setState(() => _loading = true);
    await Future.delayed(const Duration(milliseconds: 500));
    setState(() { _loading = false; _step = 2; });
    // Navigate to home if on success
    if (mounted) Navigator.of(context).pop(true);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.green.shade50,
      body: SafeArea(child: SingleChildScrollView(child: Padding(padding: const EdgeInsets.all(24), child: Column(children: [
        const SizedBox(height: 60),
        const Text('🍃', style: TextStyle(fontSize: 64)),
        const Text('UK Tea House', style: TextStyle(fontFamily: 'serif', fontSize: 28, fontWeight: FontWeight.w700)),
        const SizedBox(height: 80),
        if (_step == 0) ...[
          const Text('Sign in with your email — no password needed', textAlign: TextAlign.center, style: TextStyle(color: Colors.black54)),
          const SizedBox(height: 24),
          TextField(controller: _emailCtrl, decoration: const InputDecoration(labelText: 'Email', border: OutlineInputBorder(), prefixIcon: Icon(Icons.email))),
          const SizedBox(height: 20),
          SizedBox(width: double.infinity, child: ElevatedButton(onPressed: _loading ? null : _requestLink, style: ElevatedButton.styleFrom(backgroundColor: Colors.green.shade800, padding: const EdgeInsets.all(16)),
            child: _loading ? const CircularProgressIndicator(color: Colors.white) : const Text('📧 Send Magic Link', style: TextStyle(fontSize: 16)))),
        ] else if (_step == 1) ...[
          const Icon(Icons.mark_email_read, size: 64, color: Colors.green),
          const SizedBox(height: 16),
          Text('We sent a link to $_savedEmail\nOr paste the code we sent', textAlign: TextAlign.center, style: const TextStyle(fontSize: 15)),
          const SizedBox(height: 24),
          TextField(controller: _tokenCtrl, maxLength: 6, textAlign: TextAlign.center, decoration: const InputDecoration(labelText: '6-digit code', border: OutlineInputBorder(), counterText: '')),
          const SizedBox(height: 20),
          SizedBox(width: double.infinity, child: ElevatedButton(onPressed: _loading ? null : _verify, style: ElevatedButton.styleFrom(backgroundColor: Colors.green.shade800, padding: const EdgeInsets.all(16)),
            child: _loading ? const CircularProgressIndicator(color: Colors.white) : const Text('✅ Verify & Sign In', style: TextStyle(fontSize: 16)))),
          TextButton(onPressed: () => setState(() => _step = 0), child: const Text('Use different email')),
        ],
      ])))),
    );
  }
}
