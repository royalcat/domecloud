import 'package:flutter/widgets.dart';

import '../../models/media.dart';

class VideoInfoItem extends StatefulWidget {
  final VideoInfo info;

  VideoInfoItem({Key? key, required this.info}) : super(key: key);

  @override
  State<VideoInfoItem> createState() => _VideoInfoItemState();
}

class _VideoInfoItemState extends State<VideoInfoItem> {
  @override
  Widget build(BuildContext context) {
    return Row(
      children: [],
    );
  }
}
