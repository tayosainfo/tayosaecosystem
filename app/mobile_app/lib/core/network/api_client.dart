import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:flutter/foundation.dart';

const String _apiBaseOverride = String.fromEnvironment('API_BASE_URL', defaultValue: '');
const String insforgeBaseUrl = String.fromEnvironment(
  'INSFORGE_BASE_URL',
  defaultValue: 'https://74qj9u5z.us-east.insforge.app',
);
const String insforgeAnonKey = String.fromEnvironment(
  'INSFORGE_ANON_KEY',
  defaultValue: '',
);

// Provides the base API client configured for local Android Emulator vs Web
final apiClientProvider = Provider<Dio>((ref) {
  // Use 10.0.2.2 for Android Emulator, 127.0.0.1 for Web/iOS Simulator
  final String baseUrl = _apiBaseOverride.isNotEmpty
      ? _apiBaseOverride
      : (kIsWeb ? 'http://127.0.0.1:8080' : 'http://10.0.2.2:8080');

  final options = BaseOptions(
    baseUrl: baseUrl, 
    connectTimeout: const Duration(seconds: 15),
    receiveTimeout: const Duration(seconds: 15),
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    },
  );

  final dio = Dio(options);

  // Add Auth interceptor
  dio.interceptors.add(InterceptorsWrapper(
    onRequest: (options, handler) async {
      final prefs = await SharedPreferences.getInstance();
      final token = prefs.getString('auth_token');
      if (token != null && token.isNotEmpty) {
        options.headers['Authorization'] = 'Bearer $token';
      }
      return handler.next(options);
    },
  ));

  // Add logging interceptor
  dio.interceptors.add(LogInterceptor(
    request: true,
    requestBody: true,
    responseBody: true,
    error: true,
  ));

  return dio;
});
