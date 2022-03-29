import 'dart:convert';

import 'package:dmch_gui/models/entry.dart';
import 'package:dmch_gui/provider/dmapi.dart';
import 'package:dmch_gui/widgets/media/folder.dart';
import 'package:dmch_gui/widgets/media/video.dart';
import 'package:flutter/material.dart';
import 'package:path/path.dart' as path_utils;
import 'package:flutter/widgets.dart';
import 'package:http/http.dart' as http;
import 'package:provider/provider.dart';

import '../../models/media.dart';

const basePath = "/";

class MediaGrid extends StatefulWidget {
  const MediaGrid({Key? key}) : super(key: key);

  @override
  State<MediaGrid> createState() => _MediaGridState();
}

class _MediaGridState extends State<MediaGrid> {
  String get _path => _pathController.text;
  set _path(String newpath) => _pathController.text = newpath;
  final _pathController = TextEditingController(text: basePath);

  List<Entry> _entries = [];

  Future<void> dirUp() async {
    if (_path != "/") {
      await changePath(path_utils.dirname(_path));
    }
  }

  Future<void> dirDown(String dir) async => await changePath(path_utils.dirname(_path));

  Future<void> changePath(String newPath) async {
    _path = newPath;

    try {
      _entries = await Provider.of<DmApiClient>(context, listen: false).getEntries(_path);
      setState(() {});
    } catch (e) {
      print("exception for path: " + _path + ": " + e.toString());
    }
  }

  @override
  void initState() {
    super.initState();

    _pathController.addListener(() {
      changePath(_path);
    });

    Future((() async => await changePath(basePath)));
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        SizedBox(
          height: 80,
          child: Row(children: [
            IconButton(
              onPressed: () => changePath(path_utils.dirname(_path)),
              icon: const Icon(Icons.arrow_back),
            ),
            Expanded(
              child: TextField(
                controller: _pathController,
              ),
            ),
          ]),
        ),
        Expanded(
          child: GridView.extent(
            maxCrossAxisExtent: 200,
            children: <Widget>[
              ..._entries
                  .where((element) => element.isDir)
                  .map((e) => GestureDetector(
                        onDoubleTap: () => changePath(path_utils.joinAll([_path, e.name])),
                        child: FolderItem(entry: e),
                      ))
                  .toList(),
              ..._entries
                  .where((element) => !element.isDir)
                  .map((e) => VideoInfoItem(entry: e, dirPath: "/"))
                  .toList(),
            ],
          ),
        ),
      ],
    );
  }
}
