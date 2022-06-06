import 'package:dmch_gui/api/models/entry.dart';
import 'package:flutter/material.dart';

class FolderList extends StatelessWidget {
  final List<Entry> entries;
  final void Function(Entry entry) onOpen;

  const FolderList({Key? key, required this.entries, required this.onOpen}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ListView(
      children: entries
          .map<Widget>(
            (entry) => ListTile(
              leading: const Icon(Icons.folder),
              title: Text(entry.name),
              onTap: () => onOpen(entry),
            ),
          )
          .toList(),
    );
  }
}

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
