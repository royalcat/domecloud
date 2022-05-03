import 'package:dmch_gui/views/fs_view.dart';
import 'package:dmch_gui/views/login.dart';
import 'package:dmch_gui/views/main_view.dart';
import 'package:dmch_gui/views/media_view.dart';
import 'package:dmch_gui/widgets/media/grid.dart';
import 'package:flutter/material.dart';
import 'package:window_manager/window_manager.dart';

// class MyImageCache extends ImageCache {
//   @override
//   void clear() {
//     print('Clearing cache!');
//     super.clear();
//   }
// }

// class MyWidgetsBinding extends WidgetsFlutterBinding {
//   @override
//   ImageCache createImageCache() => MyImageCache();
// }

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  imageCache.maximumSize = 10000;
  imageCache.maximumSizeBytes = 1024 * 1024 * 1024 * 1024; // 1 GB

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
      title: 'DomeCoud',
      theme: ThemeData.light().copyWith(primaryColor: Colors.greenAccent),
      darkTheme: ThemeData.dark().copyWith(primaryColor: Colors.greenAccent),
      themeMode: ThemeMode.dark,
      home: const MainView(),
    );
  }
}
