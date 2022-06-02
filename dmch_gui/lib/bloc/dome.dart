import 'package:dmch_gui/api/dmapi.dart';
import 'package:dmch_gui/api/models/entry.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:stream_bloc/stream_bloc.dart';
import 'package:path/path.dart' as path_utils;

part 'dome.freezed.dart';

abstract class DomeEvent {}

abstract class InitDomeEvent {}

class DirDownDomeEvent implements DomeEvent {}

@freezed
class DirUpDomeEvent with _$DirUpDomeEvent implements DomeEvent {
  const factory DirUpDomeEvent({
    required String dirName,
  }) = _DirUpDomeEvent;
}

abstract class DomeState {}

@freezed
class LoadingDomeState with _$LoadingDomeState implements DomeState {
  const factory LoadingDomeState() = _LoadingDomeState;
}

@freezed
class ContentDomeState with _$ContentDomeState implements DomeState {
  const factory ContentDomeState({
    required String path,
    required List<Entry> entries,
  }) = _ContentDomeState;
}

class DomeBloc extends StreamBloc<DomeEvent, DomeState> {
  final DmApiClient apiClient;

  var currentPath = "/";

  DomeBloc(this.apiClient) : super(const LoadingDomeState());

  void _dirUp() {
    if (currentPath != "/") {
      path_utils.dirname(currentPath);
    }
  }

  void _dirDown(String dir) {
    path_utils.joinAll([currentPath, dir]);
  }

  Future<List<Entry>> _getEntries() => apiClient.getEntries(currentPath).toList();

  @override
  Stream<DomeState> mapEventToStates(DomeEvent event) async* {
    if (event is InitDomeEvent) {
    } else if (event is DirUpDomeEvent) {
      yield const LoadingDomeState();
      _dirUp();
      final entries = await _getEntries();
      yield ContentDomeState(path: currentPath, entries: entries);
    } else if (event is DirDownDomeEvent) {
      yield const LoadingDomeState();

      //_dirDown(event);
    }
  }
}
