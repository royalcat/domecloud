import 'dart:convert';

import 'package:collection/collection.dart';
import 'package:dmch_gui/api/models/mime.dart';
import 'package:path/path.dart' as path_utils;

class Entry {
  final String name;
  final bool isListable;
  final MimeType mimeType;

  final String dir;

  late final isDir = mimeType == const MimeType.directory();
  late final path = path_utils.join(dir, name);

  Entry({required this.isListable, required this.mimeType, required this.dir, required this.name});

  @override
  String toString() {
    return 'Entry(name: $name, isListable: $isListable, mimeType: $mimeType)';
  }

  factory Entry.fromMap(String dir, Map<String, dynamic> data) => Entry(
        isListable: data['isListable'] as bool,
        mimeType: MimeType(data['mimeType'] as String),
        name: data['name'] as String,
        dir: dir,
      );

  Map<String, dynamic> toMap() => {
        'isListable': isListable,
        'mimeType': mimeType,
        'name': name,
        'dir': dir,
      };

  factory Entry.fromJson(String dir, String data) {
    return Entry.fromMap(dir, json.decode(data) as Map<String, dynamic>);
  }

  String toJson() => json.encode(toMap());

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    if (other is! Entry) return false;
    final mapEquals = const DeepCollectionEquality().equals;
    return mapEquals(other.toMap(), toMap());
  }

  @override
  int get hashCode => isListable.hashCode ^ mimeType.hashCode ^ name.hashCode;
}
