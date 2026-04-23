import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'core/theme/tayosa_theme.dart';
import 'core/router/app_router.dart';
import 'core/network/supabase_client.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  try {
    await SupabaseConfig.initialize();
  } catch (e) {
    // Log the error and continue - app can still function with limited features
    debugPrint('Failed to initialize Supabase: $e');
  }
  
  runApp(
    const ProviderScope(
      child: TayosaApp(),
    ),
  );
}

class TayosaApp extends ConsumerWidget {
  const TayosaApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = ref.watch(routerProvider);
    
    return MaterialApp.router(
      title: 'TAYOSA Mobile',
      debugShowCheckedModeBanner: false,
      theme: TayosaTheme.darkTheme,
      routerConfig: router,
    );
  }
}