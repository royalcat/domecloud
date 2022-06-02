import 'package:dome_server/models/entry/entry.dart';
import 'package:mongo_dart/mongo_dart.dart';

class DomeFiles {
  final DbCollection coll;

  DomeFiles(Db db) : coll = db.collection('files');

  Future<void> set(Entry entry) async {
    await coll.replaceOne({"path", entry.path}, entry.toJson());
  }

  Stream<Entry> getMediaInDir(String targetDir) async* {}
}
