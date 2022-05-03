import 'package:dmch_gui/widgets/media/grid.dart';
import 'package:flutter/material.dart';

class MediaView extends StatefulWidget {
  const MediaView({Key? key}) : super(key: key);

  @override
  State<MediaView> createState() => _MediaViewState();
}

class _MediaViewState extends State<MediaView> {
  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      body: MediaGrid(),
    );
  }
}
