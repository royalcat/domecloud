import 'package:dmch_gui/views/fs_view.dart';
import 'package:dmch_gui/views/media_view.dart';
import 'package:dmch_gui/widgets/media/grid.dart';
import 'package:flutter/material.dart';
import 'package:window_manager/window_manager.dart';

void main() {
  // Future(() async {
  //   await windowManager.ensureInitialized();
  //   await windowManager.setTitleBarStyle(TitleBarStyle.hidden);
  // });
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData.light().copyWith(primaryColor: Colors.greenAccent),
      darkTheme: ThemeData.dark().copyWith(primaryColor: Colors.greenAccent),
      themeMode: ThemeMode.dark,
      home: Scaffold(
        appBar: AppBar(
          title: const Text('Flutter Demo'),
        ),
        body: const MediaView(),
      ),
    );
  }
}
