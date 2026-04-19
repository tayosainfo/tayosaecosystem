import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:lucide_icons/lucide_icons.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../../core/network/api_client.dart';

class KibiinaDashboardScreen extends ConsumerStatefulWidget {
  const KibiinaDashboardScreen({super.key});

  @override
  ConsumerState<KibiinaDashboardScreen> createState() => _KibiinaDashboardScreenState();
}

class _KibiinaDashboardScreenState extends ConsumerState<KibiinaDashboardScreen> {
  Map<String, dynamic>? memberHistory;
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _fetchDashboard();
  }

  Future<void> _fetchDashboard() async {
    setState(() => _isLoading = true);
    try {
      final prefs = await SharedPreferences.getInstance();
      var userId = prefs.getInt('user_id') ?? 1;

      final api = ref.read(apiClientProvider);
      
      // Port 8086 is the Kibiina Service. Assuming member's backend ID is identical to user_id for demo.
      final kibiinaBaseUrl = Uri.parse(api.options.baseUrl).replace(port: 8086).toString();
      final response = await api.get('$kibiinaBaseUrl/api/v1/kibiina/members/$userId/history');

      if (response.statusCode == 200) {
        setState(() => memberHistory = response.data);
      }
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Could not fetch data: $e')));
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return Scaffold(
        appBar: AppBar(title: const Text('Kibiina Dashboard'), backgroundColor: Colors.transparent, elevation: 0),
        body: const Center(child: CircularProgressIndicator()),
      );
    }

    if (memberHistory == null || memberHistory!['member'] == null) {
      return Scaffold(
        appBar: AppBar(title: const Text('Kibiina Dashboard'), backgroundColor: Colors.transparent, elevation: 0),
        body: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(LucideIcons.users, size: 64, color: Theme.of(context).primaryColor.withOpacity(0.5)),
              const SizedBox(height: 16),
              const Text('You have not joined a Kibiina yet.'),
              const SizedBox(height: 16),
              ElevatedButton(
                onPressed: () => context.go('/geo-onboarding'),
                child: const Text('Join a Kibiina'),
              )
            ],
          ),
        ),
      );
    }

    final member = memberHistory!['member'];
    final seats = memberHistory!['seats'] as List;
    final seat = seats.isNotEmpty ? seats.first : null;
    final kibiina = seat != null ? seat['kibiina'] : null;
    
    // Calculate values
    final contribution = member['contribution_amount'];
    final totalRounds = kibiina != null ? kibiina['total_rounds'] : 0;
    final expectedPayout = kibiina != null && seat != null 
        ? contribution * totalRounds 
        : 0;

    return Scaffold(
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        title: const Text('My Kibiina'),
        centerTitle: true,
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(24.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            // Status Header
            Container(
              padding: const EdgeInsets.all(20),
              decoration: BoxDecoration(
                color: Theme.of(context).cardColor,
                borderRadius: BorderRadius.circular(16),
                border: Border.all(color: const Color(0xFF2A3347)),
              ),
              child: Row(
                children: [
                  Container(
                    width: 48,
                    height: 48,
                    decoration: BoxDecoration(
                      color: Theme.of(context).primaryColor.withOpacity(0.2),
                      shape: BoxShape.circle,
                    ),
                    child: Icon(LucideIcons.refreshCw, color: Theme.of(context).primaryColor),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          kibiina?['name'] ?? 'Waiting for Group...',
                          style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
                        ),
                        Text(
                          'Status: ${kibiina?['status']?.toString().toUpperCase() ?? "PENDING"}',
                          style: TextStyle(
                            color: kibiina?['status'] == 'active' ? Theme.of(context).primaryColor : Colors.amber, 
                            fontWeight: FontWeight.bold,
                            fontSize: 12
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
            
            const SizedBox(height: 24),
            
            // Payout Card heavily matching Admin UI style
            Container(
              padding: const EdgeInsets.all(24),
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  colors: [Theme.of(context).primaryColor, const Color(0xFF047857)],
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                ),
                borderRadius: BorderRadius.circular(20),
                boxShadow: [
                  BoxShadow(color: Theme.of(context).primaryColor.withOpacity(0.3), blurRadius: 20, offset: const Offset(0, 10))
                ],
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text('EXPECTED PAYOUT', style: TextStyle(color: Colors.white70, fontWeight: FontWeight.bold, letterSpacing: 1.5)),
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                        decoration: BoxDecoration(color: Colors.white24, borderRadius: BorderRadius.circular(20)),
                        child: Text(
                          'Seat #${seat?['position'] ?? '?'} / $totalRounds', 
                          style: const TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 12)
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 8),
                  Text(
                    'UGX ${expectedPayout.toString().replaceAllMapped(RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'), (m) => '${m[1]},')}',
                    style: const TextStyle(color: Colors.white, fontSize: 32, fontWeight: FontWeight.bold),
                  ),
                  const SizedBox(height: 16),
                  Row(
                    children: [
                      const Icon(LucideIcons.shieldCheck, color: Colors.white, size: 16),
                      const SizedBox(width: 8),
                      Text(
                        'Collateral Built: UGX ${member['collateral_balance'].toString().replaceAllMapped(RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'), (m) => '${m[1]},')}', 
                        style: const TextStyle(color: Colors.white)
                      ),
                    ],
                  )
                ],
              ),
            ),

            const SizedBox(height: 32),
            
            // Action Buttons
            Row(
              children: [
                Expanded(
                  child: ElevatedButton.icon(
                    onPressed: kibiina?['status'] == 'active' ? () => context.go('/kibiina-deposit') : null,
                    icon: const Icon(LucideIcons.banknote),
                    label: const Text('Make Deposit'),
                    style: ElevatedButton.styleFrom(backgroundColor: Theme.of(context).primaryColor),
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: ElevatedButton.icon(
                    onPressed: () {}, // Navigate to full rounds timeline
                    icon: const Icon(LucideIcons.history),
                    label: const Text('View Timeline'),
                    style: ElevatedButton.styleFrom(
                      backgroundColor: Theme.of(context).cardColor,
                      foregroundColor: Colors.white,
                      side: const BorderSide(color: Color(0xFF2A3347)),
                    ),
                  ),
                ),
              ],
            ),

            const SizedBox(height: 32),
            const Text('Recent Deposists', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
            const SizedBox(height: 16),

            // History
            memberHistory!['deposits'] == null || (memberHistory!['deposits'] as List).isEmpty
            ? const Center(child: Padding(
                padding: EdgeInsets.all(24.0),
                child: Text('No deposits recorded yet.', style: TextStyle(color: Colors.white54)),
              ))
            : ListView.builder(
                shrinkWrap: true,
                physics: const NeverScrollableScrollPhysics(),
                itemCount: (memberHistory!['deposits'] as List).length,
                itemBuilder: (context, index) {
                  final deposit = (memberHistory!['deposits'] as List)[index];
                  final date = DateTime.parse(deposit['created_at']);
                  return Container(
                    margin: const EdgeInsets.only(bottom: 12),
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(
                      color: Theme.of(context).cardColor,
                      borderRadius: BorderRadius.circular(12),
                      border: Border.all(color: const Color(0xFF2A3347)),
                    ),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Row(
                          children: [
                            Container(
                              padding: const EdgeInsets.all(8),
                              decoration: BoxDecoration(
                                color: deposit['is_default_payment'] ? Colors.amber.withOpacity(0.1) : Theme.of(context).primaryColor.withOpacity(0.1),
                                borderRadius: BorderRadius.circular(8),
                              ),
                              child: Icon(
                                deposit['is_default_payment'] ? LucideIcons.alertTriangle : LucideIcons.checkCircle, 
                                color: deposit['is_default_payment'] ? Colors.amber : Theme.of(context).primaryColor,
                                size: 16,
                              ),
                            ),
                            const SizedBox(width: 12),
                            Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(deposit['is_default_payment'] ? 'Default Recovery' : 'Regular Deposit', style: const TextStyle(fontWeight: FontWeight.bold)),
                                Text('${date.day}/${date.month}/${date.year}', style: const TextStyle(color: Colors.white54, fontSize: 12)),
                              ],
                            ),
                          ],
                        ),
                        Text(
                          '+UGX ${deposit['amount']}',
                          style: TextStyle(color: deposit['is_default_payment'] ? Colors.amber : Theme.of(context).primaryColor, fontWeight: FontWeight.bold),
                        )
                      ],
                    ),
                  );
                },
              ),
          ],
        ),
      ),
    );
  }
}
