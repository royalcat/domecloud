import 'dart:async';
import 'dart:convert';

import 'package:dmch_gui/api/dmapi.dart';
import 'package:dmch_gui/api/models/entry.dart';
import 'package:dmch_gui/api/models/media/video.dart';
import 'package:dmch_gui/api/models/meta.dart';
import 'package:dmch_gui/widgets/app_scaffold/app_scaffold.dart';
import 'package:dmch_gui/widgets/media/folder.dart';
import 'package:dmch_gui/widgets/media/video.dart';
import 'package:flutter/material.dart';
import 'package:flutter/scheduler.dart';
import 'package:path/path.dart' as path_utils;
import 'package:provider/provider.dart';

class MediaGrid extends StatefulWidget {
  final String basePath;

  const MediaGrid({Key? key, required this.basePath}) : super(key: key);

  @override
  State<MediaGrid> createState() => _MediaGridState();
}

const _avalibleSorts = ["Name", "Size", "Date"];

class _MediaGridState extends State<MediaGrid> {
  final _pathController = TextEditingController();

  List<Entry> _entries = [];
  bool loading = true;
  bool canDirUp = false;
  String sort = _avalibleSorts[0];
  Entry? selectedEntry;

  @override
  void initState() {
    super.initState();
    updateViewForPath(widget.basePath);
  }

  void dirUp() => updateViewForPath(path_utils.dirname(_pathController.text));

  void dirDown(String dir) => updateViewForPath(path_utils.joinAll([_pathController.text, dir]));

  Future<void> updateViewForPath(String path) async {
    try {
      setState(() {
        loading = true;
      });

      final entries =
          await Provider.of<DmApiClient>(context, listen: false).getEntries(path).toList();

      setState(() {
        _pathController.text = path;
        _entries = entries;
        loading = false;
        canDirUp = path != widget.basePath;
      });
    } catch (e) {
      debugPrint("exception for path: ${_pathController.text}: $e");
    }
  }

  @override
  Widget build(BuildContext context) {
    final dmapi = Provider.of<DmApiClient>(context);
    final ThemeData theme = Theme.of(context);

    return AppScaffold(
      endDrawer: Drawer(
        //width: 300,
        child: Column(
          children: [
            UserAccountsDrawerHeader(
              currentAccountPicture: const CircleAvatar(
                child: Icon(Icons.person),
              ),
              accountName: Text(dmapi.user?.name ?? ""),
              accountEmail: Text("@${dmapi.user?.name ?? ""}"),
            ),
            if (selectedEntry != null) renderDetails(context, selectedEntry!),
          ],
        ),
      ),
      appBar: AppBar(
        leading: IconButton(
          onPressed: () => dirUp(),
          icon: const Icon(Icons.arrow_back),
        ),
        title: _PathAutoComplete(
          pathController: _pathController,
          onPathChanged: updateViewForPath,
          basePath: widget.basePath,
        ),
      ),
      floatingActionButton: Row(
        mainAxisAlignment: MainAxisAlignment.end,
        children: const [
          FloatingActionButton(
            onPressed: null,
            child: Icon(Icons.cloud_upload),
          ),
          SizedBox(width: 10),
          FloatingActionButton(
            onPressed: null,
            child: Icon(Icons.create_new_folder),
          ),
        ],
      ),
      body: !loading
          ? Row(
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
                    child: Column(
                      children: [
                        SizedBox(
                          height: 40,
                          child: Row(
                            children: [
                              const Text("Sort by: "),
                              DropdownButton<String>(
                                onChanged: (value) => setState(() {
                                  sort = value ?? _avalibleSorts[0];
                                }),
                                value: sort,
                                items: _avalibleSorts
                                    .map<DropdownMenuItem<String>>(
                                      (e) => DropdownMenuItem(
                                        value: e,
                                        child: Text(e),
                                      ),
                                    )
                                    .toList(),
                              ),
                            ],
                          ),
                        ),
                        Expanded(
                          child: GridView.extent(
                            maxCrossAxisExtent: 200,
                            mainAxisSpacing: 10,
                            crossAxisSpacing: 10,
                            children: _entries
                                .where((element) => !element.isDir)
                                .map(
                                  (e) => VideoInfoItem(
                                    entry: e,
                                    onOpenDetails: (e) => setState(() {
                                      selectedEntry = e;
                                    }),
                                  ),
                                )
                                .toList(),
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ],
            )
          : const Center(child: CircularProgressIndicator()),
    );
  }

  Widget renderDetails(BuildContext context, Entry entry) {
    final dmapi = Provider.of<DmApiClient>(context);

    return FutureBuilder<EntryMeta>(
      future: dmapi.getVideoInfo(entry.path),
      builder: (context, snapshot) {
        return Column(
          children: [
            ListTile(
              leading: const Text("Name:"),
              title: Text(entry.name),
            ),
            ...snapshot.hasData
                ? [
                    ListTile(
                      leading: const Text("Duration:"),
                      title: Text(snapshot.data!.mediaInfo?.videoInfo?.duration?.toString() ?? ""),
                    ),
                    ListTile(
                      leading: const Text("Resolution:"),
                      title:
                          Text(snapshot.data!.mediaInfo?.videoInfo?.resolution?.toString() ?? ""),
                    ),
                  ]
                : [const CircularProgressIndicator()],
          ],
        );
      },
    );
  }
}

class _PathAutoComplete extends StatelessWidget {
  const _PathAutoComplete({
    Key? key,
    required this.pathController,
    required this.onPathChanged,
    required this.basePath,
  }) : super(key: key);

  final String basePath;
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
    final dmapi = Provider.of<DmApiClient>(context, listen: false);
    if (prefix.isEmpty || prefix == basePath) {
      return [];
    }

    try {
      if (prefix.endsWith("/")) {
        return [prefix] +
            await dmapi
                .getEntries(prefix)
                .where((e) => e.isDir)
                .map<String>((e) => path_utils.joinAll([prefix, e.name]))
                .toList();
      } else {
        final dir = path_utils.dirname(prefix);
        final query = path_utils.basename(prefix);
        return await dmapi
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
