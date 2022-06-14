import 'package:flutter/material.dart';

class AppScaffold extends StatelessWidget {
  const AppScaffold({
    Key? key,
    required this.body,
    this.appBar,
    this.drawer,
    this.endDrawer,
    this.floatingActionButton,
  }) : super(key: key);

  final Widget body;

  final PreferredSizeWidget? appBar;
  final Widget? drawer;
  final Widget? endDrawer;
  final Widget? floatingActionButton;

  @override
  Widget build(BuildContext context) {
    final bool displayMobileLayout = MediaQuery.of(context).size.width < 600;
    return Row(
      children: [
        if (!displayMobileLayout && drawer != null) drawer!,
        Expanded(
          child: Scaffold(
            appBar: appBar,
            drawer: displayMobileLayout ? drawer : null,
            endDrawer: displayMobileLayout ? endDrawer : null,
            body: body,
            floatingActionButton: floatingActionButton,
          ),
        ),
        if (!displayMobileLayout && endDrawer != null) endDrawer!,
      ],
    );
  }
}
