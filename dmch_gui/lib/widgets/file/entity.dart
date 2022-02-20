import 'dart:io';

import 'package:dmch_gui/widgets/file/theme.dart';
import 'package:file/file.dart';
import 'package:flutter/material.dart';

class FSEntityItem extends StatelessWidget {
  final FileSystemEntity entity;

  const FSEntityItem({Key? key, required this.entity}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final theme = FSGridTheme.of(context);
    final IconData icon;
    if (entity is Directory) {
      icon = theme.folderIcon;
    } else if (entity is File) {
      icon = theme.fileIcon;
    } else {
      icon = Icons.question_mark;
    }

    // final icon = <Type, IconData>{
    //       Directory: theme.folderIcon,
    //       File: theme.fileIcon,
    //     }[entity.runtimeType] ??
    //     Icons.question_mark;

    return SizedBox.square(
      dimension: 100 * theme.scale,
      child: Padding(
        padding: const EdgeInsets.all(8.0),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          mainAxisSize: MainAxisSize.max,
          children: <Widget>[
            Icon(
              icon,
              size: 90 * theme.scale,
            ),
            Text(
              entity.basename,
            ),
          ],
        ),
      ),
    );
  }
}
