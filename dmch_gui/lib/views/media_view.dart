import 'package:dmch_gui/provider/dmapi.dart';
import 'package:dmch_gui/widgets/media/grid.dart';
import 'package:flutter/widgets.dart';
import 'package:provider/provider.dart';

class MediaView extends StatefulWidget {
  MediaView({Key? key}) : super(key: key);

  @override
  State<MediaView> createState() => _MediaViewState();
}

class _MediaViewState extends State<MediaView> {
  @override
  Widget build(BuildContext context) {
    return Provider<DmApi>(
      create: (context) => DmApi(),
      child: const MediaGrid(),
    );
  }
}
