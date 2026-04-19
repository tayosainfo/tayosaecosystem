import 'dart:math';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class SplashScreen extends StatefulWidget {
  const SplashScreen({super.key});

  @override
  State<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends State<SplashScreen> with TickerProviderStateMixin {
  late AnimationController _logoController;
  late AnimationController _pulseController;
  late AnimationController _particleController;
  late AnimationController _textController;
  late Animation<double> _logoScale;
  late Animation<double> _logoOpacity;
  late Animation<double> _pulseScale;
  late Animation<double> _textSlide;
  late Animation<double> _textOpacity;

  @override
  void initState() {
    super.initState();

    // Logo animation: scale + fade in
    _logoController = AnimationController(vsync: this, duration: const Duration(milliseconds: 1200));
    _logoScale = Tween<double>(begin: 0.3, end: 1.0).animate(CurvedAnimation(parent: _logoController, curve: Curves.elasticOut));
    _logoOpacity = Tween<double>(begin: 0.0, end: 1.0).animate(CurvedAnimation(parent: _logoController, curve: const Interval(0.0, 0.5, curve: Curves.easeIn)));

    // Pulse ring animation
    _pulseController = AnimationController(vsync: this, duration: const Duration(milliseconds: 2000))..repeat();
    _pulseScale = Tween<double>(begin: 1.0, end: 2.5).animate(CurvedAnimation(parent: _pulseController, curve: Curves.easeOut));

    // Particle animation
    _particleController = AnimationController(vsync: this, duration: const Duration(seconds: 4))..repeat();

    // Text animation
    _textController = AnimationController(vsync: this, duration: const Duration(milliseconds: 800));
    _textSlide = Tween<double>(begin: 30.0, end: 0.0).animate(CurvedAnimation(parent: _textController, curve: Curves.easeOutCubic));
    _textOpacity = Tween<double>(begin: 0.0, end: 1.0).animate(CurvedAnimation(parent: _textController, curve: Curves.easeIn));

    // Sequence the animations
    _logoController.forward();
    Future.delayed(const Duration(milliseconds: 800), () {
      if (mounted) _textController.forward();
    });

    // Navigate after animations complete
    Future.delayed(const Duration(milliseconds: 3500), () {
      if (mounted) context.go('/login');
    });
  }

  @override
  void dispose() {
    _logoController.dispose();
    _pulseController.dispose();
    _particleController.dispose();
    _textController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        decoration: const BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [Color(0xFF0B1120), Color(0xFF1E3A8A), Color(0xFF312E81)],
            stops: [0.0, 0.5, 1.0],
          ),
        ),
        child: Stack(
          children: [
            // Animated particles
            AnimatedBuilder(
              animation: _particleController,
              builder: (context, _) => CustomPaint(
                size: MediaQuery.of(context).size,
                painter: _ParticlePainter(_particleController.value),
              ),
            ),

            // Center logo + text
            Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  // Pulse rings
                  AnimatedBuilder(
                    animation: _pulseController,
                    builder: (context, child) => Stack(
                      alignment: Alignment.center,
                      children: [
                        // Ring 1
                        Opacity(
                          opacity: (1 - _pulseController.value).clamp(0.0, 0.5),
                          child: Container(
                            width: 120 * _pulseScale.value,
                            height: 120 * _pulseScale.value,
                            decoration: BoxDecoration(
                              shape: BoxShape.circle,
                              border: Border.all(color: Colors.blueAccent.withOpacity(0.3), width: 2),
                            ),
                          ),
                        ),
                        // Logo
                        AnimatedBuilder(
                          animation: _logoController,
                          builder: (context, child) => Opacity(
                            opacity: _logoOpacity.value,
                            child: Transform.scale(
                              scale: _logoScale.value,
                              child: Container(
                                width: 120,
                                height: 120,
                                decoration: BoxDecoration(
                                  shape: BoxShape.circle,
                                  color: Colors.white,
                                  boxShadow: [
                                    BoxShadow(color: Colors.blueAccent.withOpacity(0.6), blurRadius: 40, spreadRadius: 10),
                                    BoxShadow(color: Colors.purpleAccent.withOpacity(0.3), blurRadius: 60, spreadRadius: 15),
                                  ],
                                ),
                                child: ClipOval(
                                  child: Image.asset(
                                    'assets/images/logo.png',
                                    fit: BoxFit.cover,
                                    errorBuilder: (_, __, ___) => const Center(
                                      child: Text('T', style: TextStyle(color: Color(0xFF1E3A8A), fontSize: 56, fontWeight: FontWeight.w900)),
                                    ),
                                  ),
                                ),
                              ),
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),

                  const SizedBox(height: 48),

                  // Text
                  AnimatedBuilder(
                    animation: _textController,
                    builder: (context, child) => Opacity(
                      opacity: _textOpacity.value,
                      child: Transform.translate(
                        offset: Offset(0, _textSlide.value),
                        child: Column(
                          children: [
                            const Text(
                              'TAYOSA',
                              style: TextStyle(
                                color: Colors.white,
                                fontSize: 44,
                                fontWeight: FontWeight.w900,
                                letterSpacing: 8,
                              ),
                            ),
                            const SizedBox(height: 8),
                            ShaderMask(
                              shaderCallback: (bounds) => const LinearGradient(
                                colors: [Colors.blueAccent, Colors.purpleAccent, Colors.amberAccent],
                              ).createShader(bounds),
                              child: const Text(
                                'Private Equity Platform',
                                style: TextStyle(color: Colors.white, fontSize: 16, fontWeight: FontWeight.w500, letterSpacing: 3),
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),

                  const SizedBox(height: 80),

                  // Loading indicator
                  AnimatedBuilder(
                    animation: _textController,
                    builder: (context, _) => Opacity(
                      opacity: _textOpacity.value,
                      child: SizedBox(
                        width: 32,
                        height: 32,
                        child: CircularProgressIndicator(
                          strokeWidth: 2,
                          valueColor: AlwaysStoppedAnimation(Colors.blueAccent.withOpacity(0.6)),
                        ),
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
}

class _ParticlePainter extends CustomPainter {
  final double animationValue;
  final int particleCount = 50;
  final Random _random = Random(42); // Seeded for consistency

  _ParticlePainter(this.animationValue);

  @override
  void paint(Canvas canvas, Size size) {
    for (int i = 0; i < particleCount; i++) {
      final baseX = _random.nextDouble() * size.width;
      final baseY = _random.nextDouble() * size.height;
      final speed = 0.3 + _random.nextDouble() * 0.7;
      final radius = 1.0 + _random.nextDouble() * 2.5;

      final y = (baseY - animationValue * speed * size.height) % size.height;
      final opacity = (0.15 + _random.nextDouble() * 0.35) * (1 - (y / size.height) * 0.5);

      final paint = Paint()
        ..color = Colors.blueAccent.withOpacity(opacity.clamp(0.0, 1.0))
        ..style = PaintingStyle.fill;

      canvas.drawCircle(Offset(baseX, y), radius, paint);
    }
  }

  @override
  bool shouldRepaint(covariant _ParticlePainter oldDelegate) => true;
}
