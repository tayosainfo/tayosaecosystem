import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:lucide_icons/lucide_icons.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:dio/dio.dart';
import '../../../core/network/api_client.dart';

class LoginScreen extends ConsumerStatefulWidget {
  const LoginScreen({super.key});

  @override
  ConsumerState<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends ConsumerState<LoginScreen> {
  bool _isPhoneMethod = true;
  final _identifierController = TextEditingController();
  final _passwordController = TextEditingController();
  bool _isLoading = false;
  bool _obscurePassword = true;
  String? _error;

  Future<void> _handleLogin() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });

    final identifierText = _identifierController.text.trim();
    final passwordText = _passwordController.text;

    if (identifierText.isEmpty) {
      setState(() {
        _error = 'Please enter your ${_isPhoneMethod ? 'phone' : 'email'}';
        _isLoading = false;
      });
      return;
    }

    if (passwordText.isEmpty) {
      setState(() {
        _error = 'Please enter your password';
        _isLoading = false;
      });
      return;
    }

    try {
      final api = ref.read(apiClientProvider);
      final response = await api.post(
        '/api/v1/auth/login',
        data: {
          'identifier': identifierText,
          'password': passwordText,
        },
      );

      if (response.statusCode == 200 || response.statusCode == 201) {
        final prefs = await SharedPreferences.getInstance();
        if (response.data['user'] != null) {
          final user = response.data['user'];
          await prefs.setString('user_id', user['id']?.toString() ?? '');
          await prefs.setString('user_name', user['fullName'] ?? 'Member');
          await prefs.setString('user_email', user['contactEmail'] ?? '');
          await prefs.setString('user_phone', user['phoneE164'] ?? '');
          await prefs.setString('auth_token', response.data['session']?['accessToken'] ?? '');
        }
        if (mounted) context.go('/home');
      } else {
        setState(() => _error = response.data['error'] ?? 'Login failed');
      }
    } catch (e) {
      print('Login error: $e');
      if (e is DioException && e.response?.data != null && e.response?.data['error'] != null) {
        setState(() => _error = e.response?.data['error']);
      } else {
        setState(() => _error = 'Network or server error occurred. Please try again.');
      }
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
            stops: [0.0, 0.5, 1.0],
          ),
        ),
        child: Stack(
          children: [
            // Background Orbs
            Positioned(
              top: -100, right: -100,
              child: Container(
                width: 300, height: 300,
                decoration: BoxDecoration(shape: BoxShape.circle, color: Colors.blue.withOpacity(0.2)),
              ).blurred(),
            ),
            Positioned(
              bottom: -100, left: -100,
              child: Container(
                width: 300, height: 300,
                decoration: BoxDecoration(shape: BoxShape.circle, color: Colors.indigo.withOpacity(0.2)),
              ).blurred(),
            ),
            
            SafeArea(
              child: Center(
                child: SingleChildScrollView(
                  padding: const EdgeInsets.symmetric(horizontal: 24.0, vertical: 48),
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      // Header
                      Column(
                        children: [
                          Container(
                            padding: const EdgeInsets.all(16),
                            decoration: const BoxDecoration(
                              shape: BoxShape.circle,
                              color: Colors.white,
                              boxShadow: [
                                BoxShadow(color: Colors.blueAccent, blurRadius: 30, spreadRadius: 5)
                              ]
                            ),
                            child: const Icon(LucideIcons.shield, color: Colors.blue, size: 48),
                          ),
                          const SizedBox(height: 24),
                          const Text(
                            'Welcome Back',
                            style: TextStyle(color: Colors.white, fontSize: 36, fontWeight: FontWeight.bold),
                          ),
                          const SizedBox(height: 8),
                          Text(
                            'Sign in to your TAYOSA account',
                            style: TextStyle(color: Colors.blue[100], fontSize: 18),
                          ),
                          const SizedBox(height: 16),
                          Row(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              const Icon(LucideIcons.sparkles, color: Colors.amber, size: 16),
                              const SizedBox(width: 8),
                              Text('Private Equity Platform', style: TextStyle(color: Colors.blue[200], fontSize: 14)),
                              const SizedBox(width: 8),
                              const Icon(LucideIcons.sparkles, color: Colors.amber, size: 16),
                            ],
                          ),
                        ],
                      ),
                      const SizedBox(height: 48),

                      // Form Container
                      ClipRRect(
                        borderRadius: BorderRadius.circular(24),
                        child: BackdropFilter(
                          filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10),
                          child: Container(
                            padding: const EdgeInsets.all(24),
                            decoration: BoxDecoration(
                              color: Colors.white.withOpacity(0.1),
                              borderRadius: BorderRadius.circular(24),
                              border: Border.all(color: Colors.white.withOpacity(0.2)),
                            ),
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.stretch,
                              children: [
                                // Toggles
                                Container(
                                  padding: const EdgeInsets.all(4),
                                  decoration: BoxDecoration(
                                    color: Colors.indigo.withOpacity(0.5),
                                    borderRadius: BorderRadius.circular(12),
                                  ),
                                  child: Row(
                                    children: [
                                      Expanded(
                                        child: GestureDetector(
                                          onTap: () => setState(() => _isPhoneMethod = true),
                                          child: Container(
                                            padding: const EdgeInsets.symmetric(vertical: 12),
                                            decoration: BoxDecoration(
                                              color: _isPhoneMethod ? Colors.white : Colors.transparent,
                                              borderRadius: BorderRadius.circular(8),
                                            ),
                                            child: Text(
                                              'Phone',
                                              textAlign: TextAlign.center,
                                              style: TextStyle(
                                                color: _isPhoneMethod ? Colors.indigo[900] : Colors.blue[100],
                                                fontWeight: FontWeight.bold,
                                              ),
                                            ),
                                          ),
                                        ),
                                      ),
                                      Expanded(
                                        child: GestureDetector(
                                          onTap: () => setState(() => _isPhoneMethod = false),
                                          child: Container(
                                            padding: const EdgeInsets.symmetric(vertical: 12),
                                            decoration: BoxDecoration(
                                              color: !_isPhoneMethod ? Colors.white : Colors.transparent,
                                              borderRadius: BorderRadius.circular(8),
                                            ),
                                            child: Text(
                                              'Email',
                                              textAlign: TextAlign.center,
                                              style: TextStyle(
                                                color: !_isPhoneMethod ? Colors.indigo[900] : Colors.blue[100],
                                                fontWeight: FontWeight.bold,
                                              ),
                                            ),
                                          ),
                                        ),
                                      ),
                                    ],
                                  ),
                                ),
                                const SizedBox(height: 24),

                                // Inputs
                                Text(
                                  _isPhoneMethod ? 'Phone Number' : 'Email Address',
                                  style: const TextStyle(color: Colors.white, fontSize: 14, fontWeight: FontWeight.w500),
                                ),
                                const SizedBox(height: 8),
                                TextField(
                                  controller: _identifierController,
                                  keyboardType: _isPhoneMethod ? TextInputType.phone : TextInputType.emailAddress,
                                  style: const TextStyle(color: Colors.white),
                                  decoration: InputDecoration(
                                    prefixIcon: Icon(_isPhoneMethod ? LucideIcons.phone : LucideIcons.mail, color: Colors.blue[300]),
                                    hintText: _isPhoneMethod ? '700 123 456' : 'you@example.com',
                                    hintStyle: TextStyle(color: Colors.blue[200]),
                                    filled: true,
                                    fillColor: Colors.white.withOpacity(0.1),
                                    border: OutlineInputBorder(
                                      borderRadius: BorderRadius.circular(12),
                                      borderSide: BorderSide(color: Colors.white.withOpacity(0.3)),
                                    ),
                                    enabledBorder: OutlineInputBorder(
                                      borderRadius: BorderRadius.circular(12),
                                      borderSide: BorderSide(color: Colors.white.withOpacity(0.3)),
                                    ),
                                    focusedBorder: OutlineInputBorder(
                                      borderRadius: BorderRadius.circular(12),
                                      borderSide: const BorderSide(color: Colors.blueAccent),
                                    ),
                                  ),
                                ),
                                const SizedBox(height: 20),

                                Row(
                                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                  children: [
                                    const Text('Password', style: TextStyle(color: Colors.white, fontSize: 14, fontWeight: FontWeight.w500)),
                                    GestureDetector(
                                      onTap: () => context.push('/forgot-password'),
                                      child: Text('Forgot password?', style: TextStyle(color: Colors.blue[300], fontSize: 14)),
                                    ),
                                  ],
                                ),
                                const SizedBox(height: 8),
                                TextField(
                                  controller: _passwordController,
                                  obscureText: _obscurePassword,
                                  style: const TextStyle(color: Colors.white),
                                  decoration: InputDecoration(
                                    prefixIcon: Icon(LucideIcons.lock, color: Colors.blue[300]),
                                    suffixIcon: IconButton(
                                      icon: Icon(_obscurePassword ? LucideIcons.eyeOff : LucideIcons.eye, color: Colors.blue[300]),
                                      onPressed: () => setState(() => _obscurePassword = !_obscurePassword),
                                    ),
                                    hintText: 'Enter your password',
                                    hintStyle: TextStyle(color: Colors.blue[200]),
                                    filled: true,
                                    fillColor: Colors.white.withOpacity(0.1),
                                    border: OutlineInputBorder(
                                      borderRadius: BorderRadius.circular(12),
                                      borderSide: BorderSide(color: Colors.white.withOpacity(0.3)),
                                    ),
                                    enabledBorder: OutlineInputBorder(
                                      borderRadius: BorderRadius.circular(12),
                                      borderSide: BorderSide(color: Colors.white.withOpacity(0.3)),
                                    ),
                                    focusedBorder: OutlineInputBorder(
                                      borderRadius: BorderRadius.circular(12),
                                      borderSide: const BorderSide(color: Colors.blueAccent),
                                    ),
                                  ),
                                ),
                                const SizedBox(height: 24),

                                if (_error != null)
                                  Container(
                                    padding: const EdgeInsets.all(12),
                                    margin: const EdgeInsets.only(bottom: 24),
                                    decoration: BoxDecoration(
                                      color: Colors.red[500]!.withOpacity(0.2),
                                      border: Border.all(color: Colors.red[400]!.withOpacity(0.5)),
                                      borderRadius: BorderRadius.circular(12),
                                    ),
                                    child: Text(_error!, style: TextStyle(color: Colors.red[100])),
                                  ),

                                ElevatedButton(
                                  onPressed: _isLoading ? null : _handleLogin,
                                  style: ElevatedButton.styleFrom(
                                    backgroundColor: Colors.white,
                                    foregroundColor: Colors.blue[900],
                                    padding: const EdgeInsets.symmetric(vertical: 16),
                                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                                    elevation: 5,
                                  ),
                                  child: _isLoading 
                                    ? SizedBox(width: 20, height: 20, child: CircularProgressIndicator(color: Colors.blue[900], strokeWidth: 2))
                                    : Row(
                                        mainAxisAlignment: MainAxisAlignment.center,
                                        children: const [
                                          Text('Sign In', style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
                                          SizedBox(width: 8),
                                          Icon(LucideIcons.arrowRight, size: 20),
                                        ],
                                      ),
                                ),

                                const SizedBox(height: 24),
                                Row(
                                  mainAxisAlignment: MainAxisAlignment.center,
                                  children: [
                                    Text('New to TAYOSA? ', style: TextStyle(color: Colors.blue[200])),
                                    GestureDetector(
                                      onTap: () => context.push('/register'),
                                      child: const Text('Create account', style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold)),
                                    ),
                                  ],
                                ),
                              ],
                            ),
                          ),
                        ),
                      ),
                      
                      const SizedBox(height: 32),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Icon(LucideIcons.shieldCheck, color: Colors.blue[300], size: 16),
                          const SizedBox(width: 4),
                          Text('Bank-grade encryption • 24/7 secure', style: TextStyle(color: Colors.blue[300], fontSize: 12)),
                        ],
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

// Help extension for blur
extension BlurExtension on Widget {
  Widget blurred({double blur = 40}) {
    return ImageFilter.blur(sigmaX: blur, sigmaY: blur) != null 
        ? BackdropFilter(filter: ImageFilter.blur(sigmaX: blur, sigmaY: blur), child: this)
        : this;
  }
}
