import 'package:dmch_gui/api/dmapi.dart';
import 'package:flutter/material.dart';
import 'package:path/path.dart';
import 'package:provider/provider.dart';
import 'package:video_player/video_player.dart';
import 'package:url_launcher/url_launcher.dart';

Future<void> playVideo(BuildContext context, Uri uri) async {
  await launchUrl(uri);
  // return Navigator.push<void>(
  //   context,
  //   MaterialPageRoute<void>(
  //     builder: (context) => MyVideoPlayer(
  //       uri: uri,
  //     ),
  //   ),
  // );
}

class MyVideoPlayer extends StatefulWidget {
  final Uri uri;
  final VideoPlayerController _controller;

  MyVideoPlayer({Key? key, required this.uri})
      : _controller = VideoPlayerController.network(uri.toString()),
        super(key: key);

  @override
  State<MyVideoPlayer> createState() => _MyVideoPlayerState();
}

class _MyVideoPlayerState extends State<MyVideoPlayer> {
  @override
  Widget build(BuildContext context) {
    return VideoPlayer(widget._controller);
  }
}
