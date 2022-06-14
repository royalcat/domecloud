import 'dart:async';

import 'package:dmch_gui/api/dmapi.dart';
import 'package:dmch_gui/api/models/entry.dart';
import 'package:dmch_gui/views/video/play_video.dart';
import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';

import 'package:provider/provider.dart';

class VideoInfoItem extends StatefulWidget {
  final Entry entry;
  final void Function(Entry)? onOpenDetails;

  const VideoInfoItem({Key? key, this.onOpenDetails, required this.entry}) : super(key: key);

  @override
  State<VideoInfoItem> createState() => _VideoInfoItemState();
}

class _VideoInfoItemState extends State<VideoInfoItem> {
  @override
  Widget build(BuildContext context) {
    final dmapi = Provider.of<DmApiClient>(context, listen: false);

    return Column(
      children: [
        GestureDetector(
          onDoubleTap: () async {
            final uri = await dmapi.getPlayUrl(widget.entry.path);
            playVideo(context, uri);
          },
          onTap: () => widget.onOpenDetails?.call(widget.entry),
          child: AspectRatio(
            aspectRatio: 16 / 9,
            child: FutureBuilder<List<Entry>>(
              future: dmapi.getPreviews(widget.entry.path).toList(),
              builder: (BuildContext context, AsyncSnapshot<List<Entry>> snapshot) =>
                  snapshot.hasData && snapshot.data != null
                      ? VideoPreview(
                          previewUrls: snapshot.data!.map((e) => dmapi.getFileUri(e.path)).toList(),
                          headers: dmapi.authHeader,
                        )
                      : const SizedBox(
                          width: 100,
                          height: 100,
                          child: Center(child: CircularProgressIndicator()),
                        ),
            ),
          ),
        ),
        Text(widget.entry.name),
      ],
    );
  }
}

class VideoPreview extends StatefulWidget {
  final List<ImageProvider> previews;
  final Map<String, String> headers;

  VideoPreview({
    Key? key,
    required List<Uri> previewUrls,
    this.headers = const <String, String>{},
  })  : previews = previewUrls.map((e) => NetworkImage(e.toString(), headers: headers)).toList(),
        super(key: key);

  @override
  State<VideoPreview> createState() => _VideoPreviewState();
}

class _VideoPreviewState extends State<VideoPreview> {
  int currentPreview = 0;
  Timer? _timer;

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
    return currentPreview == widget.previews.length - 1 ? 0 : currentPreview + 1;
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    precacheNext();
  }

  void precacheNext() {
    if (currentPreview != widget.previews.length - 1) {
      precacheImage(
        widget.previews[currentPreview + 1],
        context,
      );
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
    BuildContext context,
    Widget child,
    ImageChunkEvent? loadingProgress,
  ) {
    return loadingProgress != null
        ? loadingProgress.expectedTotalBytes != null
            ? CircularProgressIndicator(
                value: loadingProgress.cumulativeBytesLoaded / loadingProgress.expectedTotalBytes!,
              )
            : const CircularProgressIndicator()
        : child;
  }
}
