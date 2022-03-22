import 'dart:convert';

import 'package:dmch_gui/models/entry.dart';
import 'package:dmch_gui/provider/dmapi.dart';
import 'package:dmch_gui/widgets/media/video.dart';
import 'package:flutter/widgets.dart';
import 'package:http/http.dart' as http;
import 'package:provider/provider.dart';

import '../../models/media.dart';

const baseUrl = 'http://localhost:5050/file';

class MediaGrid extends StatefulWidget {
  const MediaGrid({Key? key}) : super(key: key);

  @override
  State<MediaGrid> createState() => _MediaGridState();
}

class _MediaGridState extends State<MediaGrid> {
  String _path = "/";
  List<Entry> _entries = [];

  Future<void> changePath([String? newPath]) async {
    if (newPath != null) {
      _path = newPath;
    }
    _entries = await Provider.of<DmApiClient>(context, listen: false).getEntries(_path);
    setState(() {});
  }

  @override
  void initState() {
    super.initState();

    Future(changePath);
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
