import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:lucide_icons/lucide_icons.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../../core/network/api_client.dart';
import '../data/onboarding_geo_prefs.dart';

class PersonalOnboardingScreen extends ConsumerStatefulWidget {
  const PersonalOnboardingScreen({super.key});

  @override
  ConsumerState<PersonalOnboardingScreen> createState() => _PersonalOnboardingScreenState();
}

class _PersonalOnboardingScreenState extends ConsumerState<PersonalOnboardingScreen> {
  final _dobController = TextEditingController();
  final _ninController = TextEditingController();
  final _idNumberController = TextEditingController();
  final _idExpiryController = TextEditingController();
  final _addressController = TextEditingController();
  final _occupationController = TextEditingController();
  final _referralCodeController = TextEditingController();
  final _referringAgentController = TextEditingController();

  String? _selectedGender;
  String? _selectedIdType;
  String? _selectedIndustry;
  String? _selectedMembershipType;
  bool _isLoading = false;
  String? _error;
  int _currentStep = 0;

  final List<String> _genders = ['Male', 'Female', 'Other'];
  final List<String> _idTypes = ['National ID (NIN)', 'Passport', 'Driving Permit', 'Voter ID'];
  final List<String> _industries = ['Agriculture', 'Banking & Finance', 'Education', 'Healthcare', 'Technology', 'Trade & Commerce', 'Transport', 'Manufacturing', 'Other'];
  final List<String> _membershipTypes = ['Individual', 'Joint', 'Corporate', 'Agent'];

  Future<void> _handleComplete() async {
    setState(() { _isLoading = true; _error = null; });

    try {
      final prefs = await SharedPreferences.getInstance();
      final userIdStr = prefs.getString('user_id') ?? '';
      final legacyInt = prefs.getInt('user_id');
      final userId = userIdStr.isNotEmpty ? userIdStr : (legacyInt != null ? '$legacyInt' : '');
      final api = ref.read(apiClientProvider);

      await api.post(
        '/api/v1/auth/onboard',
        data: {
          'userId': userId,
          'dateOfBirth': _dobController.text.trim(),
          'nin': _ninController.text.trim(),
          'gender': _selectedGender ?? '',
          'idType': _selectedIdType ?? '',
          'idNumber': _idNumberController.text.trim(),
          'idExpiry': _idExpiryController.text.trim(),
          'address': _addressController.text.trim(),
          'occupation': _occupationController.text.trim(),
          'industry': _selectedIndustry ?? '',
          'membershipType': _selectedMembershipType ?? 'Individual',
          'referralCode': _referralCodeController.text.trim(),
          'referringAgent': _referringAgentController.text.trim(),
          'region': '',
          'subRegion': '',
          'district': prefs.getString(OnboardingGeoPrefs.district) ?? prefs.getString('user_district') ?? '',
          'county': prefs.getString(OnboardingGeoPrefs.county) ?? '',
          'subCounty': prefs.getString(OnboardingGeoPrefs.subCounty) ?? '',
          'parish': prefs.getString(OnboardingGeoPrefs.parish) ?? '',
          'village': prefs.getString(OnboardingGeoPrefs.village) ?? prefs.getString('user_village') ?? '',
          'paymentMethod': 'Mobile Money',
          'paymentNumber': prefs.getString('user_phone') ?? '',
          'feeAmount': 0,
        },
      );

      if (mounted) context.go('/home');
    } catch (e) {
      print('Onboarding save error (continuing): $e');
      // Continue to home even if save fails for now
      if (mounted) context.go('/home');
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final steps = [
      _buildPersonalInfoStep(),
      _buildIdentityStep(),
      _buildOccupationStep(),
    ];

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
                          onPressed: () {
                            if (_currentStep > 0) {
                              setState(() => _currentStep--);
                            } else {
                              context.pop();
                            }
                          },
                          icon: const Icon(LucideIcons.arrowLeft, color: Colors.white),
                        ),
                        const Expanded(
                          child: Text('Complete Your Profile', textAlign: TextAlign.center,
                            style: TextStyle(color: Colors.white, fontSize: 20, fontWeight: FontWeight.bold)),
                        ),
                        const SizedBox(width: 48),
                      ],
                    ),
                  ),

                  // Progress
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 32),
                    child: Row(
                      children: [
                        _buildProgressStep(1, 'Register', true),
                        Expanded(child: Container(height: 2, color: Colors.white.withOpacity(0.3))),
                        _buildProgressStep(2, 'Location', true),
                        Expanded(child: Container(height: 2, color: Colors.white.withOpacity(0.3))),
                        _buildProgressStep(3, 'Profile', true),
                      ],
                    ),
                  ),
                  const SizedBox(height: 8),

                  // Sub-step indicator
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 32),
                    child: Row(
                      children: List.generate(3, (i) {
                        return Expanded(
                          child: Container(
                            height: 4,
                            margin: const EdgeInsets.symmetric(horizontal: 2),
                            decoration: BoxDecoration(
                              borderRadius: BorderRadius.circular(2),
                              color: i <= _currentStep ? Colors.white : Colors.white.withOpacity(0.15),
                            ),
                          ),
                        );
                      }),
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 32, vertical: 4),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text('Step ${_currentStep + 1} of 3', style: TextStyle(color: Colors.blue[200], fontSize: 12)),
                        Text(
                          ['Personal Info', 'Identity', 'Occupation & Membership'][_currentStep],
                          style: TextStyle(color: Colors.blue[200], fontSize: 12, fontWeight: FontWeight.w600),
                        ),
                      ],
                    ),
                  ),

                  // Content
                  Expanded(
                    child: SingleChildScrollView(
                      padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
                      child: ClipRRect(
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
                            child: steps[_currentStep],
                          ),
                        ),
                      ),
                    ),
                  ),

                  // Bottom buttons
                  Padding(
                    padding: const EdgeInsets.all(24),
                    child: Row(
                      children: [
                        if (_currentStep > 0)
                          Expanded(
                            child: OutlinedButton(
                              onPressed: () => setState(() => _currentStep--),
                              style: OutlinedButton.styleFrom(
                                foregroundColor: Colors.white,
                                side: const BorderSide(color: Colors.white54),
                                padding: const EdgeInsets.symmetric(vertical: 16),
                                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                              ),
                              child: const Text('Back'),
                            ),
                          ),
                        if (_currentStep > 0) const SizedBox(width: 12),
                        Expanded(
                          flex: 2,
                          child: ElevatedButton(
                            onPressed: _isLoading
                                ? null
                                : () {
                                    if (_currentStep < 2) {
                                      setState(() => _currentStep++);
                                    } else {
                                      _handleComplete();
                                    }
                                  },
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
                                    children: [
                                      Text(
                                        _currentStep < 2 ? 'Continue' : 'Complete Setup',
                                        style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
                                      ),
                                      const SizedBox(width: 8),
                                      const Icon(LucideIcons.arrowRight, size: 20),
                                    ],
                                  ),
                          ),
                        ),
                      ],
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

  Widget _buildPersonalInfoStep() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        _buildGlassInput(controller: _dobController, label: 'Date of Birth', hint: '1990-01-15', icon: LucideIcons.calendar, keyboardType: TextInputType.datetime),
        const SizedBox(height: 20),
        _buildGlassDropdownField('Gender', _genders, _selectedGender, (v) => setState(() => _selectedGender = v)),
        const SizedBox(height: 20),
        _buildGlassInput(controller: _addressController, label: 'Address', hint: 'Your physical address', icon: LucideIcons.home),
      ],
    );
  }

  Widget _buildIdentityStep() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        _buildGlassInput(controller: _ninController, label: 'National ID Number (NIN)', hint: 'CM8XXXXXXXX', icon: LucideIcons.fingerprint),
        const SizedBox(height: 20),
        _buildGlassDropdownField('ID Document Type', _idTypes, _selectedIdType, (v) => setState(() => _selectedIdType = v)),
        const SizedBox(height: 20),
        _buildGlassInput(controller: _idNumberController, label: 'ID Document Number', hint: 'Document number', icon: LucideIcons.creditCard),
        const SizedBox(height: 20),
        _buildGlassInput(controller: _idExpiryController, label: 'ID Expiry Date', hint: '2030-12-31', icon: LucideIcons.calendarDays, keyboardType: TextInputType.datetime),
      ],
    );
  }

  Widget _buildOccupationStep() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        _buildGlassInput(controller: _occupationController, label: 'Occupation', hint: 'e.g. Software Engineer, Farmer', icon: LucideIcons.briefcase),
        const SizedBox(height: 20),
        _buildGlassDropdownField('Industry', _industries, _selectedIndustry, (v) => setState(() => _selectedIndustry = v)),
        const SizedBox(height: 20),
        _buildGlassDropdownField('Membership Type', _membershipTypes, _selectedMembershipType, (v) => setState(() => _selectedMembershipType = v)),
        const SizedBox(height: 20),
        _buildGlassInput(controller: _referralCodeController, label: 'Referral Code (Optional)', hint: 'Enter referral code', icon: LucideIcons.gift),
        const SizedBox(height: 20),
        _buildGlassInput(controller: _referringAgentController, label: 'Referring Agent (Optional)', hint: 'Agent name or code', icon: LucideIcons.userCheck),
      ],
    );
  }

  Widget _buildGlassInput({
    required TextEditingController controller,
    required String label,
    required String hint,
    required IconData icon,
    TextInputType keyboardType = TextInputType.text,
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label, style: const TextStyle(color: Colors.white, fontSize: 14, fontWeight: FontWeight.w500)),
        const SizedBox(height: 8),
        TextField(
          controller: controller,
          keyboardType: keyboardType,
          style: const TextStyle(color: Colors.white),
          decoration: InputDecoration(
            prefixIcon: Icon(icon, color: Colors.blue[300]),
            hintText: hint,
            hintStyle: TextStyle(color: Colors.blue[200]),
            filled: true,
            fillColor: Colors.white.withOpacity(0.1),
            border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide(color: Colors.white.withOpacity(0.3))),
            enabledBorder: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide(color: Colors.white.withOpacity(0.3))),
            focusedBorder: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: const BorderSide(color: Colors.blueAccent, width: 2)),
          ),
        ),
      ],
    );
  }

  Widget _buildGlassDropdownField(String label, List<String> items, String? value, Function(String?) onChanged) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label, style: const TextStyle(color: Colors.white, fontSize: 14, fontWeight: FontWeight.w500)),
        const SizedBox(height: 8),
        Container(
          decoration: BoxDecoration(
            color: Colors.white.withOpacity(0.1),
            borderRadius: BorderRadius.circular(12),
            border: Border.all(color: Colors.white.withOpacity(0.3)),
          ),
          padding: const EdgeInsets.symmetric(horizontal: 16),
          child: DropdownButtonHideUnderline(
            child: DropdownButton<String>(
              isExpanded: true,
              value: value,
              hint: Text('Select $label', style: TextStyle(color: Colors.blue[200])),
              dropdownColor: const Color(0xFF1E3A8A),
              style: const TextStyle(color: Colors.white),
              icon: Icon(LucideIcons.chevronDown, color: Colors.blue[300]),
              items: items.map((item) => DropdownMenuItem<String>(value: item, child: Text(item))).toList(),
              onChanged: onChanged,
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildProgressStep(int number, String label, bool active) {
    return Column(
      children: [
        Container(
          width: 28, height: 28,
          decoration: BoxDecoration(shape: BoxShape.circle, color: active ? Colors.white : Colors.white.withOpacity(0.15)),
          child: Center(child: Text('$number', style: TextStyle(color: active ? Colors.blue[900] : Colors.white54, fontWeight: FontWeight.bold, fontSize: 13))),
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
