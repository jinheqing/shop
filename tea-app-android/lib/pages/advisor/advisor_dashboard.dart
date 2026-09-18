import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:tea_app_android/services/api.dart';

class AdvisorDashboard extends StatefulWidget {
  const AdvisorDashboard({super.key});
  @override
  State<AdvisorDashboard> createState() => _AdvisorDashboardState();
}

class _AdvisorDashboardState extends State<AdvisorDashboard> {
  int tab = 0;
  Future<List>? orders;
  Future<List>? rooms;
  Future<List>? nodes;

  @override void initState() { super.initState(); orders = ApiService().getOrders(); rooms = ApiService().getLiveRooms(); nodes = ApiService().getNodes(); }

  @override
  Widget build(ctx) => Scaffold(
    appBar: AppBar(title: const Text('🎓 Advisor Console'), actions: [IconButton(icon: const Icon(Icons.logout), onPressed: () => context.go('/'))]),
    body: IndexedStack(index: tab, children: [_buildOrders(), _buildLive(), _buildNodes()]),
    bottomNavigationBar: NavigationBar(destinations: const [
      NavigationDestination(icon: Icon(Icons.shopping_cart), label: 'Orders'),
      NavigationDestination(icon: Icon(Icons.videocam), label: 'Live'),
      NavigationDestination(icon: Icon(Icons.cloud_queue), label: 'Nodes'),
    ], selectedIndex: tab, onDestinationSelected: (i) => setState(() => tab = i)),
  );

  Widget _buildOrders() => FutureBuilder<List>(
    future: orders,
    builder: (ctx, snap) {
      if (snap.connectionState != ConnectionState.done) return const Center(child: CircularProgressIndicator());
      if (!snap.hasData || snap.data!.isEmpty) return const Center(child: Text('No orders'));
      return ListView.builder(itemCount: snap.data!.length, itemBuilder: (_, i) {
        final o = snap.data![i];
        return Card(margin: const EdgeInsets.all(8), child: ExpansionTile(
          title: Text(o['order_no'] ?? ''),
          subtitle: Text('£${o['total_amount']}'),
          trailing: Chip(label: Text(o['state'] ?? '')),
          children: [
            ListTile(leading: const Icon(Icons.confirmation_num), title: Text('HS: ${o['hs_code']}')),
            ListTile(leading: const Icon(Icons.factory), title: Text('Origin: ${o['country_of_origin']}')),
            OverflowBar(alignment: MainAxisAlignment.end, children: [
              TextButton(child: const Text('Mark Paid'), onPressed: () async { await ApiService().updateOrderState(o['id'], 'paid'); setState(() {}); }),
              TextButton(child: const Text('Ship'), onPressed: () async { await ApiService().updateOrderState(o['id'], 'shipped'); setState(() {}); }),
            ]),
          ],
        ));
      });
    },
  );

  Widget _buildLive() => FutureBuilder<List>(
    future: rooms,
    builder: (ctx, snap) {
      if (snap.connectionState != ConnectionState.done) return const Center(child: CircularProgressIndicator());
      if (!snap.hasData || snap.data!.isEmpty) return const Center(child: Text('No rooms'));
      return ListView.builder(itemCount: snap.data!.length, itemBuilder: (_, i) {
        final r = snap.data![i];
        return Card(margin: const EdgeInsets.all(8), child: ListTile(
          leading: Icon(Icons.circle, color: r['status']=='live' ? Colors.red : Colors.grey),
          title: Text(r['room_name'] ?? ''),
          subtitle: Text('${r['room_type']} · ${r['status']}'),
          trailing: FilledButton.tonal(child: const Text('Join'), onPressed: () => context.push('/live/${r['room_id']}')),
        ));
      });
    },
  );

  Widget _buildNodes() => FutureBuilder<List>(
    future: nodes,
    builder: (ctx, snap) {
      if (snap.connectionState != ConnectionState.done) return const Center(child: CircularProgressIndicator());
      if (!snap.hasData || snap.data!.isEmpty) return const Center(child: Text('No nodes'));
      return ListView.builder(itemCount: snap.data!.length, itemBuilder: (_, i) {
        final n = snap.data![i];
        return Card(margin: const EdgeInsets.all(8), child: ListTile(
          leading: Icon(Icons.computer, color: n['status']=='online' ? Colors.green : Colors.grey),
          title: Text(n['node_name'] ?? ''),
          subtitle: Text('${n['node_type']} · ${n['public_ip']} → ${n['wireguard_ip']}'),
          trailing: Text(n['status'] ?? ''),
        ));
      });
    },
  );
}
