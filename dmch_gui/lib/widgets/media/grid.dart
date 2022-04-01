import 'dart:async';
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
  final _pathController = TextEditingController(text: basePath);
  var _suggestions = <String>[];

  List<Entry> _entries = [];

  Future<void> dirUp() async {
    if (_pathController.text != "/") {
      _pathController.text = path_utils.dirname(_pathController.text);
      await updateViewForPath();
    }
  }

  Future<void> dirDown(String dir) async {
    _pathController.text = path_utils.joinAll([_pathController.text, dir]);
    await updateViewForPath();
  }

  Future<void> updateViewForPath() async {
    try {
      _entries = await Provider.of<DmApiClient>(context, listen: false)
          .getEntries(_pathController.text)
          .toList();
      setState(() {});
    } catch (e) {
      debugPrint("exception for path: " + _pathController.text + ": " + e.toString());
      await dirUp();
    }
  }

  Future<List<String>> suggestions(String prefix) async {
    try {
      if (prefix.endsWith("/")) {
        final entries = await Provider.of<DmApiClient>(context, listen: false).getEntries(prefix);
        return entries.where((e) => !e.isDir).map<String>((e) => e.name).toList();
      } else {
        final dir = path_utils.dirname(prefix);
        final entries = await Provider.of<DmApiClient>(context, listen: false).getEntries(dir);
        final query = path_utils.basename(prefix);
        return entries
            .where((e) => !e.isDir)
            .where((element) => element.name.startsWith(query))
            .map<String>((e) => path_utils.joinAll([dir, e.name]))
            .toList();
      }
    } catch (e) {
      print("exception for path: " + _pathController.text + ": " + e.toString());
    }

    return <String>[];
  }

  @override
  void initState() {
    super.initState();

    Future((() async => await updateViewForPath()));
    _pathController.addListener(() {
      suggestions(_pathController.text).then((value) => setState(() {
            _suggestions = value;
          }));
    });
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        SizedBox(
          height: 80,
          child: Row(children: [
            IconButton(
              onPressed: () => updateViewForPath(),
              icon: const Icon(Icons.arrow_back),
            ),
            Expanded(
              child: Autocomplete<String>(
                  optionsBuilder: (textEditingValue) => suggestions(textEditingValue.text)),
            ),
            //  TextField(
            //   onEditingComplete: () => updateViewForPath(),
            //   autofillHints: _suggestions,
            //   controller: _pathController,
            // ),
          ]),
        ),
        Expanded(
          child: GridView.extent(
            maxCrossAxisExtent: 200,
            children: <Widget>[
              ..._entries
                  .where((element) => element.isDir)
                  .map((e) => GestureDetector(
                        onDoubleTap: () => dirDown(e.name),
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
