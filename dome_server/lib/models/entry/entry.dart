// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:io';

import 'package:json_annotation/json_annotation.dart';

part 'entry.g.dart';

@JsonSerializable()
class Entry {
  final String path;
  final int size;
  final DateTime modTime;

  Entry({
    required this.path,
    required this.size,
    required this.modTime,
  });

  factory Entry.fromJson(Map<String, dynamic> json) => _$EntryFromJson(json);

  Map<String, dynamic> toJson() => _$EntryToJson(this);
}

@JsonSerializable()
class MediaEntry extends Entry {
  MediaEntry({
    required super.path,
    required super.size,
    required super.modTime,
  });

  factory MediaEntry.fromJson(Map<String, dynamic> json) => _$MediaEntryFromJson(json);

  @override
  Map<String, dynamic> toJson() => _$MediaEntryToJson(this);
}

@JsonSerializable()
class MediaInfo {}

@JsonSerializable()
class Resolution {
  final int height;
  final int width;
  Resolution({
    required this.height,
    required this.width,
  });

  factory Resolution.fromJson(Map<String, dynamic> json) => _$ResolutionFromJson(json);

  Map<String, dynamic> toJson() => _$ResolutionToJson(this);
}
