import 'package:dmch_gui/scroll.dart';
import 'package:file/file.dart';
import 'package:flutter/material.dart';

import 'package:dmch_gui/widgets/file/theme.dart';

import 'entity.dart';

class FileGrid extends StatelessWidget {
  final Iterable<FileSystemEntity> entites;
  final scrollCtrl = ScrollController();

  FileGrid({
    Key? key,
    required this.entites,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return FSGridTheme(
      theme: FSGridThemeData.normal(),
      child: SmoothScroll(
        controller: scrollCtrl,
        child: GridView.extent(
          maxCrossAxisExtent: 200,
          controller: scrollCtrl,
          children: entites.map((e) => FSEntityItem(entity: e)).toList(),
        ),
      ),
    );
  }
}
