import 'dart:async';

import 'package:dmch_gui/models/entry.dart';
import 'package:dmch_gui/provider/dmapi.dart';

import 'package:extended_image/extended_image.dart';
import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';

import 'package:provider/provider.dart';

class VideoInfoItem extends StatefulWidget {
  final String dirPath;
  final Entry entry;

  const VideoInfoItem({Key? key, required this.dirPath, required this.entry}) : super(key: key);

  @override
  State<VideoInfoItem> createState() => _VideoInfoItemState();
}

class _VideoInfoItemState extends State<VideoInfoItem> {
  List<Entry> _previewEntries = [];

  @override
  void initState() {
    super.initState();

    Future(() async {
      _previewEntries = await Provider.of<DmApiClient>(context, listen: false)
          .getPreviews(widget.entry.filePath)
          .toList();

      setState(() {});
    });
  }

  @override
  Widget build(BuildContext context) {
    final dmapi = Provider.of<DmApiClient>(context, listen: false);

    return Column(
      children: [
        AspectRatio(
          aspectRatio: 16 / 9,
          child: _previewEntries.isNotEmpty
              ? VideoPreviews(
                  previewUrls:
                      _previewEntries.map((e) => dmapi.getUrlFromFilepath(e.filePath)).toList(),
                  headers: dmapi.authHeader,
                )
              : const SizedBox(
                  width: 100,
                  height: 100,
                  child: CircularProgressIndicator(),
                ),
        ),
        Text(widget.entry.name),
      ],
    );
  }
}

class VideoPreviews extends StatefulWidget {
  final List<Image> previews;
  final Map<String, String> headers;

  VideoPreviews({
    Key? key,
    required List<String> previewUrls,
    this.headers = const <String, String>{},
  })  : previews = previewUrls.map((e) => Image.network(e, headers: headers)).toList(),
        super(key: key);

  @override
  State<VideoPreviews> createState() => _VideoPreviewsState();
}

class _VideoPreviewsState extends State<VideoPreviews> {
  int currentPreview = 0;
  Timer? _timer;

  @override
  void initState() {
    super.initState();
    Future.value(precacheNext);
  }

  int get _nextPreview {
    if (currentPreview == widget.previews.length - 1) {
      return 0;
    } else {
      return currentPreview + 1;
    }
  }

  void precacheNext() async {
    if (currentPreview != widget.previews.length - 1) {
      print("precaching: " + widget.previews[currentPreview + 1].image.toString());
      final cached = await widget.previews[currentPreview + 1].image
          .obtainCacheStatus(configuration: const ImageConfiguration());
      if (cached?.untracked ?? true) {
        await precacheImage(
          widget.previews[currentPreview + 1].image,
          context,
        );
      }
    }
  }

  void startRotate(PointerEnterEvent event) {
    _timer?.cancel();
    _timer = Timer.periodic(const Duration(seconds: 1), (_) {
      setState(() {
        currentPreview = _nextPreview;
      });
      precacheNext();
    });
  }

  void stopRotate(PointerExitEvent event) {
    _timer?.cancel();
    setState(() {
      currentPreview = 0;
    });
  }

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      onEnter: startRotate,
      onExit: stopRotate,
      child: Stack(
        fit: StackFit.expand,
        children: [
          widget.previews[currentPreview],
          widget.previews[_nextPreview],
        ],
      ),
    );
  }
}
