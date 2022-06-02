extension DurationBatteries on Duration {
  int get inNanoseconds => inMicroseconds * 1000;
  static Duration fromNanoseconds(int nanoseconds) => Duration(microseconds: nanoseconds ~/ 1000);
}
