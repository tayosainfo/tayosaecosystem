import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'core/theme/tayosa_theme.dart';
import 'core/router/app_router.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();
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