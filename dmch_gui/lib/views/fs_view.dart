import 'package:dmch_gui/widgets/file/file_gird.dart';
import 'package:file/file.dart';
import 'package:file/local.dart';
import 'package:flutter/material.dart';

class FSView extends StatefulWidget {
  const FSView({Key? key}) : super(key: key);

  @override
  State<FSView> createState() => _FSViewState();
}

class _FSViewState extends State<FSView> {
  final fs = const LocalFileSystem()..currentDirectory = "/home/royalcat";

  var entities = <FileSystemEntity>[];

  @override
  void initState() {
    super.initState();
    Future(() async {
      entities = await fs.currentDirectory
          .list()
          .where((element) => !element.basename.startsWith("."))
          .toList();

      setState(() {});
    });
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Expanded(
          child: FileGrid(
            entites: entities,
          ),
        ),
      ],
    );
  }
}
