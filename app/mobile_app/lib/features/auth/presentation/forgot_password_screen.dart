import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:dio/dio.dart';
import '../../../core/network/api_client.dart';

class ForgotPasswordScreen extends ConsumerStatefulWidget {
  const ForgotPasswordScreen({super.key});

  @override
  ConsumerState<ForgotPasswordScreen> createState() => _ForgotPasswordScreenState();
}

class _ForgotPasswordScreenState extends ConsumerState<ForgotPasswordScreen> {
  final _email = TextEditingController();
  final _code = TextEditingController();
  bool _codeSent = false;
  bool _loading = false;
  String? _error;

  Future<void> _sendCode() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final apiClient = ref.read(apiClientInstanceProvider);
      await apiClient.requestPasswordReset(email: _email.text.trim());
      setState(() => _codeSent = true);
    } catch (e) {
      setState(() => _error = 'Failed to send reset code');
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  Future<void> _verifyCode() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final api = ref.read(apiClientProvider);
      final response = await api.post('/api/v1/auth/exchange-reset-password-token', data: {
        'email': _email.text.trim(),
        'code': _code.text.trim(),
      });
      final token = response.data['token']?.toString() ?? '';
      if (mounted) context.push('/reset-password?token=$token');
    } catch (e) {
      if (e is DioException && e.response?.data?['error'] != null) {
        setState(() => _error = e.response?.data?['error']);
      } else {
        setState(() => _error = 'Invalid or expired code');
      }
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Forgot Password')),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            TextField(
              controller: _email,
              decoration: const InputDecoration(labelText: 'Email'),
            ),
            const SizedBox(height: 12),
            if (_codeSent) ...[
              TextField(
                controller: _code,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(labelText: '6-digit code'),
              ),
              const SizedBox(height: 12),
            ],
            if (_error != null) Text(_error!, style: const TextStyle(color: Colors.red)),
            const SizedBox(height: 12),
            ElevatedButton(
              onPressed: _loading ? null : (_codeSent ? _verifyCode : _sendCode),
              child: Text(_codeSent ? 'Verify Code' : 'Send Reset Code'),
            ),
          ],
        ),
      ),
    );
  }
}

