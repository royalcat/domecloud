import 'package:dmch_gui/widgets/media/video.dart';
import 'package:flutter/widgets.dart';

import '../../models/media.dart';

class MediaGrid extends StatefulWidget {
  final Iterable<VideoInfo> infos;

  const MediaGrid({Key? key, required this.infos}) : super(key: key);

  @override
  State<MediaGrid> createState() => _MediaGridState();
}

class _MediaGridState extends State<MediaGrid> {
  @override
  Widget build(BuildContext context) {
    return GridView.extent(
      maxCrossAxisExtent: 200,
      children: widget.infos.map((e) => VideoInfoItem(info: e)).toList(),
    );
  }
}
