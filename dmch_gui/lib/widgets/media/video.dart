import 'dart:async';

import 'package:dmch_gui/models/entry.dart';
import 'package:dmch_gui/api/dmapi.dart';

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
                      _previewEntries.map((e) => dmapi.getUriFromFilepath(e.filePath)).toList(),
                  headers: dmapi.authHeader,
                )
              : const SizedBox(
                  width: 100,
                  height: 100,
                  child: Center(child: CircularProgressIndicator()),
                ),
        ),
        Text(widget.entry.name),
      ],
    );
  }
}

class VideoPreviews extends StatefulWidget {
  final List<ImageProvider> previews;
  final Map<String, String> headers;

  VideoPreviews({
    Key? key,
    required List<Uri> previewUrls,
    this.headers = const <String, String>{},
  })  : previews = previewUrls.map((e) => NetworkImage(e.toString(), headers: headers)).toList(),
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
    precacheNext();
  }

  @override
  void deactivate() {
    _timer?.cancel();
    super.deactivate();
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  int get _nextPreview {
    if (currentPreview == widget.previews.length - 1) {
      return 0;
    } else {
      return currentPreview + 1;
    }
  }

  Future<void> precacheNext() async {
    if (currentPreview != widget.previews.length - 1) {
      final cacheStatus = await widget.previews[currentPreview + 1]
          .obtainCacheStatus(configuration: const ImageConfiguration());
      if (cacheStatus?.untracked ?? true) {
        await precacheImage(
          widget.previews[currentPreview + 1],
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
      child: Image(
        image: widget.previews[currentPreview],
        loadingBuilder: _loadingBuilder,
        isAntiAlias: true,
      ),
    );
  }

  static Widget _loadingBuilder(
      BuildContext context, Widget child, ImageChunkEvent? loadingProgress) {
    if (loadingProgress != null) {
      if (loadingProgress.expectedTotalBytes != null) {
        return CircularProgressIndicator(
          value: loadingProgress.cumulativeBytesLoaded / loadingProgress.expectedTotalBytes!,
        );
      } else {
        return const CircularProgressIndicator();
      }
    } else {
      return child;
    }
  }
}
