import 'dart:async';

import 'package:dmch_gui/api/dmapi.dart';
import 'package:dmch_gui/api/models/entry.dart';
import 'package:dmch_gui/widgets/media/folder.dart';
import 'package:dmch_gui/widgets/media/video.dart';
import 'package:flutter/material.dart';
import 'package:flutter/scheduler.dart';
import 'package:path/path.dart' as path_utils;
import 'package:provider/provider.dart';

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
      _pathController.text = path_utils.normalize(path);
      _entries = await Provider.of<DmApiClient>(context, listen: false)
          .getEntries(_pathController.text)
          .toList();
      setState(() {});
    } catch (e) {
      debugPrint("exception for path: ${_pathController.text}: $e");
      await dirUp();
    }
  }

  @override
  void initState() {
    super.initState();

    Future(() async => updateViewForPath("/"));
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        SizedBox(
          height: 80,
          child: Row(
            children: [
              IconButton(
                onPressed: () => dirUp(),
                icon: const Icon(Icons.arrow_back),
              ),
              Expanded(
                child: _PathAutoComplete(
                  pathController: _pathController,
                  onPathChanged: updateViewForPath,
                ),
              ),
            ],
          ),
        ),
        const Divider(
          height: 1,
          thickness: 1,
        ),
        Expanded(
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Flexible(
                flex: 2,
                child: FolderList(
                  entries: _entries.where((element) => element.isDir).toList(),
                  onOpen: (entry) => dirDown(entry.name),
                ),
              ),
              const VerticalDivider(
                width: 1,
                thickness: 1,
              ),
              Flexible(
                flex: 10,
                child: Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: GridView.extent(
                    maxCrossAxisExtent: 200,
                    mainAxisSpacing: 10,
                    crossAxisSpacing: 10,
                    children: _entries
                        .where((element) => !element.isDir)
                        .map((e) => VideoInfoItem(
                              entry: e,
                              dirPath: _pathController.text,
                            ))
                        .toList(),
                  ),
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }
}

class _PathAutoComplete extends StatelessWidget {
  const _PathAutoComplete({Key? key, required this.pathController, required this.onPathChanged})
      : super(key: key);

  final TextEditingController pathController;
  final void Function(String) onPathChanged;

  @override
  Widget build(BuildContext context) {
    return RawAutocomplete<String>(
      focusNode: FocusNode(),
      textEditingController: pathController,
      onSelected: onPathChanged,
      fieldViewBuilder: (
        BuildContext context,
        TextEditingController textEditingController,
        FocusNode focusNode,
        VoidCallback onFieldSubmitted,
      ) {
        return TextFormField(
          controller: textEditingController,
          focusNode: focusNode,
          decoration: const InputDecoration(
            border: InputBorder.none,
          ),
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
      optionsBuilder: (textEditingValue) => suggestions(context, textEditingValue.text),
    );
  }

  Future<List<String>> suggestions(BuildContext context, String prefix) async {
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
      debugPrint("exception for path: ${pathController.text}: $e");
    }

    return <String>[];
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
