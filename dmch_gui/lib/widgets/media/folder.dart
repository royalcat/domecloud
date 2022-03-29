import 'package:dmch_gui/models/entry.dart';
import 'package:flutter/material.dart';

class FolderItem extends StatelessWidget {
  final Entry entry;

  const FolderItem({Key? key, required this.entry}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        const Icon(
          Icons.folder,
          size: 100,
        ),
        Text(entry.name),
      ],
    );
  }
}
