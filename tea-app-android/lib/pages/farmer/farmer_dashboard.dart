import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:tea_app_android/services/api.dart';

class FarmerDashboard extends StatefulWidget {
  const FarmerDashboard({super.key});
  @override
  State<FarmerDashboard> createState() => _FarmerDashboardState();
}

class _FarmerDashboardState extends State<FarmerDashboard> {
  int tab = 0;
  Future<List>? orders;
  Future<List>? presets;

  @override void initState() { super.initState(); orders = ApiService().getOrders(); presets = ApiService().getSlowPresets(); }

  @override
  Widget build(ctx) => Scaffold(
    appBar: AppBar(title: const Text('👨‍🌾 Farmer Dashboard'), actions: [IconButton(icon: const Icon(Icons.logout), onPressed: () => context.go('/'))]),
    body: IndexedStack(index: tab, children: [_buildOrders(), _buildLive()]),
    bottomNavigationBar: NavigationBar(destinations: const [
      NavigationDestination(icon: Icon(Icons.inventory_2), label: 'Orders'),
      NavigationDestination(icon: Icon(Icons.videocam), label: 'Live Streams'),
    ], selectedIndex: tab, onDestinationSelected: (i) => setState(() => tab = i)),
  );

  Widget _buildOrders() => FutureBuilder<List>(
    future: orders,
    builder: (ctx, snap) {
      if (snap.connectionState != ConnectionState.done) return const Center(child: CircularProgressIndicator());
      if (!snap.hasData || snap.data!.isEmpty) return const Center(child: Text('No orders yet'));
      return ListView.builder(itemCount: snap.data!.length, itemBuilder: (_, i) {
        final o = snap.data![i];
        return Card(margin: const EdgeInsets.all(8), child: ListTile(
          title: Text(o['order_no'] ?? ''),
          subtitle: Text('£${o['total_amount']} · ${o['state']}'),
          trailing: o['state']=='paid' ? const Icon(Icons.check_circle, color: Colors.green) : null,
        ));
      });
    },
  );

  Widget _buildLive() => FutureBuilder<List>(
    future: presets,
    builder: (ctx, snap) {
      if (snap.connectionState != ConnectionState.done) return const Center(child: CircularProgressIndicator());
      if (!snap.hasData || snap.data!.isEmpty) return const Center(child: Text('No camera presets'));
      return GridView.count(crossAxisCount: 2, padding: const EdgeInsets.all(8), crossAxisSpacing: 8, mainAxisSpacing: 8,
        children: snap.data!.map((p) => Card(clipBehavior: Clip.antiAlias, child: Column(children: [
          Container(height: 100, color: Colors.brown.shade900, alignment: Alignment.center, child: const Icon(Icons.videocam, color: Colors.white54, size: 48)),
          Padding(padding: const EdgeInsets.all(8), child: Text(p['name'] ?? '', style: const TextStyle(fontWeight: FontWeight.w600))),
        ]))).toList());
    },
  );
}
