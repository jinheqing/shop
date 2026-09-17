import 'package:flutter/material.dart';

/// Tea App Android — Step 1 scaffolding
/// 
/// One APK, two roles:
///   role=tea_farmer  → Farmer dashboard (Go Live, orders)
///   role=advisor     → Advisor dashboard (IM, quotes, Go Live, nodes)
void main() {
  runApp(const TeaApp());
}

class TeaApp extends StatelessWidget {
  const TeaApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Tea System',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: const Color(0xFFB8A47C)),
        useMaterial3: true,
      ),
      home: const Scaffold(
        body: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text(
                'Tea System',
                style: TextStyle(fontSize: 32, fontWeight: FontWeight.bold),
              ),
              SizedBox(height: 12),
              Text(
                'Step 1 scaffolding ✅',
                style: TextStyle(fontSize: 16, color: Colors.grey),
              ),
              SizedBox(height: 4),
              Text(
                'Flutter 3.38 + livekit_client',
                style: TextStyle(fontSize: 14, color: Colors.grey),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
