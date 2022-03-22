import 'dart:convert';

import 'package:dmch_gui/models/entry.dart';
import 'package:flutter/cupertino.dart';

import 'package:path/path.dart' as path_utils;

import 'package:http/http.dart' as http;

import '../models/media.dart';

class DmApiClient {
  final baseUrl = 'http://localhost:5050/file/';

  final path_utils.Context ctx = path_utils.posix;
  final http.Client _client = http.Client();

  Future<VideoInfo> getVideoInfo(String fpath) async {
    final resp = await _request(ctx.join(fpath, "info"));
    return VideoInfo.fromJson(resp.body);
  }

  Stream<VideoInfo> getVideoInfos(String dir, List<Entry> entries) async* {
    for (final entry in entries.where((e) => e.mediaType == MediaTypes.video)) {
      yield await getVideoInfo(ctx.joinAll([dir, entry.name]));
    }
  }

  Future<List<Entry>> getPreviews(String fpath) async {
    return getEntries(ctx.joinAll([fpath, "previews"]));
  }

  Future<List<Entry>> getEntries(String dir) async {
    final resp = await _request(dir);
    return (json.decode(resp.body) as List<dynamic>)
        .map((e) => Entry.fromMap(e, ctx.joinAll([dir, e])))
        .toList();
  }

  Future<http.Response> _request(path) async {
    final uri = Uri.parse(getUrlFromFilepath(path));
    return _client.get(uri);
  }

  String getUrlFromFilepath(String fpath) => path_utils.url.joinAll([
        baseUrl,
        fpath == "/" ? "" : fpath,
      ]);
}
