// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'entry.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Entry _$EntryFromJson(Map<String, dynamic> json) => Entry(
      path: json['path'] as String,
      size: json['size'] as int,
      modTime: DateTime.parse(json['modTime'] as String),
    );

Map<String, dynamic> _$EntryToJson(Entry instance) => <String, dynamic>{
      'path': instance.path,
      'size': instance.size,
      'modTime': instance.modTime.toIso8601String(),
    };

MediaEntry _$MediaEntryFromJson(Map<String, dynamic> json) => MediaEntry(
      path: json['path'] as String,
      size: json['size'] as int,
      modTime: DateTime.parse(json['modTime'] as String),
    );

Map<String, dynamic> _$MediaEntryToJson(MediaEntry instance) =>
    <String, dynamic>{
      'path': instance.path,
      'size': instance.size,
      'modTime': instance.modTime.toIso8601String(),
    };

MediaInfo _$MediaInfoFromJson(Map<String, dynamic> json) => MediaInfo();

Map<String, dynamic> _$MediaInfoToJson(MediaInfo instance) =>
    <String, dynamic>{};

Resolution _$ResolutionFromJson(Map<String, dynamic> json) => Resolution(
      height: json['height'] as int,
      width: json['width'] as int,
    );

Map<String, dynamic> _$ResolutionToJson(Resolution instance) =>
    <String, dynamic>{
      'height': instance.height,
      'width': instance.width,
    };
