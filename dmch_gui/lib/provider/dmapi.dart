import 'dart:convert';

import 'package:dmch_gui/models/entry.dart';
import 'package:flutter/cupertino.dart';

import 'package:path/path.dart' as path_utils;

import 'package:http/http.dart' as http;

import '../models/media.dart';

class DmApiClient {
  final host = 'http://localhost:5050';
  final baseUrl = "/file";

  final path_utils.Context ctx = path_utils.url;
  final http.Client _client = http.Client();

  Future<VideoInfo> getVideoInfo(String fpath) async {
    final resp = await _request(ctx.join(fpath, "info.json"));
    return VideoInfo.fromJson(resp.body);
  }

  Stream<VideoInfo> getVideoInfos(String dir, List<Entry> entries) async* {
    for (final entry in entries) {
      yield await getVideoInfo(ctx.joinAll([dir, entry.name]));
    }
  }

  Stream<Entry> getPreviews(String fpath) {
    return getEntries(ctx.joinAll([fpath, "previews"]));
  }

  Stream<Entry> getEntries(String dir) async* {
    final resp = await _request(dir);

    yield* Stream.fromIterable(
      (json.decode(resp.body) as List<dynamic>).map((e) => Entry.fromMap(dir, e)),
    );
  }

  Future<http.Response> _request(path) async {
    return await _client.get(Uri.parse(getUrlFromFilepath(path)));
  }

  String getUrlFromFilepath(String fpath) => host + ctx.joinAll([baseUrl, fpath.trimLeading("/")]);
}

extension StringTrim on String {
  String trimLeading(String pattern) {
    int i = 0;
    while (startsWith(pattern, i)) {
      i += pattern.length;
    }
    return substring(i);
  }
}
