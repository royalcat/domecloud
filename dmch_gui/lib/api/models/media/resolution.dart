import 'dart:convert';

class Resolution {
  final int height;
  final int width;

  Resolution(this.height, this.width);

  Map<String, dynamic> toMap() {
    return {
      'height': height,
      'width': width,
    };
  }

  factory Resolution.fromMap(Map<String, dynamic> map) {
    return Resolution(
      map['height'] as int? ?? 0,
      map['width'] as int? ?? 0,
    );
  }

  String toJson() => json.encode(toMap());

  factory Resolution.fromJson(String source) =>
      Resolution.fromMap(json.decode(source) as Map<String, dynamic>);
}
