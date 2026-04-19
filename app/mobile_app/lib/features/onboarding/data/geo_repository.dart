import 'dart:convert';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

final geoDataProvider = FutureProvider<Map<String, dynamic>>((ref) async {
  final jsonString = await rootBundle.loadString('assets/data/ugandaGeoData.json');
  return json.decode(jsonString) as Map<String, dynamic>;
});
