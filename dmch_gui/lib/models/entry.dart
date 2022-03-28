import 'dart:convert';
import 'package:path/path.dart' as path_utils;

enum MediaTypes { none, video, photo }

class Entry {
  final String name;
  final bool isDir;
  final String mimeType;
  final String filePath;

  Entry({required this.name, required this.isDir, required this.mimeType, required this.filePath});

  Map<String, dynamic> toMap() {
    return {
      'name': name,
      'isDir': isDir,
      'mimeType': mimeType.toString(),
    };
  }

  factory Entry.fromMap(Map<String, dynamic> map, String dir) {
    return Entry(
      name: map['name'] ?? '',
      isDir: map['isDir'] ?? false,
      mimeType: map['mimeType'],
      filePath: path_utils.posix.joinAll([dir, map['name'] ?? '']),
    );
  }
}
