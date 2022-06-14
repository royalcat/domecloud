import 'dart:convert';

import 'package:dmch_gui/api/models/media/media.dart';
import 'package:dmch_gui/api/models/media/resolution.dart';
import 'package:dmch_gui/utils/duration.dart';

class VideoInfo extends MediaInfo {
  final Duration duration;
  final Resolution resolution;
  VideoInfo({
    required this.duration,
    required this.resolution,
  });

  VideoInfo copyWith({
    Duration? duration,
    Resolution? resolution,
  }) {
    return VideoInfo(
      duration: duration ?? this.duration,
      resolution: resolution ?? this.resolution,
    );
  }

  Map<String, dynamic> toMap() {
    return <String, dynamic>{
      'duration': duration.inNanoseconds,
      'resolution': resolution.toMap(),
    };
  }

  factory VideoInfo.fromMap(Map<String, dynamic> map) {
    return VideoInfo(
      duration: DurationBatteries.fromNanoseconds(map['duration'] as int),
      resolution: Resolution.fromMap(map['resolution'] as Map<String, dynamic>),
    );
  }

  String toJson() => json.encode(toMap());

  factory VideoInfo.fromJson(String source) =>
      VideoInfo.fromMap(json.decode(source) as Map<String, dynamic>);

  @override
  String toString() => 'VideoInfo(duration: $duration, resolution: $resolution)';

  @override
  int get hashCode => duration.hashCode ^ resolution.hashCode;

  @override
  bool operator ==(Object other) =>
      other is VideoInfo &&
      other.runtimeType == runtimeType &&
      other.duration == duration &&
      other.resolution == resolution;
}
