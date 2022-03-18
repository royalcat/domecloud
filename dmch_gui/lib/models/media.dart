import 'dart:convert';

class VideoInfo {
  final String path;
  final int size;
  final DateTime modificationTime;
  final Duration duration;
  final Resolution resolution;

  VideoInfo(
    this.path,
    this.size,
    this.modificationTime,
    this.duration,
    this.resolution,
  );

  String basename() {
    return path.split("/").last;
  }

  Map<String, dynamic> toMap() {
    return {
      'path': path,
      'size': size,
      'modTime': modificationTime.millisecondsSinceEpoch,
      'duration': duration.inMicroseconds * 1000,
      'resolution': resolution.toMap(),
    };
  }

  factory VideoInfo.fromMap(Map<String, dynamic> map) {
    return VideoInfo(
      map['path'] ?? '',
      map['size']?.toInt() ?? 0,
      DateTime.fromMillisecondsSinceEpoch(map['modTime']),
      Duration(microseconds: ((map['duration'] / 1000) as double).floor()),
      Resolution.fromMap(map['resolution']),
    );
  }

  String toJson() => json.encode(toMap());

  factory VideoInfo.fromJson(String source) => VideoInfo.fromMap(json.decode(source));
}

class Resolution {
  final int height;
  final int width;

  Resolution(this.height, this.width);

  Map<String, dynamic> toMap() {
    return {
      'height': height,
      'width': width,
    };
  }

  factory Resolution.fromMap(Map<String, dynamic> map) {
    return Resolution(
      map['height']?.toInt() ?? 0,
      map['width']?.toInt() ?? 0,
    );
  }

  String toJson() => json.encode(toMap());

  factory Resolution.fromJson(String source) => Resolution.fromMap(json.decode(source));
}
