import 'dart:async';
import 'dart:convert';

import 'package:dmch_gui/models/entry.dart';
import 'package:dmch_gui/provider/dmapi.dart';
import 'package:dmch_gui/widgets/media/grid.dart';
import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:path/path.dart' as path;

import 'package:http/http.dart' as http;
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
      for (final e in _previewEntries) {
        precacheImage(
          NetworkImage(
            Provider.of<DmApiClient>(context, listen: false).getUrlFromFilepath(e.filePath),
          ),
          context,
        );
      }

      setState(() {});
    });
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        AspectRatio(
          aspectRatio: 16 / 9,
          child: _previewEntries.isNotEmpty
              ? VideoPreviews(
                  previewUrls: _previewEntries
                      .map(
                        (e) => Provider.of<DmApiClient>(context, listen: false)
                            .getUrlFromFilepath(e.filePath),
                      )
                      .toList(),
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
  final List<String> previewUrls;

  const VideoPreviews({Key? key, required this.previewUrls}) : super(key: key);

  @override
  State<VideoPreviews> createState() => _VideoPreviewsState();
}

class _VideoPreviewsState extends State<VideoPreviews> {
  int currentPreview = 0;
  Timer? _timer;

  void startRotate(PointerEnterEvent event) {
    _timer?.cancel();
    _timer = Timer.periodic(const Duration(seconds: 1), (_) {
      setState(() {
        if (currentPreview == widget.previewUrls.length - 1) {
          currentPreview = 0;
        } else {
          currentPreview++;
        }
      });
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
      child: Image.network(
        // TODO сделать кастомный кеш
        widget.previewUrls[currentPreview],
        width: 100,
        height: 100,
      ),
    );
  }
}
