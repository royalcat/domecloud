import 'package:bloc/bloc.dart';

import 'package:dmch_gui/models/entry.dart';
import 'package:dmch_gui/models/media.dart';
import 'package:dmch_gui/provider/dmapi.dart';
import 'package:dmch_gui/widgets/media/video.dart';

class VideoItemBloc extends Bloc<VideoItemEvent, VideoItemState> {
  final DmApiClient client;

  VideoItemBloc(this.client, Entry entry) : super(VideoItemState.init(entry)) {
    on<VideoItemEventItemShown>((event, emit) async {
      final info = await client.getVideoInfo(state.entry.filePath);
      emit(state.copyWith(info: info));
    });

    on<VideoItemEventHover>((event, emit) async {
      final previews = await client.getPreviews(state.entry.filePath);
      emit(state.copyWith(previews: previews));
    });
  }
}

class VideoItemState {
  final VideoInfo? info;
  final List<Entry>? previews;
  final Entry entry;
  const VideoItemState({
    required this.entry,
    this.info,
    this.previews,
  });

  const VideoItemState.init(this.entry)
      : info = null,
        previews = null;

  VideoItemState copyWith({
    VideoInfo? info,
    List<Entry>? previews,
    Entry? entry,
  }) {
    return VideoItemState(
      info: info ?? this.info,
      previews: previews ?? this.previews,
      entry: entry ?? this.entry,
    );
  }
}

class VideoItemEventHover extends VideoItemEvent {}

abstract class VideoItemEvent {}

class VideoItemEventItemShown extends VideoItemEvent {}
