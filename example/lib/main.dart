import 'package:flutter/material.dart';

import 'package:serial_port_flutter_example/app.dart';
import 'package:serial_port_flutter_example/stores/app.dart';

import 'package:provider/provider.dart';

void main() {
  runApp(
    ChangeNotifierProvider(
      builder: (context) => AppModel(),
      child: MyApp(),
    ),
  );
}