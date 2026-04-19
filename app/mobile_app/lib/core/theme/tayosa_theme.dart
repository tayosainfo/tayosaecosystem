import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

class TayosaTheme {
  static const Color background = Color(0xFF0B1120);
  static const Color surface = Color(0xFF131A2A);
  static const Color surfaceBorder = Color(0xFF2A3347);
  static const Color emerald500 = Color(0xFF10B981);
  static const Color emerald400 = Color(0xFF34D399);
  static const Color blue500 = Color(0xFF3B82F6);
  static const Color rose500 = Color(0xFFF43F5E);
  static const Color amber500 = Color(0xFFF59E0B);
  
  static const Color textPrimary = Colors.white;
  static const Color textSecondary = Color(0xFF94A3B8); // Slate 400

  static ThemeData get darkTheme {
    return ThemeData(
      useMaterial3: true,
      brightness: Brightness.dark,
      scaffoldBackgroundColor: background,
      primaryColor: emerald500,
      colorScheme: const ColorScheme.dark(
        primary: emerald500,
        secondary: blue500,
        surface: surface,
        background: background,
        error: rose500,
      ),
      textTheme: GoogleFonts.dmSansTextTheme(ThemeData.dark().textTheme).copyWith(
        displayLarge: GoogleFonts.outfit(fontWeight: FontWeight.bold, color: textPrimary),
        displayMedium: GoogleFonts.outfit(fontWeight: FontWeight.bold, color: textPrimary),
        titleLarge: GoogleFonts.outfit(fontWeight: FontWeight.w600, color: textPrimary),
        bodyLarge: GoogleFonts.dmSans(color: textPrimary),
        bodyMedium: GoogleFonts.dmSans(color: textSecondary),
      ),
      cardTheme: const CardThemeData(
        color: surface,
        elevation: 0,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.all(Radius.circular(16)),
          side: BorderSide(color: surfaceBorder),
        ),
      ),
      inputDecorationTheme: InputDecorationTheme(
        filled: true,
        fillColor: const Color(0xFF0F1623),
        contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: const BorderSide(color: surfaceBorder),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: const BorderSide(color: emerald500, width: 2),
        ),
        hintStyle: GoogleFonts.dmSans(color: textSecondary),
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          backgroundColor: emerald500,
          foregroundColor: Colors.white,
          elevation: 0,
          padding: const EdgeInsets.symmetric(vertical: 16, horizontal: 24),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(12),
          ),
          textStyle: GoogleFonts.dmSans(fontWeight: FontWeight.bold, fontSize: 16),
        ),
      ),
    );
  }
}
