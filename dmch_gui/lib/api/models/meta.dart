import 'package:dmch_gui/utils/duration.dart';

class EntryMeta {
  String path;
  int size;
  DateTime modTime;
  String mimeType;
  MediaInfo? mediaInfo;

  EntryMeta({
    required this.path,
    required this.size,
    required this.modTime,
    required this.mimeType,
    this.mediaInfo,
  });

  EntryMeta.fromMap(Map<String, dynamic> json)
      : path = json['path'] as String,
        size = json['size'] as int,
        modTime = DateTime.parse(json['modTime'] as String),
        mimeType = json['MimeType'] as String,
        mediaInfo = json['mediaInfo'] != null
            ? MediaInfo.fromJson(json['mediaInfo'] as Map<String, dynamic>)
            : null;

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['path'] = path;
    data['size'] = size;
    data['modTime'] = modTime;
    data['MimeType'] = mimeType;
    if (mediaInfo != null) {
      data['mediaInfo'] = mediaInfo?.toJson();
    }
    return data;
  }
}

class MediaInfo {
  String? mediaType;
  VideoInfo? videoInfo;

  MediaInfo({this.mediaType, this.videoInfo});

  MediaInfo.fromJson(Map<String, dynamic> json) {
    mediaType = json['mediaType'] as String;
    videoInfo = json['videoInfo'] != null
        ? VideoInfo.fromJson(json['videoInfo'] as Map<String, dynamic>)
        : null;
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['mediaType'] = mediaType;
    if (videoInfo != null) {
      data['videoInfo'] = videoInfo?.toJson();
    }
    return data;
  }
}

class VideoInfo {
  Duration? duration;
  Resolution? resolution;

  VideoInfo({this.duration, this.resolution});

  VideoInfo.fromJson(Map<String, dynamic> json) {
    duration = DurationBatteries.fromNanoseconds(json['duration'] as int);
    resolution = json['resolution'] != null
        ? Resolution.fromJson(json['resolution'] as Map<String, dynamic>)
        : null;
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['duration'] = duration;
    if (resolution != null) {
      data['resolution'] = resolution?.toJson();
    }
    return data;
  }
}

class Resolution {
  int height;
  int width;

  Resolution({required this.height, required this.width});

  Resolution.fromJson(Map<String, dynamic> json)
      : height = json['height'] as int,
        width = json['width'] as int;

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['height'] = height;
    data['width'] = width;
    return data;
  }

  @override
  String toString() => '${width}x$height';
}
