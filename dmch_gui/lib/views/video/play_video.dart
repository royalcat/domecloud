import 'dart:convert';
import 'dart:io';
import 'dart:ui' as ui;

import 'package:dart_vlc/dart_vlc.dart';
import 'package:dmch_gui/views/video/my_video.dart';
import 'package:flutter/material.dart';

void playVideo(BuildContext context, Uri uri, Map<String, String> headers) async {
  Navigator.push(
    context,
    MaterialPageRoute<void>(
      builder: (context) => VLCPlayer(Media.network(uri, extras: headers)),
    ),
  );
}

class VLCPlayer extends StatefulWidget {
  final Media media;

  const VLCPlayer(this.media);

  @override
  VLCPlayerState createState() => VLCPlayerState();
}

class VLCPlayerState extends State<VLCPlayer> {
  Player player = Player(
    id: 0,
    videoDimensions: const VideoDimensions(640, 360),
    registerTexture: !Platform.isWindows,
  );
  // MediaType mediaType = MediaType.file;
  // CurrentState current = CurrentState();
  // PositionState position = PositionState();
  // PlaybackState playback = PlaybackState();
  // GeneralState general = GeneralState();
  // VideoDimensions videoDimensions = const VideoDimensions(0, 0);
  // List<Device> devices = <Device>[];
  // TextEditingController controller = TextEditingController();
  // TextEditingController metasController = TextEditingController();
  // double bufferingProgress = 0.0;
  // Media? metasMedia;

  @override
  void initState() {
    super.initState();
    // if (mounted) {
    //   player.currentStream.listen((current) {
    //     setState(() => this.current = current);
    //   });
    //   player.positionStream.listen((position) {
    //     setState(() => this.position = position);
    //   });
    //   player.playbackStream.listen((playback) {
    //     setState(() => this.playback = playback);
    //   });
    //   player.generalStream.listen((general) {
    //     setState(() => this.general = general);
    //   });
    //   player.videoDimensionsStream.listen((videoDimensions) {
    //     setState(() => this.videoDimensions = videoDimensions);
    //   });
    //   player.bufferingProgressStream.listen(
    //     (bufferingProgress) {
    //       setState(() => this.bufferingProgress = bufferingProgress);
    //     },
    //   );
    //   player.errorStream.listen((event) {
    //     debugPrint('libvlc error: $event');
    //   });
    //   devices = Devices.all;
    //   final Equalizer equalizer = Equalizer.createMode(EqualizerMode.live);
    //   equalizer.setPreAmp(10.0);
    //   equalizer.setBandAmp(31.25, 10.0);
    //   player.setEqualizer(equalizer);
    // }

    player.open(widget.media);
  }

  @override
  void dispose() {
    player.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: MyVideo(
        player: player,
      ),
    );
  }
}
