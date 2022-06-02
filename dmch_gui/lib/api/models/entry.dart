import 'dart:convert';

import 'package:dmch_gui/api/models/mime.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:path/path.dart' as path_utils;

class Entry {
  final String path;
  final int size;
  final DateTime modType;
  final MimeType mimeType;

  Entry({
    required this.path,
    required this.size,
    required this.modType,
    required this.mimeType,
  });

  String get name => path_utils.basename(path);
  bool get isDir => mimeType == const MimeType.directory();

  Map<String, dynamic> toMap() {
    return <String, dynamic>{
      'path': path,
      'size': size,
      'modType': modType.toUtc().toIso8601String(),
      'mimeType': mimeType.toString(),
    };
  }

  factory Entry.fromMap(Map<String, dynamic> map) {
    return Entry(
      path: map['path'] as String,
      size: map['size'] as int,
      modType: DateTime.parse(map['modType'] as String),
      mimeType: MimeType(map['mimeType'] as String),
    );
  }

  String toJson() => json.encode(toMap());

  factory Entry.fromJson(String source) =>
      Entry.fromMap(json.decode(source) as Map<String, dynamic>);
}
