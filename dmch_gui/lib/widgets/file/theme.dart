import 'package:equatable/equatable.dart';
import 'package:flutter/material.dart';

class FSGridTheme extends InheritedWidget {
  final FSGridThemeData theme;

  const FSGridTheme({
    Key? key,
    required Widget child,
    required this.theme,
  }) : super(key: key, child: child);

  static FSGridThemeData of(BuildContext context) {
    final FSGridTheme? result = context.dependOnInheritedWidgetOfExactType<FSGridTheme>();
    assert(result != null, 'No FrogColor found in context');
    return result!.theme;
  }

  @override
  bool updateShouldNotify(FSGridTheme oldWidget) => theme != oldWidget.theme;
}

class FSGridThemeData extends Equatable {
  final IconData fileIcon;
  final IconData folderIcon;
  final double scale;

  factory FSGridThemeData.normal() => const FSGridThemeData(
        fileIcon: Icons.insert_drive_file,
        folderIcon: Icons.folder,
        scale: 1,
      );

  factory FSGridThemeData.sharp() => const FSGridThemeData(
        fileIcon: Icons.insert_drive_file_sharp,
        folderIcon: Icons.folder_sharp,
        scale: 1,
      );

  const FSGridThemeData({
    required this.fileIcon,
    required this.folderIcon,
    required this.scale,
  });

  @override
  List<Object?> get props => [fileIcon, folderIcon, scale];
}
