import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../features/splash/presentation/splash_screen.dart';
import '../../features/auth/presentation/login_screen.dart';
import '../../features/auth/presentation/register_screen.dart';
import '../../features/auth/presentation/forgot_password_screen.dart';
import '../../features/auth/presentation/reset_password_screen.dart';
import '../../features/onboarding/presentation/geo_onboarding_screen.dart';
import '../../features/onboarding/presentation/personal_onboarding_screen.dart';
import '../../features/home/presentation/home_dashboard_screen.dart';
import '../../features/onboarding/presentation/kibiina_registration_screen.dart';
import '../../features/kibiina/presentation/kibiina_dashboard_screen.dart';
import '../../features/kibiina/presentation/kibiina_deposit_screen.dart';

final routerProvider = Provider<GoRouter>((ref) {
  return GoRouter(
    initialLocation: '/',
    routes: [
      // Splash (app entry)
      GoRoute(
        path: '/',
        builder: (context, state) => const SplashScreen(),
      ),

      // Auth
      GoRoute(
        path: '/login',
        builder: (context, state) => const LoginScreen(),
      ),
      GoRoute(
        path: '/register',
        builder: (context, state) => const RegisterScreen(),
      ),
      GoRoute(
        path: '/forgot-password',
        builder: (context, state) => const ForgotPasswordScreen(),
      ),
      GoRoute(
        path: '/reset-password',
        builder: (context, state) => ResetPasswordScreen(
          token: state.uri.queryParameters['token'] ?? '',
        ),
      ),

      // Global Onboarding (Registration Flow)
      GoRoute(
        path: '/geo-onboarding',
        builder: (context, state) => const GeoOnboardingScreen(),
      ),
      GoRoute(
        path: '/personal-onboarding',
        builder: (context, state) => const PersonalOnboardingScreen(),
      ),

      // Main App
      GoRoute(
        path: '/home',
        builder: (context, state) => const HomeDashboardScreen(),
      ),

      // Kibiina (Service-specific, accessed from Dashboard)
      GoRoute(
        path: '/kibiina-register',
        builder: (context, state) => const KibiinaRegistrationScreen(),
      ),
      GoRoute(
        path: '/kibiina-dashboard',
        builder: (context, state) => const KibiinaDashboardScreen(),
      ),
      GoRoute(
        path: '/kibiina-deposit',
        builder: (context, state) => const KibiinaDepositScreen(),
      ),
    ],
  );
});
