import 'dart:ui';
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:lucide_icons/lucide_icons.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../../core/network/api_client.dart';
import '../data/geo_repository.dart';
import '../data/onboarding_geo_prefs.dart';

class GeoOnboardingScreen extends ConsumerStatefulWidget {
  const GeoOnboardingScreen({super.key});

  @override
  ConsumerState<GeoOnboardingScreen> createState() => _GeoOnboardingScreenState();
}

class _GeoOnboardingScreenState extends ConsumerState<GeoOnboardingScreen> {
  String? selectedDistrict;
  String? selectedCounty;
  String? selectedSubcounty;
  String? selectedParish;
  String? selectedVillage;
  bool _isSubmitting = false;
  String? _submitError;

  Future<void> _persistGeoAndSubmitPhase2() async {
    final d = selectedDistrict;
    final c = selectedCounty;
    final sc = selectedSubcounty;
    final p = selectedParish;
    final v = selectedVillage;
    if (d == null || c == null || sc == null || p == null || v == null) {
      setState(() => _submitError = 'Please select district through village.');
      return;
    }

    setState(() {
      _isSubmitting = true;
      _submitError = null;
    });

    try {
      final prefs = await SharedPreferences.getInstance();
      final userId = prefs.getString('user_id')?.trim() ?? '';
      final token = prefs.getString('auth_token')?.trim() ?? '';
      if (userId.isEmpty) {
        setState(() => _submitError = 'Missing user id. Register again, then return here.');
        return;
      }
      if (token.isEmpty) {
        setState(() => _submitError = 'Missing session. Register or sign in again so we can save your location.');
        return;
      }

      await prefs.setString(OnboardingGeoPrefs.district, d);
      await prefs.setString(OnboardingGeoPrefs.county, c);
      await prefs.setString(OnboardingGeoPrefs.subCounty, sc);
      await prefs.setString(OnboardingGeoPrefs.parish, p);
      await prefs.setString(OnboardingGeoPrefs.village, v);
      // Legacy keys used elsewhere in the app
      await prefs.setString('user_district', d);
      await prefs.setString('user_village', v);

      final api = ref.read(apiClientProvider);
      await api.post(
        '/api/v1/onboarding/phase',
        data: {
          'userId': userId,
          'phase': 2,
          'geo': {
            'district': d,
            'county': c,
            'sub_county': sc,
            'parish': p,
            'village': v,
          },
        },
      );

      if (mounted) context.go('/personal-onboarding');
    } on DioException catch (e) {
      final body = e.response?.data;
      String msg = 'Could not save location. Try again.';
      if (body is Map && body['error'] != null) {
        msg = body['error'].toString();
      } else if (e.message != null && e.message!.isNotEmpty) {
        msg = e.message!;
      }
      if (mounted) setState(() => _submitError = msg);
    } catch (e) {
      if (mounted) setState(() => _submitError = e.toString());
    } finally {
      if (mounted) setState(() => _isSubmitting = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final geoAsync = ref.watch(geoDataProvider);

    return Scaffold(
      body: Container(
        decoration: const BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [Color(0xFF1E3A8A), Color(0xFF1E40AF), Color(0xFF312E81)],
          ),
        ),
        child: Stack(
          children: [
            // Decorative orbs
            Positioned(top: -100, right: -100, child: _buildOrb(Colors.blue.withOpacity(0.15))),
            Positioned(bottom: -100, left: -100, child: _buildOrb(Colors.indigo.withOpacity(0.15))),

            SafeArea(
              child: Column(
                children: [
                  // App bar
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                    child: Row(
                      children: [
                        IconButton(
                          onPressed: () => context.pop(),
                          icon: const Icon(LucideIcons.arrowLeft, color: Colors.white),
                        ),
                        const Expanded(
                          child: Text(
                            'Bind Your Community',
                            textAlign: TextAlign.center,
                            style: TextStyle(color: Colors.white, fontSize: 20, fontWeight: FontWeight.bold),
                          ),
                        ),
                        const SizedBox(width: 48), // balance the back button
                      ],
                    ),
                  ),

                  // Progress indicator
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 32),
                    child: Row(
                      children: [
                        _buildStep(1, 'Register', true),
                        Expanded(child: Container(height: 2, color: Colors.white.withOpacity(0.3))),
                        _buildStep(2, 'Location', true),
                        Expanded(child: Container(height: 2, color: Colors.white.withOpacity(0.15))),
                        _buildStep(3, 'Profile', false),
                      ],
                    ),
                  ),

                  const SizedBox(height: 16),

                  // Content
                  Expanded(
                    child: geoAsync.when(
                      loading: () => const Center(child: CircularProgressIndicator(color: Colors.white)),
                      error: (err, stack) => Center(child: Text('Error: $err', style: const TextStyle(color: Colors.white))),
                      data: (geoTree) {
                        final districts = (geoTree.keys.toList()..sort());
                        final counties = selectedDistrict != null
                            ? (geoTree[selectedDistrict]!.keys.toList()..sort())
                            : <String>[];
                        final subcounties = selectedCounty != null
                            ? (geoTree[selectedDistrict]![selectedCounty]!.keys.toList()..sort())
                            : <String>[];
                        final parishes = selectedSubcounty != null
                            ? (geoTree[selectedDistrict]![selectedCounty]![selectedSubcounty]!.keys.toList()..sort())
                            : <String>[];
                        final villages = selectedParish != null
                            ? (List<String>.from(geoTree[selectedDistrict]![selectedCounty]![selectedSubcounty]![selectedParish]!)..sort())
                            : <String>[];

                        return SingleChildScrollView(
                          padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.stretch,
                            children: [
                              // Info banner
                              ClipRRect(
                                borderRadius: BorderRadius.circular(16),
                                child: BackdropFilter(
                                  filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10),
                                  child: Container(
                                    padding: const EdgeInsets.all(16),
                                    decoration: BoxDecoration(
                                      color: Colors.blue.withOpacity(0.15),
                                      borderRadius: BorderRadius.circular(16),
                                      border: Border.all(color: Colors.blue.withOpacity(0.3)),
                                    ),
                                    child: Row(
                                      children: [
                                        const Icon(LucideIcons.mapPin, color: Colors.blueAccent, size: 24),
                                        const SizedBox(width: 12),
                                        Expanded(
                                          child: Text(
                                            'Your location helps TAYOSA connect you with local services, community groups, and agents near you.',
                                            style: TextStyle(color: Colors.blue[100], fontSize: 14),
                                          ),
                                        ),
                                      ],
                                    ),
                                  ),
                                ),
                              ),
                              const SizedBox(height: 24),

                              _buildGlassDropdown('District', districts, selectedDistrict, (v) {
                                setState(() {
                                  selectedDistrict = v;
                                  selectedCounty = selectedSubcounty = selectedParish = selectedVillage = null;
                                });
                              }),
                              const SizedBox(height: 16),

                              _buildGlassDropdown('County / Municipality', counties, selectedCounty, (v) {
                                setState(() {
                                  selectedCounty = v;
                                  selectedSubcounty = selectedParish = selectedVillage = null;
                                });
                              }, disabled: selectedDistrict == null),
                              const SizedBox(height: 16),

                              _buildGlassDropdown('Subcounty / Division', subcounties, selectedSubcounty, (v) {
                                setState(() {
                                  selectedSubcounty = v;
                                  selectedParish = selectedVillage = null;
                                });
                              }, disabled: selectedCounty == null),
                              const SizedBox(height: 16),

                              _buildGlassDropdown('Parish / Ward', parishes, selectedParish, (v) {
                                setState(() {
                                  selectedParish = v;
                                  selectedVillage = null;
                                });
                              }, disabled: selectedSubcounty == null),
                              const SizedBox(height: 16),

                              _buildGlassDropdown('Village / Cell', villages, selectedVillage, (v) {
                                setState(() => selectedVillage = v);
                              }, disabled: selectedParish == null),

                              const SizedBox(height: 24),

                              if (_submitError != null)
                                Padding(
                                  padding: const EdgeInsets.only(bottom: 16),
                                  child: Material(
                                    color: Colors.red.shade900.withOpacity(0.35),
                                    borderRadius: BorderRadius.circular(12),
                                    child: Padding(
                                      padding: const EdgeInsets.all(12),
                                      child: Row(
                                        crossAxisAlignment: CrossAxisAlignment.start,
                                        children: [
                                          Icon(LucideIcons.alertCircle, color: Colors.red.shade100, size: 22),
                                          const SizedBox(width: 10),
                                          Expanded(
                                            child: Text(
                                              _submitError!,
                                              style: TextStyle(color: Colors.red.shade50, fontSize: 14),
                                            ),
                                          ),
                                        ],
                                      ),
                                    ),
                                  ),
                                ),

                              ElevatedButton(
                                onPressed: (selectedVillage == null || _isSubmitting) ? null : _persistGeoAndSubmitPhase2,
                                style: ElevatedButton.styleFrom(
                                  backgroundColor: Colors.white,
                                  foregroundColor: Colors.blue[900],
                                  disabledBackgroundColor: Colors.white.withOpacity(0.15),
                                  disabledForegroundColor: Colors.white.withOpacity(0.3),
                                  padding: const EdgeInsets.symmetric(vertical: 16),
                                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                                  elevation: 5,
                                ),
                                child: Row(
                                  mainAxisAlignment: MainAxisAlignment.center,
                                  children: [
                                    if (_isSubmitting)
                                      SizedBox(
                                        width: 20,
                                        height: 20,
                                        child: CircularProgressIndicator(color: Colors.blue[900], strokeWidth: 2),
                                      )
                                    else ...[
                                      const Text('Continue to Profile', style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
                                      const SizedBox(width: 8),
                                      const Icon(LucideIcons.arrowRight, size: 20),
                                    ],
                                  ],
                                ),
                              ),
                            ],
                          ),
                        );
                      },
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

  Widget _buildGlassDropdown(String label, List items, String? value, Function(String?) onChanged, {bool disabled = false}) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label, style: const TextStyle(color: Colors.white, fontSize: 14, fontWeight: FontWeight.w600)),
        const SizedBox(height: 8),
        Container(
          decoration: BoxDecoration(
            color: disabled ? Colors.white.withOpacity(0.05) : Colors.white.withOpacity(0.1),
            borderRadius: BorderRadius.circular(12),
            border: Border.all(color: Colors.white.withOpacity(disabled ? 0.1 : 0.3)),
          ),
          padding: const EdgeInsets.symmetric(horizontal: 16),
          child: DropdownButtonHideUnderline(
            child: DropdownButton<String>(
              isExpanded: true,
              value: value,
              hint: Text('Select $label', style: TextStyle(color: disabled ? Colors.white24 : Colors.blue[200])),
              dropdownColor: const Color(0xFF1E3A8A),
              style: const TextStyle(color: Colors.white),
              icon: Icon(LucideIcons.chevronDown, color: disabled ? Colors.white24 : Colors.blue[300]),
              items: disabled
                  ? []
                  : items.map((item) {
                      return DropdownMenuItem<String>(
                        value: item.toString(),
                        child: Text(item.toString()),
                      );
                    }).toList(),
              onChanged: disabled ? null : onChanged,
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildStep(int number, String label, bool active) {
    return Column(
      children: [
        Container(
          width: 28, height: 28,
          decoration: BoxDecoration(
            shape: BoxShape.circle,
            color: active ? Colors.white : Colors.white.withOpacity(0.15),
          ),
          child: Center(
            child: Text(
              '$number',
              style: TextStyle(
                color: active ? Colors.blue[900] : Colors.white54,
                fontWeight: FontWeight.bold,
                fontSize: 13,
              ),
            ),
          ),
        ),
        const SizedBox(height: 4),
        Text(label, style: TextStyle(color: active ? Colors.white : Colors.white54, fontSize: 11)),
      ],
    );
  }

  Widget _buildOrb(Color color) {
    return Container(width: 300, height: 300, decoration: BoxDecoration(shape: BoxShape.circle, color: color));
  }
}
