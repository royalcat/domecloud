import 'dart:convert';

enum MediaTypes { none, video, photo }

class Entry {
  final String name;
  final bool isDir;
  final MediaTypes mediaType;
  final String filePath;

  Entry({required this.name, required this.isDir, required this.mediaType, required this.filePath});

  Map<String, dynamic> toMap() {
    return {
      'name': name,
      'isDir': isDir,
      'mediaType': mediaType.toString(),
    };
  }

  factory Entry.fromMap(Map<String, dynamic> map, String filePath) {
    return Entry(
        name: map['name'] ?? '',
        isDir: map['isDir'] ?? false,
        mediaType: MediaTypes.values.firstWhere((e) => e.toString() == map['mediaType']),
        filePath: filePath);
  }
}
