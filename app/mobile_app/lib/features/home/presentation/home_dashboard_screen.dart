import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:lucide_icons/lucide_icons.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';

class HomeDashboardScreen extends ConsumerStatefulWidget {
  const HomeDashboardScreen({super.key});

  @override
  ConsumerState<HomeDashboardScreen> createState() => _HomeDashboardScreenState();
}

class _HomeDashboardScreenState extends ConsumerState<HomeDashboardScreen> {
  String _userName = 'User';

  @override
  void initState() {
    super.initState();
    _loadUserName();
  }

  Future<void> _loadUserName() async {
    final prefs = await SharedPreferences.getInstance();
    setState(() {
      _userName = prefs.getString('user_name') ?? 'User';
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        decoration: const BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
            colors: [Color(0xFF0B1120), Color(0xFF131A2A)],
          ),
        ),
        child: SafeArea(
          child: SingleChildScrollView(
            padding: const EdgeInsets.all(20),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Header
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text('Welcome back,', style: TextStyle(color: Colors.blue[200], fontSize: 14)),
                        const SizedBox(height: 4),
                        Text(_userName, style: const TextStyle(color: Colors.white, fontSize: 24, fontWeight: FontWeight.bold)),
                      ],
                    ),
                    Row(
                      children: [
                        _buildHeaderIcon(LucideIcons.bell),
                        const SizedBox(width: 12),
                        _buildHeaderIcon(LucideIcons.settings),
                      ],
                    ),
                  ],
                ),
                const SizedBox(height: 24),

                // Balance Card
                ClipRRect(
                  borderRadius: BorderRadius.circular(20),
                  child: BackdropFilter(
                    filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10),
                    child: Container(
                      width: double.infinity,
                      padding: const EdgeInsets.all(24),
                      decoration: BoxDecoration(
                        gradient: const LinearGradient(
                          begin: Alignment.topLeft,
                          end: Alignment.bottomRight,
                          colors: [Color(0xFF1E3A8A), Color(0xFF312E81)],
                        ),
                        borderRadius: BorderRadius.circular(20),
                        border: Border.all(color: Colors.white.withOpacity(0.15)),
                      ),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Text('Total Balance', style: TextStyle(color: Colors.blue[200], fontSize: 14)),
                              Container(
                                padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                                decoration: BoxDecoration(
                                  color: Colors.green.withOpacity(0.2),
                                  borderRadius: BorderRadius.circular(20),
                                ),
                                child: const Row(
                                  children: [
                                    Icon(LucideIcons.trendingUp, color: Colors.greenAccent, size: 14),
                                    SizedBox(width: 4),
                                    Text('+0.0%', style: TextStyle(color: Colors.greenAccent, fontSize: 12, fontWeight: FontWeight.bold)),
                                  ],
                                ),
                              ),
                            ],
                          ),
                          const SizedBox(height: 8),
                          const Text(
                            'UGX 0',
                            style: TextStyle(color: Colors.white, fontSize: 36, fontWeight: FontWeight.w900, letterSpacing: 1),
                          ),
                          const SizedBox(height: 20),
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceAround,
                            children: [
                              _buildQuickAction(LucideIcons.plus, 'Deposit'),
                              _buildQuickAction(LucideIcons.arrowUpRight, 'Send'),
                              _buildQuickAction(LucideIcons.arrowDownLeft, 'Receive'),
                              _buildQuickAction(LucideIcons.history, 'History'),
                            ],
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
                const SizedBox(height: 28),

                // Services Section
                const Text('Services', style: TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold)),
                const SizedBox(height: 16),
                GridView.count(
                  crossAxisCount: 3,
                  shrinkWrap: true,
                  physics: const NeverScrollableScrollPhysics(),
                  mainAxisSpacing: 16,
                  crossAxisSpacing: 16,
                  childAspectRatio: 0.95,
                  children: [
                    _buildServiceTile(
                      icon: LucideIcons.wallet,
                      label: 'Wallet',
                      color: Colors.blueAccent,
                      onTap: () {},
                    ),
                    _buildServiceTile(
                      icon: LucideIcons.refreshCw,
                      label: 'Kibiina',
                      color: const Color(0xFF10B981),
                      onTap: () => context.push('/kibiina-register'),
                    ),
                    _buildServiceTile(
                      icon: LucideIcons.barChart3,
                      label: 'Shares',
                      color: Colors.purpleAccent,
                      onTap: () {},
                    ),
                    _buildServiceTile(
                      icon: LucideIcons.users,
                      label: 'Groups',
                      color: Colors.orangeAccent,
                      onTap: () {},
                    ),
                    _buildServiceTile(
                      icon: LucideIcons.landmark,
                      label: 'Loans',
                      color: Colors.tealAccent,
                      onTap: () {},
                    ),
                    _buildServiceTile(
                      icon: LucideIcons.shield,
                      label: 'Insurance',
                      color: Colors.redAccent,
                      onTap: () {},
                    ),
                    _buildServiceTile(
                      icon: LucideIcons.receipt,
                      label: 'Fees',
                      color: Colors.amberAccent,
                      onTap: () {},
                    ),
                    _buildServiceTile(
                      icon: LucideIcons.zap,
                      label: 'Utilities',
                      color: Colors.lightBlueAccent,
                      onTap: () {},
                    ),
                    _buildServiceTile(
                      icon: LucideIcons.moreHorizontal,
                      label: 'More',
                      color: Colors.grey,
                      onTap: () {},
                    ),
                  ],
                ),
                const SizedBox(height: 28),

                // Recent Activity
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text('Recent Activity', style: TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold)),
                    Text('See all', style: TextStyle(color: Colors.blue[300], fontSize: 14, fontWeight: FontWeight.w500)),
                  ],
                ),
                const SizedBox(height: 16),
                _buildEmptyState(),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildHeaderIcon(IconData icon) {
    return Container(
      padding: const EdgeInsets.all(10),
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.08),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: Colors.white.withOpacity(0.1)),
      ),
      child: Icon(icon, color: Colors.white70, size: 20),
    );
  }

  Widget _buildQuickAction(IconData icon, String label) {
    return Column(
      children: [
        Container(
          padding: const EdgeInsets.all(12),
          decoration: BoxDecoration(
            color: Colors.white.withOpacity(0.15),
            borderRadius: BorderRadius.circular(12),
          ),
          child: Icon(icon, color: Colors.white, size: 20),
        ),
        const SizedBox(height: 8),
        Text(label, style: const TextStyle(color: Colors.white70, fontSize: 12)),
      ],
    );
  }

  Widget _buildServiceTile({
    required IconData icon,
    required String label,
    required Color color,
    required VoidCallback onTap,
  }) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        decoration: BoxDecoration(
          color: Colors.white.withOpacity(0.06),
          borderRadius: BorderRadius.circular(16),
          border: Border.all(color: Colors.white.withOpacity(0.1)),
        ),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: color.withOpacity(0.15),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Icon(icon, color: color, size: 24),
            ),
            const SizedBox(height: 10),
            Text(label, style: const TextStyle(color: Colors.white, fontSize: 13, fontWeight: FontWeight.w500)),
          ],
        ),
      ),
    );
  }

  Widget _buildEmptyState() {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(32),
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.04),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: Colors.white.withOpacity(0.08)),
      ),
      child: Column(
        children: [
          Icon(LucideIcons.inbox, color: Colors.white.withOpacity(0.2), size: 48),
          const SizedBox(height: 12),
          Text('No recent activity', style: TextStyle(color: Colors.white.withOpacity(0.4), fontSize: 16)),
          const SizedBox(height: 4),
          Text('Your transactions will appear here', style: TextStyle(color: Colors.white.withOpacity(0.25), fontSize: 13)),
        ],
      ),
    );
  }
}
