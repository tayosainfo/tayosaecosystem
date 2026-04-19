import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:lucide_icons/lucide_icons.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../../core/network/api_client.dart';

class KibiinaDepositScreen extends ConsumerStatefulWidget {
  const KibiinaDepositScreen({super.key});

  @override
  ConsumerState<KibiinaDepositScreen> createState() => _KibiinaDepositScreenState();
}

class _KibiinaDepositScreenState extends ConsumerState<KibiinaDepositScreen> {
  final _amountController = TextEditingController();
  String _selectedProvider = 'MTN Mobile Money';
  bool _isProcessing = false;
  bool _isLoadingMember = true;

  // Kibiina member context — fetched on init to ensure correct IDs
  int? _kibiinaMemberId;
  int? _kibiinaId;
  int? _seatId;

  @override
  void initState() {
    super.initState();
    _fetchMemberContext();
  }

  /// Fetches the logged-in user's Kibiina member record to obtain
  /// the correct member_id, kibiina_id, and seat_id needed for deposits.
  /// This prevents IDOR mismatches between auth user_id and kibiina member_id.
  Future<void> _fetchMemberContext() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final userId = prefs.getInt('user_id') ?? 1;

      final api = ref.read(apiClientProvider);
      final kibiinaBaseUrl = Uri.parse(api.options.baseUrl).replace(port: 8086).toString();
      final response = await api.get('$kibiinaBaseUrl/api/v1/kibiina/members/$userId/history');

      if (response.statusCode == 200 && response.data['member'] != null) {
        final member = response.data['member'];
        final seats = response.data['seats'] as List?;
        final seat = (seats != null && seats.isNotEmpty) ? seats.first : null;

        setState(() {
          _kibiinaMemberId = member['id'];
          _kibiinaId = member['kibiina_id'];
          _seatId = seat?['id'];
        });
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Could not load Kibiina data: $e'), backgroundColor: Colors.red),
        );
      }
    } finally {
      if (mounted) setState(() => _isLoadingMember = false);
    }
  }

  Future<void> _processPayment() async {
    if (_amountController.text.isEmpty) return;
    if (_kibiinaMemberId == null || _kibiinaId == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Kibiina membership not found. Please join a Kibiina first.'), backgroundColor: Colors.red),
      );
      return;
    }

    setState(() => _isProcessing = true);

    try {
      // Simulate MTN / Airtel Network Wait (STK Push delay)
      await Future.delayed(const Duration(seconds: 2));

      final api = ref.read(apiClientProvider);

      // Route to port 8086 (kibiina-service), NOT the default 8080 (auth-service)
      final kibiinaBaseUrl = Uri.parse(api.options.baseUrl).replace(port: 8086).toString();

      final response = await api.post(
        '$kibiinaBaseUrl/api/v1/kibiina/deposits',
        data: {
          'kibiina_id': _kibiinaId,            // Required: which Kibiina group
          'member_id': _kibiinaMemberId,        // Required: actual Kibiina member ID (matches JWT for IDOR check)
          'amount': int.parse(_amountController.text),  // Required: deposit amount
          'deposit_type': 'contribution',       // Required: regular contribution
          'channel': _selectedProvider == 'MTN Mobile Money' ? 'mtn_momo' : 'airtel_money',
          if (_seatId != null) 'seat_id': _seatId,
        },
      );

      if (response.statusCode == 200 || response.statusCode == 201) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Deposit Successful!'), backgroundColor: Colors.green));
          context.pop(); // Go back to dashboard
        }
      } else {
        if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: ${response.data['error']}')));
      }

    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Payment Failed: $e')));
    } finally {
      if (mounted) setState(() => _isProcessing = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Make a Deposit'),
        backgroundColor: Colors.transparent,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(LucideIcons.arrowLeft, color: Colors.white),
          onPressed: () => context.pop(),
        ),
      ),
      body: _isLoadingMember
          ? const Center(child: CircularProgressIndicator())
          : SingleChildScrollView(
              padding: const EdgeInsets.all(24.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  Container(
                    padding: const EdgeInsets.all(24),
                    decoration: BoxDecoration(
                      color: Theme.of(context).cardColor,
                      borderRadius: BorderRadius.circular(16),
                      border: Border.all(color: const Color(0xFF2A3347)),
                    ),
                    child: Column(
                      children: [
                        Icon(LucideIcons.smartphoneNfc, size: 48, color: Theme.of(context).primaryColor),
                        const SizedBox(height: 16),
                        const Text('Direct Mobile Money Transfer', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                        const SizedBox(height: 8),
                        const Text('We will send a prompt to your phone. Enter your PIN to approve the Kibiina deposit.', textAlign: TextAlign.center, style: TextStyle(color: Colors.white54)),
                      ],
                    ),
                  ),

                  // Show active Kibiina context
                  if (_kibiinaId != null) ...[
                    const SizedBox(height: 16),
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                      decoration: BoxDecoration(
                        color: Theme.of(context).primaryColor.withOpacity(0.1),
                        borderRadius: BorderRadius.circular(12),
                        border: Border.all(color: Theme.of(context).primaryColor.withOpacity(0.3)),
                      ),
                      child: Row(
                        children: [
                          Icon(LucideIcons.checkCircle, color: Theme.of(context).primaryColor, size: 16),
                          const SizedBox(width: 10),
                          Expanded(
                            child: Text(
                              'Depositing to Kibiina #$_kibiinaId • Member #$_kibiinaMemberId',
                              style: TextStyle(color: Theme.of(context).primaryColor, fontSize: 12, fontWeight: FontWeight.w600),
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                  
                  const SizedBox(height: 32),
                  Text('Select Network', style: Theme.of(context).textTheme.titleMedium),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      Expanded(
                        child: InkWell(
                          onTap: () => setState(() => _selectedProvider = 'MTN Mobile Money'),
                          borderRadius: BorderRadius.circular(12),
                          child: Container(
                            padding: const EdgeInsets.symmetric(vertical: 16),
                            decoration: BoxDecoration(
                              color: _selectedProvider == 'MTN Mobile Money' ? Colors.amber.withOpacity(0.1) : Theme.of(context).cardColor,
                              borderRadius: BorderRadius.circular(12),
                              border: Border.all(color: _selectedProvider == 'MTN Mobile Money' ? Colors.amber : const Color(0xFF2A3347), width: _selectedProvider == 'MTN Mobile Money' ? 2 : 1),
                            ),
                            child: Center(
                              child: Text('MTN MoMo', style: TextStyle(fontWeight: FontWeight.bold, color: _selectedProvider == 'MTN Mobile Money' ? Colors.amber : Colors.white)),
                            )
                          ),
                        )
                      ),
                      const SizedBox(width: 16),
                      Expanded(
                        child: InkWell(
                          onTap: () => setState(() => _selectedProvider = 'Airtel Money'),
                          borderRadius: BorderRadius.circular(12),
                          child: Container(
                            padding: const EdgeInsets.symmetric(vertical: 16),
                            decoration: BoxDecoration(
                              color: _selectedProvider == 'Airtel Money' ? Colors.red.withOpacity(0.1) : Theme.of(context).cardColor,
                              borderRadius: BorderRadius.circular(12),
                              border: Border.all(color: _selectedProvider == 'Airtel Money' ? Colors.red : const Color(0xFF2A3347), width: _selectedProvider == 'Airtel Money' ? 2 : 1),
                            ),
                            child: Center(
                              child: Text('Airtel Money', style: TextStyle(fontWeight: FontWeight.bold, color: _selectedProvider == 'Airtel Money' ? Colors.red : Colors.white)),
                            )
                          ),
                        )
                      ),
                    ],
                  ),

                  const SizedBox(height: 32),
                  Text('Amount to Deposit (UGX)', style: Theme.of(context).textTheme.titleMedium),
                  const SizedBox(height: 8),
                  TextField(
                    controller: _amountController,
                    keyboardType: TextInputType.number,
                    style: const TextStyle(fontSize: 24, fontWeight: FontWeight.bold, letterSpacing: 2),
                    decoration: const InputDecoration(
                      prefixText: 'UGX ',
                      prefixStyle: TextStyle(fontSize: 24, color: Colors.white54, fontWeight: FontWeight.normal),
                      hintText: '0',
                    ),
                  ),
                  
                  const SizedBox(height: 48),
                  ElevatedButton(
                    onPressed: (_isProcessing || _kibiinaId == null) ? null : _processPayment,
                    child: _isProcessing 
                      ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2))
                      : Text('Request payment via $_selectedProvider'),
                  ),
                ],
              ),
            ),
    );
  }
}
