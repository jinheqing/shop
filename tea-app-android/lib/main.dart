import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:tea_app_android/services/api.dart';
import 'pages/login.dart';
import 'pages/home/home_page.dart';
import 'pages/farmer/farmer_dashboard.dart';
import 'pages/advisor/advisor_dashboard.dart';
import 'pages/live/live_page.dart';

void main() { runApp(const TeaApp()); }

final router = GoRouter(initialLocation: '/', routes: [
  GoRoute(path: '/', builder: (_, __) => const LoginPage()),
  GoRoute(path: '/home', builder: (_, __) => const HomePage()),
  GoRoute(path: '/farmer', builder: (_, __) => const FarmerDashboard()),
  GoRoute(path: '/advisor', builder: (_, __) => const AdvisorDashboard()),
  GoRoute(path: '/live/:room', builder: (_, s) => LivePage(room: s.pathParameters['room']!)),
]);

class TeaApp extends StatelessWidget {
  const TeaApp({super.key});
  @override
  Widget build(BuildContext ctx) => MaterialApp.router(
    title: 'UK Tea House',
    theme: ThemeData(colorSchemeSeed: const Color(0xFF5a3a1c), useMaterial3: true),
    routerConfig: router,
  );
}
