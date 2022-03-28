import 'dart:convert';

import 'package:dmch_gui/models/entry.dart';
import 'package:dmch_gui/provider/dmapi.dart';
import 'package:dmch_gui/widgets/media/video.dart';
import 'package:flutter/material.dart';
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
  String _path = basePath;
  final _pathController = TextEditingController(text: basePath);

  List<Entry> _entries = [];

  Future<void> changePath([String? newPath]) async {
    if (newPath != null) {
      _path = newPath;
    }
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
      changePath(_pathController.text);
    });

    Future(changePath);
  }

  @override
  Widget build(BuildContext context) {
    // return Column(
    //   children: [
    //     Container(
    //       height: 100,
    //       padding: const EdgeInsets.all(10),
    //       child: Text(_path),
    //     ),
    //   ],
    // );

    return Column(
      children: [
        SizedBox(
          height: 80,
          child: TextField(
            controller: _pathController,
          ),
        ),
        Expanded(
          child: GridView.extent(
            maxCrossAxisExtent: 200,
            children: _entries
                .where((element) => !element.isDir)
                .map((e) => VideoInfoItem(entry: e, dirPath: "/"))
                .toList(),
          ),
        ),
      ],
    );
  }
}
