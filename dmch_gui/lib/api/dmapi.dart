import 'dart:convert';

import 'package:dmch_gui/api/models/entry.dart';
import 'package:dmch_gui/api/models/media/media.dart';
import 'package:dmch_gui/api/models/media/video.dart';
import 'package:dmch_gui/api/models/meta.dart';
import 'package:dmch_gui/api/models/users.dart';
import 'package:http/http.dart' as http;
import 'package:path/path.dart' as path_utils;

class DmApiClient {
  Map<String, String> authHeader = <String, String>{};

  User? _user;
  User? get user => _user;
  bool get isLoggedIn => _user != null;

  final scheme = "http";
  final host = 'localhost';
  final port = 5050;
  final filePathBase = "file";
  final playerPathBase = "player";

  String get userDir => "/user/${user?.name}";

  final path_utils.Context ctx = path_utils.url;
  final http.Client _client = http.Client();

  DmApiClient();

  Future<EntryMeta> getVideoInfo(String fpath) async {
    final resp = await _requestFile(ctx.join(fpath, "info.json"));

    return EntryMeta.fromMap(json.decode(resp.body) as Map<String, dynamic>);
  }

  Stream<EntryMeta> getVideoInfos(String dir, List<Entry> entries) async* {
    for (final entry in entries) {
      yield await getVideoInfo(ctx.joinAll([dir, entry.name]));
    }
  }

  Stream<Entry> getPreviews(String fpath) {
    return getEntries(ctx.joinAll([fpath, "previews"]));
  }

  Stream<Entry> listMediaPreviews(String fpath) async* {
    final previewsDir = ctx.joinAll([fpath, "previews"]);
    final resp = await _requestApi("listMedia", ctx.joinAll([fpath, "previews"]));

    yield* Stream.fromIterable(
      (json.decode(resp.body) as List<dynamic>)
          .map((e) => Entry.fromMap(previewsDir, e as Map<String, dynamic>)),
    );
  }

  Stream<Entry> getEntries(String dir) async* {
    final resp = await _requestFile(dir);

    yield* Stream.fromIterable(
      (json.decode(resp.body) as List<dynamic>)
          .map((e) => Entry.fromMap(dir, e as Map<String, dynamic>)),
    );
  }

  Future<Uri> getPlayUrl(String fpath) async {
    final uri = _getPlayerUri("regtoken", fpath);
    final resp = await _request(uri);

    return _getPlayerUri("play", fpath, token: resp.body);
  }

  Future<http.Response> _requestApi(String command, String path) async {
    if (!isLoggedIn) {
      throw NotAuthenticatedException();
    }

    return _request(_getUriApiFromFilepath(command, path));
  }

  Future<http.Response> _requestFile(String path) async {
    if (!isLoggedIn) {
      throw NotAuthenticatedException();
    }

    return _request(getFileUri(path));
  }

  Future<http.Response> _request(Uri uri) async {
    final resp = await _client.get(
      uri,
      headers: {
        ...authHeader,
      },
    );

    switch (resp.statusCode) {
      case 401:
        throw NotAuthenticatedException();
      default:
        return resp;
    }
  }

  Uri getFileUri(String fpath) => Uri(
        scheme: scheme,
        host: host,
        port: port,
        path: ctx.joinAll([filePathBase, fpath.trimLeading("/")]),
      );

  Uri _getUriApiFromFilepath(String command, String fpath) => Uri(
        scheme: scheme,
        host: host,
        port: port,
        path: ctx.joinAll([filePathBase, command, fpath.trimLeading("/")]),
      );

  Uri _getPlayerUri(String command, String fpath, {String? token}) => Uri(
        scheme: scheme,
        host: host,
        port: port,
        path: ctx.joinAll([playerPathBase, command, fpath.trimLeading("/")]),
        queryParameters: token != null
            ? {
                "token": token,
              }
            : null,
      );

  Future<User?> logIn({
    required String username,
    required String password,
  }) async {
    authHeader["Authorization"] = 'Basic ${base64Encode(utf8.encode('$username:$password'))}';
    try {
      final resp = await _request(
        Uri(
          scheme: scheme,
          host: host,
          port: port,
          path: "/login",
        ),
      );
      _user = User.fromJson(resp.body);
      return user!;
    } on NotAuthenticatedException {
      return null;
    }
  }

  Future<void> logOut() async {
    authHeader.remove("Authorization");
    _user = null;
  }
}

class DmApiException implements Exception {}

class DmApiRequestException implements DmApiException {
  final Uri uri;
  final http.Response response;
  DmApiRequestException({
    required this.uri,
    required this.response,
  });
}

class NotAuthenticatedException implements DmApiException {}

extension _StringTrim on String {
  String trimLeading(String pattern) {
    int i = 0;
    while (startsWith(pattern, i)) {
      i += pattern.length;
    }
    return substring(i);
  }
}
