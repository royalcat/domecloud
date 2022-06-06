// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:convert';

class MimeType {
  final String string;
  const MimeType(this.string);
  const MimeType.directory() : string = "inode/directory";
  const MimeType.json() : string = "application/json";

  @override
  String toString() => string;

  MimeType.fromString(this.string);

  @override
  int get hashCode => string.hashCode;
  @override
  bool operator ==(Object other) =>
      other is MimeType && other.runtimeType == runtimeType && other.string == string;
}

enum MediaType {
  none("none"),
  video("video"),
  photo("photo");

  final String string;
  const MediaType(this.string);

  @override
  String toString() => string;
}
