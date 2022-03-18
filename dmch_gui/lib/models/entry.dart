import 'dart:convert';

class Entry {
  final String name;
  final bool isDir;

  Entry(this.name, this.isDir);

  Map<String, dynamic> toMap() {
    return {
      'name': name,
      'isDir': isDir,
    };
  }

  factory Entry.fromMap(Map<String, dynamic> map) {
    return Entry(
      map['name'] ?? '',
      map['isDir'] ?? false,
    );
  }

  String toJson() => json.encode(toMap());

  factory Entry.fromJson(String source) => Entry.fromMap(json.decode(source));
}
