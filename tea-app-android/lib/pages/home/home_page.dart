import 'package:flutter/material.dart';

class HomePage extends StatelessWidget {
  const HomePage({super.key});
  @override
  Widget build(ctx) => Scaffold(
    appBar: AppBar(title: const Text('UK Tea House')),
    body: const Center(child: Text('Home')),
  );
}
