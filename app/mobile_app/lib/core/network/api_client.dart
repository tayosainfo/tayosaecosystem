import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:flutter/foundation.dart';
import 'supabase_client.dart';

const String _apiBaseOverride = String.fromEnvironment('API_BASE_URL', defaultValue: '');

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

  // Add Supabase Auth interceptor
  dio.interceptors.add(InterceptorsWrapper(
    onRequest: (options, handler) async {
      // Get current session from Supabase
      final session = SupabaseConfig.client.auth.currentSession;
      if (session != null) {
        options.headers['Authorization'] = 'Bearer ${session.accessToken}';
      } else {
        // Fallback to SharedPreferences for backward compatibility
        final prefs = await SharedPreferences.getInstance();
        final token = prefs.getString('auth_token');
        if (token != null && token.isNotEmpty) {
          options.headers['Authorization'] = 'Bearer $token';
        }
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

/// API Client class that integrates with Supabase authentication
class ApiClient {
  final Dio _dio;

  ApiClient(this._dio);

  /// Register a new user through the backend API
  /// The backend handles Supabase signup internally
  Future<Map<String, dynamic>> register({
    required String email,
    required String password,
    required String fullName,
    required String phone,
  }) async {
    final response = await _dio.post('/api/v1/auth/register', data: {
      'email': email,
      'password': password,
      'fullName': fullName,
      'phone': phone,
    });
    return response.data;
  }

  /// Login user through the backend API
  /// The backend handles Supabase authentication internally
  Future<Map<String, dynamic>> login({
    required String identifier,
    required String password,
  }) async {
    final response = await _dio.post('/api/v1/auth/login', data: {
      'identifier': identifier,
      'password': password,
    });
    return response.data;
  }

  /// Logout user through the backend API
  /// The backend handles Supabase session invalidation
  Future<void> logout() async {
    await _dio.post('/api/v1/auth/logout');
  }

  /// Request password reset through the backend API
  /// The backend handles Supabase password reset email
  Future<void> requestPasswordReset({required String email}) async {
    await _dio.post('/api/v1/auth/forgot-password', data: {
      'email': email,
    });
  }

  /// Reset password through the backend API
  /// The backend handles Supabase password reset completion
  Future<void> resetPassword({
    required String token,
    required String newPassword,
  }) async {
    await _dio.post('/api/v1/auth/reset-password', data: {
      'token': token,
      'newPassword': newPassword,
    });
  }

  /// Get current user profile
  Future<Map<String, dynamic>> getCurrentUser() async {
    final response = await _dio.get('/api/v1/users/me');
    return response.data;
  }

  /// Update user profile
  Future<Map<String, dynamic>> updateProfile(Map<String, dynamic> data) async {
    final response = await _dio.put('/api/v1/users/me', data: data);
    return response.data;
  }
}

/// Provider for the ApiClient instance
final apiClientInstanceProvider = Provider<ApiClient>((ref) {
  final dio = ref.watch(apiClientProvider);
  return ApiClient(dio);
});
