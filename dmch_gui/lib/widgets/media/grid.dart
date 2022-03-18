import 'dart:convert';

import 'package:dmch_gui/models/entry.dart';
import 'package:dmch_gui/widgets/media/video.dart';
import 'package:flutter/widgets.dart';
import 'package:http/http.dart' as http;

import '../../models/media.dart';

const baseUrl = 'http://localhost:5050/file';

class MediaGrid extends StatefulWidget {
  const MediaGrid({Key? key}) : super(key: key);

  @override
  State<MediaGrid> createState() => _MediaGridState();
}

class _MediaGridState extends State<MediaGrid> {
  List<Entry> _entries = [];

  @override
  void initState() {
    super.initState();
    Future(() async {
      final resp = await http.get(Uri.parse(baseUrl));
      setState(() {
        _entries = (json.decode(resp.body) as List<dynamic>).map((e) => Entry.fromMap(e)).toList();
      });
    });
  }

  @override
  Widget build(BuildContext context) {
    return GridView.extent(
      maxCrossAxisExtent: 200,
      children: _entries
          .where((element) => !element.isDir)
          .map((e) => VideoInfoItem(entry: e, dirPath: "/"))
          .toList(),
    );
  }
}
