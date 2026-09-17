import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:tea_app_android/services/api.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});
  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final email = TextEditingController();
  final password = TextEditingController();
  bool loading = false;
  String? error;

  Future<void> _login() async {
    setState(() { loading = true; error = null });
    try {
      final data = await ApiService().login(email.text.trim(), password.text.trim());
      final role = data['role'] ?? 'staff';
      if (!mounted) return;
      if (role == 'farmer') context.go('/farmer');
      else context.go('/advisor');
    } catch (e) {
      setState(() => error = 'Login failed: ${e.toString().split(':').last}');
    } finally {
      if (mounted) setState(() => loading = false);
    }
  }

  @override
  Widget build(BuildContext ctx) => Scaffold(
    body: Container(
      decoration: const BoxDecoration(
        gradient: LinearGradient(begin: Alignment.topLeft, end: Alignment.bottomRight,
          colors: [Color(0xFF311d0c), Color(0xFF5a3a1c), Color(0xFFa66e23)]),
      ),
      child: SafeArea(
        child: Center(
          child: SingleChildScrollView(
            padding: const EdgeInsets.all(24),
            child: Card(
              elevation: 8,
              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(20)),
              child: Padding(padding: const EdgeInsets.all(28), child: Column(mainAxisSize: MainAxisSize.min, children: [
                const Text('🍃', style: TextStyle(fontSize: 64)),
                const Text('UK Tea House', style: TextStyle(fontSize: 22, fontWeight: FontWeight.w600)),
                const SizedBox(height: 4),
                Text('Staff Portal · Pu\'er Private Domain', style: TextStyle(color: Colors.brown.shade700, fontSize: 12)),
                const SizedBox(height: 24),
                TextField(controller: email, decoration: const InputDecoration(labelText: 'Email', border: OutlineInputBorder())),
                const SizedBox(height: 12),
                TextField(controller: password, obscureText: true, decoration: const InputDecoration(labelText: 'Password', border: OutlineInputBorder())),
                if (error != null) ...[const SizedBox(height: 12), Text(error!, style: const TextStyle(color: Colors.red))],
                const SizedBox(height: 24),
                SizedBox(width: double.infinity, height: 48, child: FilledButton(onPressed: loading ? null : _login, child: loading ? const CircularProgressIndicator() : const Text('Sign In'))),
                const SizedBox(height: 12),
                Wrap(spacing: 8, runSpacing: 4, alignment: WrapAlignment.center, children: [
                  ChoiceChip(label: const Text('Advisor'), selected: email.text=='advisor@test.com', onSelected: (_) { email.text='advisor@test.com'; password.text='Advisor@12345'; setState(() {}); }),
                  ChoiceChip(label: const Text('Farmer'), selected: email.text=='farmer@test.com', onSelected: (_) { email.text='farmer@test.com'; password.text='Farmer@12345'; setState(() {}); }),
                ]),
              ])),
            ),
          ),
        ),
      ),
    ),
  );
}
