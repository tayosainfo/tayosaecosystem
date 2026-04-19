import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:lucide_icons/lucide_icons.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../../core/network/api_client.dart';

class KibiinaRegistrationScreen extends ConsumerStatefulWidget {
  const KibiinaRegistrationScreen({super.key});

  @override
  ConsumerState<KibiinaRegistrationScreen> createState() => _KibiinaRegistrationScreenState();
}

class _KibiinaRegistrationScreenState extends ConsumerState<KibiinaRegistrationScreen> {
  String _selectedCycle = 'daily';
  int _contributionAmount = 50000;
  bool _isLoading = false;
  String? _error;

  final List<Map<String, String>> _cycles = [
    {'key': 'daily', 'label': 'Daily Savings', 'desc': 'Round is 7 days'},
    {'key': 'weekly', 'label': 'Weekly Savings', 'desc': 'Round is 4 weeks'},
    {'key': 'monthly', 'label': 'Monthly Savings', 'desc': 'Round is 12 months'},
    {'key': 'annual', 'label': 'Annual Savings', 'desc': 'Round is 52 weeks'},
  ];

  final List<int> _amounts = [10000, 20000, 50000, 100000, 200000, 500000];

  Future<void> _handleRegister() async {
    setState(() { _isLoading = true; _error = null; });

    try {
      final prefs = await SharedPreferences.getInstance();
      final userId = prefs.getInt('user_id') ?? 0;
      final userName = prefs.getString('user_name') ?? 'Unknown';
      final userPhone = prefs.getString('user_phone') ?? '';
      final userDistrict = prefs.getString('user_district') ?? '';
      final userVillage = prefs.getString('user_village') ?? '';
      final api = ref.read(apiClientProvider);

      // Kibiina service runs on port 8086 (separate from auth on 8080)
      final kibiinaBaseUrl = Uri.parse(api.options.baseUrl).replace(port: 8086).toString();

      final response = await api.post(
        '$kibiinaBaseUrl/api/v1/kibiina/members/register',
        data: {
          'user_id': userId,
          'full_name': userName,
          'phone': userPhone,
          'district': userDistrict,
          'village': userVillage,
          'savings_type': _selectedCycle,
          'contribution_amount': _contributionAmount,
        },
      );

      if (response.statusCode == 200 || response.statusCode == 201) {
        if (mounted) context.go('/kibiina-dashboard');
      } else {
        setState(() => _error = response.data['error'] ?? 'Registration failed');
      }
    } catch (e) {
      setState(() => _error = 'Network error: $e');
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        decoration: const BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [Color(0xFF1E3A8A), Color(0xFF1E40AF), Color(0xFF312E81)],
          ),
        ),
        child: Stack(
          children: [
            Positioned(top: -100, right: -100, child: _buildOrb(Colors.blue.withOpacity(0.15))),
            Positioned(bottom: -100, left: -100, child: _buildOrb(Colors.indigo.withOpacity(0.15))),

            SafeArea(
              child: Column(
                children: [
                  // App bar
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                    child: Row(
                      children: [
                        IconButton(
                          onPressed: () => context.pop(),
                          icon: const Icon(LucideIcons.arrowLeft, color: Colors.white),
                        ),
                        const Expanded(
                          child: Text('Join a Kibiina', textAlign: TextAlign.center,
                            style: TextStyle(color: Colors.white, fontSize: 20, fontWeight: FontWeight.bold)),
                        ),
                        const SizedBox(width: 48),
                      ],
                    ),
                  ),

                  Expanded(
                    child: SingleChildScrollView(
                      padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.stretch,
                        children: [
                          // Info card
                          ClipRRect(
                            borderRadius: BorderRadius.circular(16),
                            child: BackdropFilter(
                              filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10),
                              child: Container(
                                padding: const EdgeInsets.all(20),
                                decoration: BoxDecoration(
                                  color: Colors.white.withOpacity(0.1),
                                  borderRadius: BorderRadius.circular(16),
                                  border: Border.all(color: Colors.white.withOpacity(0.2)),
                                ),
                                child: Column(
                                  children: [
                                    Container(
                                      padding: const EdgeInsets.all(12),
                                      decoration: BoxDecoration(
                                        color: const Color(0xFF10B981).withOpacity(0.15),
                                        borderRadius: BorderRadius.circular(12),
                                      ),
                                      child: const Icon(LucideIcons.refreshCw, color: Color(0xFF10B981), size: 32),
                                    ),
                                    const SizedBox(height: 16),
                                    const Text('Setup your Merry-Go-Round',
                                      style: TextStyle(color: Colors.white, fontSize: 20, fontWeight: FontWeight.bold)),
                                    const SizedBox(height: 8),
                                    Text('Choose your cycle and how much you want to contribute. You will be matched with others locally.',
                                      textAlign: TextAlign.center,
                                      style: TextStyle(color: Colors.blue[200], fontSize: 14)),
                                  ],
                                ),
                              ),
                            ),
                          ),
                          const SizedBox(height: 28),

                          // Savings Cycle
                          const Text('Savings Cycle', style: TextStyle(color: Colors.white, fontSize: 16, fontWeight: FontWeight.bold)),
                          const SizedBox(height: 12),
                          ...(_cycles.map((cycle) {
                            final isSelected = _selectedCycle == cycle['key'];
                            return Padding(
                              padding: const EdgeInsets.only(bottom: 10),
                              child: GestureDetector(
                                onTap: () => setState(() => _selectedCycle = cycle['key']!),
                                child: Container(
                                  padding: const EdgeInsets.all(16),
                                  decoration: BoxDecoration(
                                    color: isSelected
                                        ? const Color(0xFF10B981).withOpacity(0.15)
                                        : Colors.white.withOpacity(0.06),
                                    borderRadius: BorderRadius.circular(12),
                                    border: Border.all(
                                      color: isSelected
                                          ? const Color(0xFF10B981).withOpacity(0.5)
                                          : Colors.white.withOpacity(0.1),
                                      width: isSelected ? 2 : 1,
                                    ),
                                  ),
                                  child: Row(
                                    children: [
                                      Container(
                                        width: 24, height: 24,
                                        decoration: BoxDecoration(
                                          shape: BoxShape.circle,
                                          color: isSelected ? const Color(0xFF10B981) : Colors.transparent,
                                          border: Border.all(
                                            color: isSelected ? const Color(0xFF10B981) : Colors.white30,
                                            width: 2,
                                          ),
                                        ),
                                        child: isSelected
                                            ? const Icon(Icons.check, color: Colors.white, size: 14)
                                            : null,
                                      ),
                                      const SizedBox(width: 16),
                                      Column(
                                        crossAxisAlignment: CrossAxisAlignment.start,
                                        children: [
                                          Text(cycle['label']!, style: TextStyle(color: isSelected ? Colors.white : Colors.white70, fontSize: 15, fontWeight: FontWeight.w600)),
                                          Text(cycle['desc']!, style: TextStyle(color: Colors.blue[200], fontSize: 13)),
                                        ],
                                      ),
                                    ],
                                  ),
                                ),
                              ),
                            );
                          })),

                          const SizedBox(height: 24),

                          // Contribution Amount
                          const Text('Contribution Amount', style: TextStyle(color: Colors.white, fontSize: 16, fontWeight: FontWeight.bold)),
                          const SizedBox(height: 12),
                          Container(
                            decoration: BoxDecoration(
                              color: Colors.white.withOpacity(0.1),
                              borderRadius: BorderRadius.circular(12),
                              border: Border.all(color: Colors.white.withOpacity(0.3)),
                            ),
                            padding: const EdgeInsets.symmetric(horizontal: 16),
                            child: DropdownButtonHideUnderline(
                              child: DropdownButton<int>(
                                isExpanded: true,
                                value: _contributionAmount,
                                dropdownColor: const Color(0xFF1E3A8A),
                                style: const TextStyle(color: Colors.white, fontSize: 16),
                                icon: Icon(LucideIcons.chevronDown, color: Colors.blue[300]),
                                items: _amounts.map((amount) {
                                  return DropdownMenuItem<int>(
                                    value: amount,
                                    child: Text('UGX ${_formatAmount(amount)}'),
                                  );
                                }).toList(),
                                onChanged: (v) => setState(() => _contributionAmount = v!),
                              ),
                            ),
                          ),

                          const SizedBox(height: 20),

                          // Fee notice
                          Container(
                            padding: const EdgeInsets.all(14),
                            decoration: BoxDecoration(
                              color: Colors.amber.withOpacity(0.1),
                              borderRadius: BorderRadius.circular(12),
                              border: Border.all(color: Colors.amber.withOpacity(0.3)),
                            ),
                            child: Row(
                              children: [
                                const Icon(LucideIcons.alertTriangle, color: Colors.amber, size: 18),
                                const SizedBox(width: 12),
                                Expanded(
                                  child: Text(
                                    'TAYOSA charges a one-time non-refundable Kibiina registration fee of UGX 5,000.',
                                    style: TextStyle(color: Colors.amber[200], fontSize: 13),
                                  ),
                                ),
                              ],
                            ),
                          ),

                          if (_error != null) ...[
                            const SizedBox(height: 16),
                            Container(
                              padding: const EdgeInsets.all(12),
                              decoration: BoxDecoration(
                                color: Colors.red[500]!.withOpacity(0.2),
                                border: Border.all(color: Colors.red[400]!.withOpacity(0.5)),
                                borderRadius: BorderRadius.circular(12),
                              ),
                              child: Text(_error!, style: TextStyle(color: Colors.red[100])),
                            ),
                          ],

                          const SizedBox(height: 32),

                          ElevatedButton(
                            onPressed: _isLoading ? null : _handleRegister,
                            style: ElevatedButton.styleFrom(
                              backgroundColor: const Color(0xFF10B981),
                              foregroundColor: Colors.white,
                              padding: const EdgeInsets.symmetric(vertical: 16),
                              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                              elevation: 5,
                            ),
                            child: _isLoading
                                ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2))
                                : Row(
                                    mainAxisAlignment: MainAxisAlignment.center,
                                    children: const [
                                      Text('Join Kibiina', style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
                                      SizedBox(width: 8),
                                      Icon(LucideIcons.arrowRight, size: 20),
                                    ],
                                  ),
                          ),
                        ],
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  String _formatAmount(int amount) {
    final str = amount.toString();
    final buffer = StringBuffer();
    for (int i = 0; i < str.length; i++) {
      if (i > 0 && (str.length - i) % 3 == 0) buffer.write(',');
      buffer.write(str[i]);
    }
    return buffer.toString();
  }

  Widget _buildOrb(Color color) {
    return Container(width: 300, height: 300, decoration: BoxDecoration(shape: BoxShape.circle, color: color));
  }
}
