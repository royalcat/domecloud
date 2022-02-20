import 'package:dmch_gui/widgets/file/theme.dart';
import 'package:file/file.dart';
import 'package:file/local.dart';
import 'package:flutter/material.dart';

import 'entity.dart';

class FSGrid extends StatefulWidget {
  final FileSystem fs;

  const FSGrid({Key? key, required this.fs}) : super(key: key);

  @override
  State<FSGrid> createState() => _FSGridState();
}

class _FSGridState extends State<FSGrid> {
  Iterable<FileSystemEntity> entites = <FileSystemEntity>[];

  @override
  void initState() {
    super.initState();

    widget.fs.directory("/home/royalcat/").list().toList().then((list) => setState(() {
          entites = list.where((element) => !element.basename.startsWith("."));
        }));
  }

  @override
  Widget build(BuildContext context) {
    return FSGridTheme(
      theme: FSGridThemeData.normal(),
      child: GridView.extent(
        maxCrossAxisExtent: 200,
        children: entites.map((e) => FSEntityItem(entity: e)).toList(),
      ),
    );
  }
}
