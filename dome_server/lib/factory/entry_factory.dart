import 'dart:io';

import '../models/entry/entry.dart';

class EntryFactory {
  Future<Entry> fromFile(File file) async {
    final stat = await file.stat();
    return Entry(
      path: file.path,
      size: stat.size,
      modTime: stat.modified,
    );
  }
}
