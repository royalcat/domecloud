import 'dart:async';
import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter/scheduler.dart';
import 'package:flutter/widgets.dart';
import 'package:http/http.dart' as http;
import 'package:path/path.dart' as path_utils;
import 'package:provider/provider.dart';

import 'package:dmch_gui/api/dmapi.dart';
import 'package:dmch_gui/models/entry.dart';
import 'package:dmch_gui/widgets/media/folder.dart';
import 'package:dmch_gui/widgets/media/video.dart';

import '../../models/media.dart';

const basePath = "/";

class MediaGrid extends StatefulWidget {
  const MediaGrid({Key? key}) : super(key: key);

  @override
  State<MediaGrid> createState() => _MediaGridState();
}

class _MediaGridState extends State<MediaGrid> {
  final _pathController = TextEditingController(text: basePath);

  List<Entry> _entries = [];

  Future<void> dirUp() async {
    if (_pathController.text != "/") {
      await updateViewForPath(path_utils.dirname(_pathController.text));
    }
  }

  Future<void> dirDown(String dir) async {
    await updateViewForPath(path_utils.joinAll([_pathController.text, dir]));
  }

  Future<void> updateViewForPath(String path) async {
    try {
      if (path != null) {
        _pathController.text = path_utils.normalize(path);
      }
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
        return [prefix] +
            await Provider.of<DmApiClient>(context, listen: false)
                .getEntries(prefix)
                .where((e) => e.isDir)
                .map<String>((e) => path_utils.joinAll([prefix, e.name]))
                .toList();
      } else {
        final dir = path_utils.dirname(prefix);
        final query = path_utils.basename(prefix);
        return await Provider.of<DmApiClient>(context, listen: false)
            .getEntries(dir)
            .where((e) => e.isDir)
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

    Future(() async => await updateViewForPath("/"));
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        SizedBox(
          height: 80,
          child: Row(children: [
            IconButton(
              onPressed: () => dirUp(),
              icon: const Icon(Icons.arrow_back),
            ),
            Expanded(
              child: RawAutocomplete<String>(
                  focusNode: FocusNode(),
                  textEditingController: _pathController,
                  onSelected: (option) {
                    updateViewForPath(option);
                  },
                  fieldViewBuilder: (
                    BuildContext context,
                    TextEditingController textEditingController,
                    FocusNode focusNode,
                    VoidCallback onFieldSubmitted,
                  ) {
                    return TextFormField(
                      controller: textEditingController,
                      focusNode: focusNode,
                      onFieldSubmitted: (String value) {
                        onFieldSubmitted();
                      },
                    );
                  },
                  optionsViewBuilder: (
                    BuildContext context,
                    AutocompleteOnSelected<String> onSelected,
                    Iterable<String> options,
                  ) {
                    return _AutocompleteOptions<String>(
                      displayStringForOption: RawAutocomplete.defaultStringForOption,
                      onSelected: onSelected,
                      options: options,
                      maxOptionsHeight: 200,
                    );
                  },
                  optionsBuilder: (textEditingValue) => suggestions(textEditingValue.text)),
            ),
          ]),
        ),
        Expanded(
          child: GridView.extent(
            maxCrossAxisExtent: 200,
            mainAxisSpacing: 10,
            crossAxisSpacing: 10,
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

// The default Material-style Autocomplete options.
class _AutocompleteOptions<T extends Object> extends StatelessWidget {
  const _AutocompleteOptions({
    Key? key,
    required this.displayStringForOption,
    required this.onSelected,
    required this.options,
    required this.maxOptionsHeight,
  }) : super(key: key);

  final AutocompleteOptionToString<T> displayStringForOption;

  final AutocompleteOnSelected<T> onSelected;

  final Iterable<T> options;
  final double maxOptionsHeight;

  @override
  Widget build(BuildContext context) {
    return Align(
      alignment: Alignment.topLeft,
      child: Material(
        elevation: 4.0,
        child: ConstrainedBox(
          constraints: BoxConstraints(maxHeight: maxOptionsHeight),
          child: ListView.builder(
            padding: EdgeInsets.zero,
            shrinkWrap: true,
            itemCount: options.length,
            itemBuilder: (BuildContext context, int index) {
              final T option = options.elementAt(index);
              return InkWell(
                onTap: () {
                  onSelected(option);
                },
                child: Builder(builder: (BuildContext context) {
                  final bool highlight = AutocompleteHighlightedOption.of(context) == index;
                  if (highlight) {
                    SchedulerBinding.instance.addPostFrameCallback((Duration timeStamp) {
                      Scrollable.ensureVisible(context, alignment: 0.5);
                    });
                  }
                  return Container(
                    color: highlight ? Theme.of(context).focusColor : null,
                    padding: const EdgeInsets.all(16.0),
                    child: Text(displayStringForOption(option)),
                  );
                }),
              );
            },
          ),
        ),
      ),
    );
  }
}
