import 'package:flutter/foundation.dart' show debugDefaultTargetPlatformOverride;
import 'package:flutter/material.dart';

import 'package:serial_port_flutter_example/stores/app.dart';
import 'package:serial_port_flutter_example/app.dart';

import 'package:provider/provider.dart';

void main() {
  debugDefaultTargetPlatformOverride = TargetPlatform.fuchsia;
  runApp(
    ChangeNotifierProvider(
      builder: (context) => AppModel(),
      child: MyApp(),
    ),
  );
}
